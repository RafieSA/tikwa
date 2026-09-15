package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ValidateInput: cek edge cases sebelum proses.
// - file tidak ada
// - bukan file video (cek ekstensi + cek ffmpeg)
// - path traversal / injection (kita tidak pakai shell, pakai exec.Command arg terpisah = aman dari injection)
func ValidateInput(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path kosong — pilih video dulu")
	}
	clean := filepath.Clean(path)
	info, err := os.Stat(clean)
	if err != nil {
		return fmt.Errorf("file tidak ditemukan: %s", clean)
	}
	if info.IsDir() {
		return fmt.Errorf("yang dipilih adalah folder, bukan file video: %s", clean)
	}
	if info.Size() == 0 {
		return fmt.Errorf("file kosong (0 byte): %s", clean)
	}
	// Cek ekstensi umum — bukan validasi final, hanya early feedback
	ext := strings.ToLower(filepath.Ext(clean))
	allowed := map[string]bool{".mp4": true, ".mov": true, ".m4v": true, ".mkv": true, ".avi": true, ".webm": true, ".3gp": true}
	if !allowed[ext] {
		return fmt.Errorf("format %s belum didukung — pakai .mp4/.mov/.mkv/.avi/.webm", ext)
	}
	return nil
}

// CheckFFmpeg: pastikan ffmpeg ada di PATH (audit awal yang kamu minta).
func CheckFFmpeg() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg tidak ditemukan — install dulu: brew install ffmpeg")
	}
	return nil
}

// BuildArgs: susun argumen FFmpeg Level 1 Basic Polish.
// - scale + pad ke target (jaga aspect ratio, tidak stretch)
// - unsharp = tajamkan dikit (biar tidak burem setelah kompres)
// - eq = terangkan kalau gelap (brightness 0.06 = +6%)
// YAGNI: belum pakai AI upscale / denoise berat di v1.
func BuildArgs(input, output string, p Profile, f Filter) []string {
	// vf = scale+pad + filter gaya (Natural/Dramatis/Cinematic/Original)
	baseVF := fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease,pad=%d:%d:(%d-iw)/2:(%d-ih)/2:color=black", p.Width, p.Height, p.Width, p.Height, p.Width, p.Height)
	extra := filterVF(f)
	vf := baseVF
	if extra != "" {
		vf = baseVF + "," + extra
	}

	return []string{
		"-y", // overwrite output
		"-i", input,
		"-vf", vf,
		"-c:v", "libx264",
		"-preset", "medium", // tradeoff: medium = seimbang cepat & kualitas
		"-crf", "23", // 18=jernih besar, 28=kecil burem. 23 = sweet spot Level 1
		"-maxrate", p.MaxBitrate,
		"-bufsize", p.MaxBitrate,
		"-c:a", "aac",
		"-b:a", p.AudioBitrate,
		"-movflags", "+faststart", // biar bisa streaming langsung pas upload
		output,
	}
}

// Run: eksekusi FFmpeg dengan context (bisa cancel) + log aman.
// Tidak pakai shell, jadi aman dari injection (SSRF/XSS tidak relevan di CLI, tapi command injection kita cegah).
func Run(ctx context.Context, input, output string, p Profile, f Filter) error {
	if f == FilterRequest {
		return fmt.Errorf("filter request tidak diproses FFmpeg — pakai SaveRequest")
	}
	if !f.IsValid() {
		return fmt.Errorf("filter tidak valid: %s", f)
	}
	if err := ValidateInput(input); err != nil {
		return err
	}
	if err := CheckFFmpeg(); err != nil {
		return err
	}
	args := BuildArgs(input, output, p, f)
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	// gabung stderr ke output biar bisa di-log kalau gagal
	out, err := cmd.CombinedOutput()
	if err != nil {
		// jangan bocorkan path absolut sensitif ke user? tapi untuk CLI lokal, tampilkan agar debug mudah.
		return fmt.Errorf("ffmpeg gagal: %v\n%s", err, string(out))
	}
	return nil
}
