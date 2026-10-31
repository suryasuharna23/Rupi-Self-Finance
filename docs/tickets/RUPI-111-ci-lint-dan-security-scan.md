# RUPI-111 · CI: golangci-lint, batas cakupan, dan pemindaian Trivy

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | S | NFR-12, SEC-12 |

## Cakupan
- Tambah `golangci-lint` dengan konfigurasi di repo.
- Gagalkan CI bila cakupan `internal/parser` atau `internal/transaction` < 70%.
- Pindai dependensi dan image dengan Trivy.
- Build image di CI dengan tag SHA commit (bukan `latest`).

## Kriteria penerimaan
- [ ] PR dengan error lint atau cakupan < 70% gagal di CI.
- [ ] Temuan Trivy berkategori HIGH/CRITICAL menggagalkan CI.
- [ ] Image ter-push dengan tag SHA.

## Dependensi
Tidak ada.
