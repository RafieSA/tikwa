# Contributing ke TIKWA

Terima kasih mau kontribusi! Biar rapi dan tidak berantakan, ikuti langkah ini.

## Cara Kontribusi (Super Presisi)

1. **Fork** repo ini → `RafieSA/tikwa` klik Fork
2. **Clone** fork kamu:
   ```bash
   git clone https://github.com/<username>/tikwa.git
   cd tikwa
   ```
3. **Bikin branch** baru (jangan langsung di `main`):
   ```bash
   git checkout -b feat/nama-fitur
   ```
4. **Ngoding** — patuhi:
   - Clean Code, SRP, DRY, KISS
   - `go vet ./...` harus lulus
   - `go test ./...` harus lulus
   - Jangan bikin dari 0 kalau sudah ada library
5. **Test di local:**
   ```bash
   go vet ./...
   go test ./... -v
   go build -o tikwa .
   ./tikwa  # cek TUI
   ```
6. **Commit** jelas:
   ```bash
   git commit -m "feat: tambah profil Instagram"
   ```
7. **Push & PR:**
   ```bash
   git push origin feat/nama-fitur
   ```
   Lalu buka GitHub → Create Pull Request → jelaskan perubahan + screenshot TUI kalau ada.

## Aturan

- No shortcut, no happy-path saja — handle empty, error, file besar, nama spasi
- No spaghetti, no secret bocor
- Satu PR = satu fitur/fix (jangan campur)
- Tulis deskripsi PR yang jelas: masalah apa, solusi apa, sudah dites apa

## Butuh Bantuan?

Buka Issue → pilih template Bug atau Feature. Tanya saya lebih detail dan mendalam di Issue.
