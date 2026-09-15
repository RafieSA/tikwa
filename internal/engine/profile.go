package engine

// Profile = cetakan setting FFmpeg untuk tiap platform.
// Satu fungsi = satu tanggung jawab (SRP).
type Profile struct {
	Name        string // "tiktok" | "whatsapp"
	Width       int
	Height      int
	MaxBitrate  string // contoh "4000k"
	AudioBitrate string
	MaxSizeMB   int
}

// TikTokProfile: vertikal 9:16, 1080x1920, tajam & optimal agar tidak diperas lagi.
// Tradeoff: file lebih besar dari WA, tapi tajam maksimal di TikTok.
func TikTokProfile() Profile {
	return Profile{
		Name:        "tiktok",
		Width:       1080,
		Height:      1920,
		MaxBitrate:  "5000k",
		AudioBitrate: "128k",
		MaxSizeMB:   50,
	}
}

// WhatsAppProfile: 720p, kecil, tetap jernih. Lolos limit 64MB dengan aman.
// Tradeoff: resolusi lebih rendah dari TikTok, tapi kecil & cepat kekirim.
func WhatsAppProfile() Profile {
	return Profile{
		Name:        "whatsapp",
		Width:       720,
		Height:      1280,
		MaxBitrate:  "2500k",
		AudioBitrate: "96k",
		MaxSizeMB:   25,
	}
}

// OutputPath: bikin nama file output yang aman.
// contoh: /tmp/video.mp4 + tiktok -> /tmp/video-tiktok.mp4
func OutputPath(inputPath string, p Profile) string {
	ext := ".mp4"
	base := inputPath
	// potong ekstensi jika ada
	for i := len(inputPath) - 1; i >= 0; i-- {
		if inputPath[i] == '.' {
			ext = inputPath[i:]
			base = inputPath[:i]
			break
		}
		if inputPath[i] == '/' {
			break
		}
	}
	return base + "-" + p.Name + ext
}
