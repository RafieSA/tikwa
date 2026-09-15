# ▶ TIKWA — Bikin Video HP Tidak Burem

**TIKWA** = **TIK**Tok + **WA** Enhancer. TUI premium di terminal untuk memoles video dari HP (misal Samsung A23 5G) agar tidak burem saat upload ke TikTok & WhatsApp.

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃  ▶  T I K W A  ◀         ┃
┃  Bikin Video Jadi Cling ┃
┃  TikTok • WhatsApp      ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

## Fitur v1 (Level 1 Basic Polish)

- **Profil TikTok:** 1080×1920 vertikal, tajam, terang — tidak diperas lagi oleh TikTok
- **Profil WhatsApp:** 720p kecil <25MB, tetap jernih — lolos limit WA
- **TUI cantik:** welcoming animation, menu angka [1]/[2], loading bar, empty & error state jelas
- **Aman:** tanpa shell injection, validasi file, cek FFmpeg

## Install

```bash
git clone https://github.com/rafiesafarazaribowo/tikwa.git
cd tikwa
go build -o tikwa .
./tikwa
```

Butuh **FFmpeg**: `brew install ffmpeg` (macOS) / `sudo apt install ffmpeg` (Linux)

## Pakai

1. Jalankan `./tikwa`
2. Tekan `1` untuk TikTok, `2` untuk WhatsApp
3. Paste path video (drag & drop juga bisa)
4. Enter → tunggu loading → jadi!

Output: `video-tiktok.mp4` atau `video-wa.mp4` di folder yang sama.

## Teknologi

Go + Bubble Tea + Lipgloss + FFmpeg (tidak bikin dari 0, manfaatkan yang sudah ada).

## Roadmap

- v1: Basic Polish (sekarang)
- v2: Smart Clean (denoise + stabilizer)
- v3: AI Upscale (opsional, butuh GPU)

## License

MIT — bebas pakai.
