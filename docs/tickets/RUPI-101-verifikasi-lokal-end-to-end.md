# RUPI-101 · Verifikasi bot berjalan end-to-end secara lokal

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | S | FR-01 – FR-16 |

## Latar
Setup Fase 1 baru diverifikasi sampai tahap compile dan unit test. Kode database
(`internal/transaction`, `user`, `report`, `conversation/store.go`) dan alur
Telegram belum pernah dijalankan dengan PostgreSQL dan bot sungguhan.

## Cakupan
- Pasang Docker Desktop, jalankan `make db`, isi `.env`, jalankan `make run`.
- Uji manual seluruh alur di Telegram dan perbaiki bug yang muncul.
- Dokumentasikan langkah khusus Windows/PowerShell di README.

## Kriteria penerimaan
- [ ] Migrasi berjalan bersih di database kosong dan aman dijalankan ulang.
- [ ] `/start` → kebijakan privasi → tanya nama → onboarding selesai.
- [ ] "beli bakso 15rb" tersimpan dan dibalas dengan tombol Ubah / Batalkan.
- [ ] "bakso 15rb sama es teh goceng" menghasilkan tepat 2 transaksi.
- [ ] "makan siang" → "Berapa nominalnya?" → "15rb" menyimpan transaksi makan siang.
- [ ] "yang bakso tadi 18rb" mengubah nominal dan menulis baris `transaction_edits`.
- [ ] Tombol Batalkan, `/undo`, `/ringkas`, `/hapusakun` bekerja.
- [ ] "bulan ini habis berapa buat makan?" dijawab dengan total yang benar.
- [ ] `GET /healthz` mengembalikan 200 saat database hidup dan 503 saat mati.

## Dependensi
Tidak ada. Kerjakan paling awal.
