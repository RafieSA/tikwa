# Security Policy

## Versi yang didukung

| Versi | Didukung |
|---|:---:|
| v0.1.x | ✅ |

## Lapor Vulnerability

Jangan buka Issue publik untuk security. Hubungi:

- Email: **rafiesafaraz@student.telkomuniversity.ac.id**
- Subjek: `[SECURITY] TIKWA - deskripsi singkat`

Kita akan balas dalam 48 jam dan fix dalam 7 hari kalau valid.

## Yang kita jaga

- No command injection (pakai `exec.Command` tanpa shell, sudah di-test)
- Validasi path (`filepath.Clean`, cek ekstensi & size)
- No secret bocor di log

Terima kasih sudah bantu jaga TIKWA tetap aman.
