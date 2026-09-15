# ▶ TIKWA — Bikin Video HP Tidak Burem

<p align="center">
  <a href="https://github.com/RafieSA/tikwa/actions/workflows/ci.yml"><img src="https://github.com/RafieSA/tikwa/actions/workflows/ci.yml/badge.svg" alt="CI" /></a>
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go" alt="Go 1.26" />
  <img src="https://img.shields.io/badge/TUI-Bubble_Tea-FF75B5?style=for-the-badge" alt="Bubble Tea" />
  <img src="https://img.shields.io/badge/FFmpeg-8.1-007808?style=for-the-badge" alt="FFmpeg" />
  <img src="https://img.shields.io/badge/License-MIT-FFD700?style=for-the-badge" alt="MIT" />
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=for-the-badge" alt="Platform" />
</p>

<p align="center">
  <strong>TIKWA</strong> = <strong>TIK</strong>Tok + <strong>WA</strong> Enhancer<br/>
  TUI premium di terminal — poles video HP (Samsung A23 5G, dll) biar <em>tidak burem</em> pas upload ke TikTok & WhatsApp.
</p>

<p align="center">

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃   ▶  T I K W A  ◀              ┃
┃   Bikin Video Jadi Cling ✨   ┃
┃   TikTok • WhatsApp • Go      ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

</p>

---

## ✨ Preview TUI

```
┌──────────────────────────────────────────────┐
│  ▶  TIKWA  ◀  — Bikin Video HP Tidak Burem   │
│  TikTok • WhatsApp                           │
├──────────────────────────────────────────────┤
│  [1] Untuk TikTok  — 1080×1920 vertikal      │
│  [2] Untuk WhatsApp — 720p kecil (<25MB)     │
│  [q] Keluar                                  │
├──────────────────────────────────────────────┤
│  ↑↓ pilih  •  1/2 atau Enter  •  q keluar   │
└──────────────────────────────────────────────┘
         ↓ pilih 1
┌──────────────────────────────────────────────┐
│  → Mode: TIKTOK (1080×1920)                  │
│  Path video: [/Users/rafie/Movies/VID.mp4 ]  │
│  📭 Contoh: drag & drop file ke terminal     │
└──────────────────────────────────────────────┘
         ↓ Enter
┌──────────────────────────────────────────────┐
│  ⏳ Memoles: VID.mp4 → VID-tiktok.mp4        │
│  [█████████████░░░░░░] 65%                   │
└──────────────────────────────────────────────┘
         ↓
┌──────────────────────────────────────────────┐
│  ✅ Berhasil! Hasil: VID-tiktok.mp4 (42MB)   │
└──────────────────────────────────────────────┘
```

> Empty state `📭`, Loading bar, Error `❌ File bukan video` — semua dipoles, tidak cuma `print`.

---

## 🎯 Fitur v1 (Level 1 — Basic Polish)

| Profil | Resolusi | Bitrate | Hasil | Cocok untuk |
|---|---|---|---|---|
| **TikTok** | 1080×1920 vertikal (pad + scale, tidak stretch) | 5000k + CRF 23 | Tajam, terang (+6% brightness, +5% contrast) | Upload TikTok tidak diperas lagi |
| **WhatsApp** | 720×1280 | 2500k + CRF 23 | Kecil <25MB, tetap jernih | Lolos limit 64MB WA |

- **Aman:** `exec.Command` tanpa shell (anti command injection), `filepath.Clean`, validasi ekstensi & size
- **Cepat:** `preset medium`, single binary Go, tidak butuh Python/Node
- **Ringan:** jalan di MBA M4 sampai laptop kentang

## 📦 Install

### Opsi A — Go install (paling gampang kalau sudah ada Go)
```bash
go install github.com/RafieSA/tikwa@latest
tikwa
```

### Opsi B — Download binary (tanpa Go)
Download di [Releases](https://github.com/RafieSA/tikwa/releases) → pilih sesuai OS:
- `tikwa-darwin-arm64` (M1/M2/M3/M4 Mac)
- `tikwa-darwin-amd64` (Intel Mac)
- `tikwa-linux-amd64`
- `tikwa-windows-amd64.exe`

```bash
chmod +x tikwa-darwin-arm64
./tikwa-darwin-arm64
```

### Opsi C — Build dari source
```bash
git clone https://github.com/RafieSA/tikwa.git
cd tikwa
go build -o tikwa .
./tikwa
```

**Syarat:** FFmpeg wajib ada
```bash
brew install ffmpeg        # macOS
sudo apt install ffmpeg    # Ubuntu/Debian
winget install ffmpeg      # Windows
ffmpeg -version            # cek
```

## 🚀 Pakai

1. Jalankan `tikwa`
2. Tekan `1` untuk TikTok atau `2` untuk WhatsApp (atau `↑↓` + Enter)
3. Paste path video — **drag & drop** file ke terminal juga bisa
4. Enter → tunggu loading → jadi!

Output di folder yang sama:
- `VID_20260915-tiktok.mp4` (TikTok)
- `VID_20260915-whatsapp.mp4` (WA)

Edge cases yang sudah di-handle:
- File kosong / tidak ada → `❌ file tidak ditemukan`
- Salah pilih folder / .jpg → `❌ format belum didukung`
- Nama file ada spasi `video gue.mp4` → aman
- FFmpeg belum install → `❌ ffmpeg tidak ditemukan — brew install ffmpeg`

## 🛠 Teknologi

- **Go 1.26** + **Bubble Tea** + **Lipgloss** + **Bubbles** — TUI modern, single binary
- **FFmpeg 8.1** — tidak bikin dari 0, manfaatkan yang sudah ada (YAGNI, DRY)
- Clean Code, SRP, testable, logging — no spaghetti, no shortcut


## 🤝 Kontribusi

PR & Issue welcome! Jalanin `go vet ./...` dan `go test ./...` sebelum PR.

## 📄 License

MIT — bebas pakai, lihat [LICENSE](LICENSE).
