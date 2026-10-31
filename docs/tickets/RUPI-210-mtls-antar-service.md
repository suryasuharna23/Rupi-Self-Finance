# RUPI-210 · mTLS antar service

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | M | SEC-03 |

## Cakupan
- Pilih antara cert-manager (mTLS di kode) atau Linkerd (mesh). Linkerd juga
  menyelesaikan load balancing HTTP/2.
- Semua komunikasi antar service terenkripsi dan terautentikasi.

## Kriteria penerimaan
- [ ] Pod tanpa sertifikat valid tidak bisa memanggil service internal.
- [ ] Rotasi sertifikat tidak memutus layanan.

## Dependensi
RUPI-209.
