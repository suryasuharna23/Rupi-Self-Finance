# RUPI-211 · End-to-end test di kind (gerbang MVP)

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | M | Pengujian: End-to-end, gerbang MVP |

## Cakupan
- Skrip Telegram tiruan yang mengirim webhook ke Channel Adapter dan menangkap balasan.
- LLM di-mock dengan respons tetap.
- Jalankan semua kriteria penerimaan P0 dari dokumen.

## Kriteria penerimaan
- [ ] "bakso 15rb sama es teh goceng" menghasilkan tepat 2 transaksi Rp15.000 dan Rp5.000.
- [ ] Webhook yang sama dikirim dua kali hanya menghasilkan satu set transaksi.
- [ ] "yang bakso tadi 18rb" mengubah transaksi bakso terakhir dan mencatat jejak di `transaction_edits`.
- [ ] Saat LLM tidak merespons dalam 5 detik, pengguna menerima pertanyaan klarifikasi, bukan error.
- [ ] Pengguna A tidak dapat melihat atau mengubah transaksi pengguna B melalui RPC apa pun.
- [ ] Test berjalan di CI.

## Dependensi
RUPI-209, RUPI-210.
