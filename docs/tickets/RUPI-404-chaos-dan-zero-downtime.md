# RUPI-404 · Chaos sederhana dan rolling update tanpa downtime

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 4 · Siap produksi | P1 | M | NFR-07, NFR-08, NFR-10 |

## Cakupan
- Graceful shutdown di semua service: berhenti menerima request, selesaikan
  yang berjalan, maksimal 30 detik.
- Skrip `kubectl delete pod` acak terjadwal selama load test.
- Rolling update Parser Service di tengah load test.

## Kriteria penerimaan
- [ ] Tidak ada transaksi hilang atau ganda selama chaos.
- [ ] Rolling update Parser Service selama load test tidak menyebabkan pesan hilang.

## Dependensi
RUPI-403.
