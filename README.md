# SENTINEL

SENTINEL adalah sistem monitoring sederhana yang ditulis dalam bahasa Go. Aplikasi ini dapat memantau status berbagai layanan web melalui HTTP dan melaporkan statusnya secara berkala.

## Fitur

- Monitoring berbagai layanan web melalui HTTP
- Konfigurasi sederhana menggunakan format YAML
- Pelaporan status UP/DOWN dengan waktu respons
- Pemeriksaan otomatis setiap menit

## Cara Menggunakan

1. Pastikan Go sudah terinstal di sistem Anda
2. Clone repositori ini
3. Sesuaikan file `sentinel.yaml` dengan layanan yang ingin Anda pantau
4. Jalankan aplikasi dengan perintah `go run main.go`

## Struktur Konfigurasi

File konfigurasi `sentinel.yaml` memiliki format sebagai berikut:

```yaml
services:
  - name: "Google"
    url: "https://www.google.com"
  - name: "GitHub"
    url: "https://github.com"
  - name: "Example"
    url: "https://example.com"
```

## Struktur Proyek

- `main.go` - File utama program
- `config/config.go` - Package untuk mengelola konfigurasi
- `checker/checker.go` - Package untuk melakukan pengecekan layanan
- `sentinel.yaml` - File konfigurasi

## Roadmap

- [x] Fase 1: Core Engine - Implementasi dasar
- [ ] Fase 2: Making it "Pintar" - Konkurensi dan CLI yang lebih baik
- [ ] Fase 3: Advanced Features - Notifikasi dan tipe pengecekan lanjutan
- [ ] Fase 4: Open Source Readiness - Dokumentasi dan pengujian

## Lisensi

MIT
