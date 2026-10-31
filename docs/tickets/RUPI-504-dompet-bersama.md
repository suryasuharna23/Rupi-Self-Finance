# RUPI-504 · Dompet bersama

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 5 · Fitur lanjutan | P2 | L | FR-40, FR-41, FR-42 |

## Cakupan
- Buat dompet bersama dan undang anggota lewat tautan (tabel `wallets`, `wallet_members`).
- Transaksi dompet bersama mencatat siapa yang membayar dan pembagiannya.
- Hitung siapa berutang ke siapa (penyederhanaan utang) dan kirim pengingat.

## Kriteria penerimaan
- [ ] Anggota hanya melihat transaksi di dompet bersama tempat ia terdaftar.
- [ ] Total utang seluruh anggota selalu berjumlah nol.
- [ ] Pengingat tagihan bisa dimatikan per anggota.

## Dependensi
Gerbang Fase 4.
