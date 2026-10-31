# Rupi

Asisten pencatatan keuangan pribadi untuk mahasiswa lewat chat Telegram. Ketik
"beli bakso 15rb" dan Rupi mencatat, mengategorikan, dan melaporkannya.

Kebutuhan lengkap ada di [`docs/Rupi - Dokumen Kebutuhan Sistem.pdf`](docs/).

## Status: Fase 1, inti percakapan (monolit)

Sesuai roadmap, Fase 1 sengaja berupa satu service Go agar kualitas parser bisa
diasah sebelum sistem dipecah. Paket di `internal/` sudah mengikuti batas
service Fase 2, jadi tiap paket nanti menjadi satu service gRPC.

| Paket                    | Calon service (Fase 2) | Isi                                                       |
| ------------------------ | ---------------------- | --------------------------------------------------------- |
| `internal/telegram`      | Channel Adapter        | Long polling / webhook dengan `secret_token`              |
| `internal/conversation`  | Conversation Service   | Onboarding, klarifikasi, koreksi, idempotensi, rate limit |
| `internal/parser`        | Parser Service         | Aturan (nominal, waktu, kategori) + antarmuka LLM         |
| `internal/transaction`   | Transaction Service    | Simpan, koreksi + jejak audit, undo                       |
| `internal/user`          | User Service           | Profil, zona waktu, hapus akun                            |
| `internal/report`        | Report Service         | Total per periode/kategori, `/ringkas`                    |
| `internal/storage`       | —                      | Koneksi PostgreSQL dan migrasi (satu schema per service)  |

## Menjalankan secara lokal

Prasyarat: Go 1.27+, Docker, dan token bot dari [@BotFather](https://t.me/BotFather).

```sh
cp .env.example .env        # isi TELEGRAM_BOT_TOKEN
make db                     # PostgreSQL 16 lewat docker compose
make run                    # migrasi otomatis lalu mulai long polling
```

Tanpa `make` di PowerShell:

```powershell
docker compose up -d postgres
Get-Content .env | Where-Object { $_ -match "^[A-Z_]+=" } | ForEach-Object {
  $k, $v = $_ -split "=", 2; Set-Item "env:$k" $v
}
go run ./cmd/rupi
```

## Pengujian

```sh
make test   # unit test
make e2e    # alur chat lengkap: Postgres sungguhan + Telegram tiruan
make eval   # akurasi parser terhadap testdata/eval/dataset.jsonl
```

Gerbang Fase 1 adalah akurasi ekstraksi ≥ 95% pada dataset evaluasi 300 kalimat
(NFR-03). Dataset baru berisi 20 contoh dari dokumen kebutuhan; tambahkan
kalimat nyata ke `testdata/eval/dataset.jsonl`.

## Perintah bot

| Perintah                            | Fungsi                                  |
| ----------------------------------- | --------------------------------------- |
| `/start`                            | Daftar dan onboarding                   |
| `beli bakso 15rb`                   | Catat transaksi                         |
| `yang bakso tadi 18rb`              | Koreksi nominal                         |
| `/undo`, `batalin yang terakhir`    | Batalkan transaksi terakhir (≤ 24 jam)  |
| `bulan ini habis berapa buat makan?`| Tanya total                             |
| `/ringkas`                          | Ringkasan hari, minggu, dan bulan ini   |
| `/hapusakun`                        | Hapus semua data                        |

## Belum dikerjakan

- Implementasi LLM: penyedianya masih pertanyaan terbuka di dokumen. Antarmuka
  `parser.LLM` dan validasi keluarannya (SEC-05) sudah ada; sementara ini
  pesan yang tidak dipahami aturan dijawab dengan pertanyaan balik.
- Saldo awal saat onboarding (opsional di FR-02) dan pengubahan zona waktu.
- Purge akun yang sudah di-soft delete setelah 30 hari.
- Integration test dengan Testcontainers.
