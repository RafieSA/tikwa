package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RequestFile: lokasi simpan titipan ide (di folder config user, bukan di repo).
// Biar tidak bocor ke git, dan tidak perlu permission tulis di /Applications.
func RequestFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "tikwa-requests.txt"
	}
	return filepath.Join(home, ".tikwa-requests.txt")
}

// SaveRequest: simpan ide filter user dengan validasi & sanitasi.
// - trim spasi
// - tolak kosong
// - batasi panjang 200 biar tidak bom file
// - sanitasi newline supaya 1 baris = 1 request (no injection log)
// - append dengan timestamp
func SaveRequest(idea string) (string, error) {
	idea = strings.TrimSpace(idea)
	if idea == "" {
		return "", fmt.Errorf("ide tidak boleh kosong — tulis misal: vintage BW 90an")
	}
	if len(idea) > 200 {
		return "", fmt.Errorf("ide kepanjangan (>200 huruf) — perpendek ya")
	}
	// sanitasi: ganti newline jadi spasi
	idea = strings.ReplaceAll(idea, "\n", " ")
	idea = strings.ReplaceAll(idea, "\r", " ")

	path := RequestFile()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", fmt.Errorf("gagal simpan request: %w", err)
	}
	defer f.Close()

	line := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04"), idea)
	if _, err := f.WriteString(line); err != nil {
		return "", err
	}
	return path, nil
}
