# RUPI-209 · Helm chart dan deploy ke kind

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | L | bagian Infrastruktur Kubernetes |

## Cakupan
- Cluster lokal kind dengan namespace `rupi-dev` dan `rupi-staging`.
- Helm chart untuk semua service dengan replika dan request/limit sesuai tabel di dokumen.
- PostgreSQL dan Redis sebagai StatefulSet + PVC (atau CloudNativePG).
- Migrasi skema sebagai Job Helm pre-upgrade hook.
- `readinessProbe` dan `livenessProbe` memakai `grpc.health.v1`.
- Konfigurasi lewat ConfigMap; token bot, API key LLM, kredensial DB lewat Secret.
- Load balancing gRPC lewat headless Service + client-side round robin.
- Satu Dockerfile multi-target atau satu Dockerfile per service; image non-root.

## Kriteria penerimaan
- [ ] `helm install` di kind menjalankan semua service sampai Ready.
- [ ] Tidak ada rahasia di values file atau image.
- [ ] Beban gRPC tersebar ke kedua replika (terlihat di log).

## Dependensi
RUPI-203 – RUPI-208.
