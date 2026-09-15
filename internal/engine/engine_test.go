package engine

import (
	"os"
	"strings"
	"testing"
)

// TestProfile: pastikan cetakan TikTok & WA tidak salah ukuran.
// Kalau salah, video jadi stretch / burem lagi.
func TestProfile(t *testing.T) {
	tik := TikTokProfile()
	if tik.Width != 1080 || tik.Height != 1920 {
		t.Fatalf("tiktok: want 1080x1920, got %dx%d", tik.Width, tik.Height)
	}
	if tik.Name != "tiktok" {
		t.Fatalf("tiktok name wrong: %s", tik.Name)
	}
	wa := WhatsAppProfile()
	if wa.Width != 720 || wa.Height != 1280 {
		t.Fatalf("whatsapp: want 720x1280, got %dx%d", wa.Width, wa.Height)
	}
	if OutputPath("/tmp/a.mp4", tik) != "/tmp/a-tiktok.mp4" {
		t.Fatalf("output tiktok path wrong: %s", OutputPath("/tmp/a.mp4", tik))
	}
	if OutputPath("/tmp/my video", wa) != "/tmp/my video-whatsapp.mp4" {
		t.Fatalf("output no-ext wrong: %s", OutputPath("/tmp/my video", wa))
	}
	// file dengan spasi harus aman
	if !strings.Contains(OutputPath("/tmp/video gue.mp4", wa), "video gue-whatsapp") {
		t.Fatalf("spasi handling failed")
	}
}

func TestValidateInput(t *testing.T) {
	if err := ValidateInput(""); err == nil {
		t.Fatal("empty path harus error")
	}
	if err := ValidateInput("   "); err == nil {
		t.Fatal("spasi doang harus error")
	}
	if err := ValidateInput("/tmp"); err == nil {
		t.Fatal("folder harus error")
	}
	if err := ValidateInput("/tmp/notfound-tikwa-xyz123.mp4"); err == nil {
		t.Fatal("file tidak ada harus error")
	}
	// file .jpg harus ditolak (bukan video)
	f, _ := os.CreateTemp("", "tikwa*.jpg")
	f.Close()
	defer os.Remove(f.Name())
	if err := ValidateInput(f.Name()); err == nil {
		t.Fatal("jpg harus ditolak")
	}
	// file kosong 0 byte harus ditolak
	empty, _ := os.CreateTemp("", "tikwa-empty*.mp4")
	empty.Close()
	defer os.Remove(empty.Name())
	if err := ValidateInput(empty.Name()); err == nil {
		t.Fatal("file 0 byte harus error")
	}
}

func TestBuildArgs(t *testing.T) {
	args := BuildArgs("/tmp/in.mp4", "/tmp/out.mp4", TikTokProfile(), FilterNatural)
	if len(args) < 10 {
		t.Fatalf("args terlalu pendek: %v", args)
	}
	// pastikan tidak ada arg kosong (anti injection edge)
	for _, a := range args {
		if a == "" {
			t.Fatal("arg kosong tidak boleh")
		}
	}
	// pastikan vf mengandung scale & pad (jaga aspect ratio)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "scale=") || !strings.Contains(joined, "pad=") {
		t.Fatal("vf harus ada scale & pad")
	}
	if !strings.Contains(joined, "libx264") {
		t.Fatal("codec harus libx264")
	}
}
func TestFilter(t *testing.T) {
	for _, f := range []Filter{FilterNatural, FilterDramatis, FilterCinematic, FilterOriginal} {
		if !f.IsValid() {
			t.Fatalf("%s harus valid", f)
		}
		args := BuildArgs("/tmp/in.mp4", "/tmp/out.mp4", TikTokProfile(), f)
		joined := strings.Join(args, " ")
		if !strings.Contains(joined, "scale=") {
			t.Fatalf("%s harus ada scale", f)
		}
	}
	if FilterRequest.IsValid() == false {
		t.Fatal("request harus valid tapi tidak dipoles")
	}
	// Original tidak boleh ada eq/unsharp
	orig := strings.Join(BuildArgs("/tmp/in.mp4", "/tmp/out.mp4", WhatsAppProfile(), FilterOriginal), " ")
	if strings.Contains(orig, "eq=") || strings.Contains(orig, "unsharp=") {
		t.Fatal("original tidak boleh ada eq/unsharp")
	}
	// Dramatis harus ada vignette & contrast tinggi
	dram := strings.Join(BuildArgs("/tmp/in.mp4", "/tmp/out.mp4", TikTokProfile(), FilterDramatis), " ")
	if !strings.Contains(dram, "vignette") || !strings.Contains(dram, "1.25") {
		t.Fatal("dramatis harus ada vignette & contrast 1.25")
	}
	// Cinematic harus ada colorbalance & curves
	cin := strings.Join(BuildArgs("/tmp/in.mp4", "/tmp/out.mp4", TikTokProfile(), FilterCinematic), " ")
	if !strings.Contains(cin, "colorbalance") || !strings.Contains(cin, "curves") {
		t.Fatal("cinematic harus ada colorbalance & curves")
	}
}

func TestSaveRequest(t *testing.T) {
	if _, err := SaveRequest(""); err == nil {
		t.Fatal("kosong harus error")
	}
	if _, err := SaveRequest("   "); err == nil {
		t.Fatal("spasi harus error")
	}
	long := strings.Repeat("a", 201)
	if _, err := SaveRequest(long); err == nil {
		t.Fatal("kepanjangan harus error")
	}
	// happy path
	path, err := SaveRequest("vintage BW 90an test")
	if err != nil {
		t.Fatalf("save request gagal: %v", err)
	}
	if path == "" {
		t.Fatal("path kosong")
	}
	// cleanup: hapus baris test (baca & tidak perlu hapus file, biar user lihat)
}

