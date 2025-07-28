# SENTINEL

SENTINEL adalah sistem monitoring sederhana yang ditulis dalam bahasa Go. Aplikasi ini dapat memantau status berbagai layanan web melalui HTTP dan melaporkan statusnya secara berkala.

## Fitur

- Monitoring berbagai layanan web melalui HTTP
- Konfigurasi sederhana menggunakan format YAML
- Pelaporan status UP/DOWN dengan waktu respons
- Pemeriksaan otomatis dengan interval yang dapat dikonfigurasi
- Konkurensi untuk pemeriksaan yang efisien
- CLI yang fleksibel dengan berbagai perintah

## Cara Menggunakan

1. Pastikan Go sudah terinstal di sistem Anda
2. Clone repositori ini
3. Sesuaikan file `sentinel.yaml` dengan layanan yang ingin Anda pantau
4. Jalankan aplikasi dengan salah satu perintah berikut:
   - `go run main.go run` - Mulai monitoring secara terus-menerus
   - `go run main.go once` - Jalankan pengecekan sekali saja
   - `go run main.go validate` - Validasi file konfigurasi

## Struktur Konfigurasi

File konfigurasi `sentinel.yaml` memiliki format sebagai berikut:

```yaml
# Global settings
settings:
  check_interval: 1m        # Interval antara pengecekan (e.g., 30s, 1m, 5m)
  default_timeout: 5s       # Default timeout untuk semua layanan

# Daftar layanan yang akan dipantau
services:
  - name: "Google"
    url: "https://www.google.com"
  
  - name: "GitHub"
    url: "https://github.com"
  
  - name: "Example"
    url: "https://example.com"
    timeout: 10s            # Custom timeout untuk layanan ini
```

## Struktur Proyek

- `main.go` - File utama program
- `config/config.go` - Package untuk mengelola konfigurasi
- `checker/checker.go` - Package untuk melakukan pengecekan layanan
- `sentinel.yaml` - File konfigurasi

## Roadmap

- [x] Fase 1: Core Engine - Implementasi dasar
- [x] Fase 2: Making it "Pintar" - Konkurensi dan CLI yang lebih baik
- [ ] Fase 3: Advanced Features - Notifikasi dan tipe pengecekan lanjutan
- [ ] Fase 4: Open Source Readiness - Dokumentasi dan pengujian

## Lisensi

MIT
