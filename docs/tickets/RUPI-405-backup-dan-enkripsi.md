# RUPI-405 · Backup harian dan enkripsi at rest

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 4 · Siap produksi | P1 | M | SEC-09, SEC-10 |

## Cakupan
- Backup PostgreSQL harian dengan retensi 7 hari (CloudNativePG atau CronJob `pg_dump`).
- Volume database dan backup terenkripsi.
- Prosedur restore terdokumentasi dan pernah diuji.

## Kriteria penerimaan
- [ ] Restore dari backup kemarin berhasil di namespace terpisah.
- [ ] Backup berumur > 7 hari terhapus otomatis, sehingga data akun yang dihapus
      hilang dari backup dalam 30 hari (SEC-10).

## Dependensi
RUPI-209.
