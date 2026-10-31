# RUPI-108 · Onboarding lengkap: zona waktu dan saldo awal

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | S | FR-02 |

## Latar
Onboarding saat ini hanya menanyakan nama; zona waktu selalu WIB dan saldo awal
belum ada.

## Cakupan
- Langkah pilih zona waktu dengan tombol WIB / WITA / WIT (default WIB).
- Langkah saldo awal opsional dengan tombol "Lewati"; disimpan sebagai pemasukan
  kategori Lainnya berdeskripsi "Saldo awal".
- Perintah `/zonawaktu` untuk mengubah zona waktu setelah onboarding.

## Kriteria penerimaan
- [ ] Pengguna WITA yang mencatat "kemarin" mendapat tanggal kemarin menurut WITA.
- [ ] Saldo awal bisa dilewati tanpa membuat transaksi.
- [ ] `/zonawaktu` memperbarui `user_svc.users.timezone`.

## Dependensi
RUPI-101.
