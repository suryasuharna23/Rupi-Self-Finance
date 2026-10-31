# RUPI-402 · HPA, PodDisruptionBudget, dan NetworkPolicy

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 4 · Siap produksi | P1 | M | NFR-05, NFR-06, Ketentuan wajib Deployment |

## Cakupan
- HPA Parser Service 2–10 Pod, target CPU 70%.
- `PodDisruptionBudget` dengan `minAvailable: 1` untuk setiap Deployment.
- `NetworkPolicy` default deny; hanya jalur panggilan di diagram arsitektur yang diizinkan.
- Semua container berjalan non-root dengan filesystem read-only (SEC-12).

## Kriteria penerimaan
- [ ] Panggilan yang tidak ada di diagram (misalnya Channel Adapter → PostgreSQL) diblokir.
- [ ] `kubectl drain` sebuah node tidak membuat service mana pun turun ke 0 replika.
- [ ] HPA naik saat beban dan turun kembali setelahnya.

## Dependensi
RUPI-209.
