# Tiket Rupi

Tiket disusun dari [dokumen kebutuhan sistem](../Rupi%20-%20Dokumen%20Kebutuhan%20Sistem.pdf)
dan mengikuti roadmap lima fase. Setiap fase baru dimulai setelah gerbang fase
sebelumnya terpenuhi. Ukuran: **S** ≈ 1–2 hari, **M** ≈ 3–5 hari, **L** ≈ 1–2 minggu.

**Bisa dikerjakan sekarang** (tanpa dependensi): RUPI-101, RUPI-103, RUPI-104, RUPI-107, RUPI-111.
RUPI-105 bisa dimulai setelah penyedia LLM diputuskan.

## Fase 1 · Inti percakapan (monolit)
Gerbang: akurasi parser ≥ 95% pada dataset evaluasi 300 kalimat (NFR-03).

| Tiket | Judul | Prioritas | Ukuran | Dependensi |
| --- | --- | --- | --- | --- |
| [RUPI-101](RUPI-101-verifikasi-lokal-end-to-end.md) | Verifikasi bot berjalan end-to-end secara lokal | P0 | S | — |
| [RUPI-102](RUPI-102-integration-test-testcontainers.md) | Integration test store dengan Testcontainers | P0 | M | 101 |
| [RUPI-103](RUPI-103-unit-test-conversation.md) | Unit test Conversation Service dengan fake | P0 | M | — |
| [RUPI-104](RUPI-104-dataset-evaluasi-300.md) | Dataset evaluasi parser 300 kalimat | P0 | L | — |
| [RUPI-105](RUPI-105-implementasi-llm-fallback.md) | Implementasi LLM sebagai fallback parser | P0 | L | keputusan LLM |
| [RUPI-106](RUPI-106-batas-dan-biaya-llm.md) | Batas harian dan pencatatan biaya LLM | P0 | S | 105 |
| [RUPI-107](RUPI-107-parser-logs.md) | Tabel parser_logs dengan teks disamarkan | P0 | M | — |
| [RUPI-108](RUPI-108-onboarding-lengkap.md) | Onboarding lengkap: zona waktu dan saldo awal | P0 | S | 101 |
| [RUPI-109](RUPI-109-purge-akun.md) | Purge permanen akun setelah 30 hari | P0 | S | 102 |
| [RUPI-110](RUPI-110-koreksi-kategori-deskripsi.md) | Koreksi kategori, deskripsi, dan tanggal | P0 | M | 101 |
| [RUPI-111](RUPI-111-ci-lint-dan-security-scan.md) | CI: golangci-lint, batas cakupan, Trivy | P0 | S | — |

## Fase 2 · Pecah ke microservices gRPC
Gerbang MVP: semua kriteria penerimaan P0 lolos end-to-end di kind.

| Tiket | Judul | Prioritas | Ukuran | Dependensi |
| --- | --- | --- | --- | --- |
| [RUPI-201](RUPI-201-kontrak-proto-dan-buf.md) | Kontrak protobuf dan buf | P0 | M | gerbang Fase 1 |
| [RUPI-202](RUPI-202-grpc-interceptor-dan-konvensi.md) | Interceptor gRPC dan konvensi error | P0 | M | 201 |
| [RUPI-203](RUPI-203-parser-service.md) | Pecah Parser Service | P0 | M | 201, 202 |
| [RUPI-204](RUPI-204-transaction-service.md) | Pecah Transaction Service | P0 | M | 201, 202, 102 |
| [RUPI-205](RUPI-205-user-service.md) | Pecah User Service | P0 | S | 201, 202 |
| [RUPI-206](RUPI-206-report-service.md) | Pecah Report Service | P0 | M | 204 |
| [RUPI-207](RUPI-207-conversation-service-redis.md) | Conversation Service dengan konteks di Redis | P0 | L | 203–206 |
| [RUPI-208](RUPI-208-channel-adapter.md) | Channel Adapter sebagai service terpisah | P0 | M | 207 |
| [RUPI-209](RUPI-209-helm-chart-kind.md) | Helm chart dan deploy ke kind | P0 | L | 203–208 |
| [RUPI-210](RUPI-210-mtls-antar-service.md) | mTLS antar service | P0 | M | 209 |
| [RUPI-211](RUPI-211-e2e-test-kind.md) | End-to-end test di kind (gerbang MVP) | P0 | M | 209, 210 |

## Fase 3 · Event, budget, dan laporan
Gerbang: peringatan budget terkirim tepat sekali per periode.

| Tiket | Judul | Prioritas | Ukuran | Dependensi |
| --- | --- | --- | --- | --- |
| [RUPI-301](RUPI-301-nats-jetstream-event.md) | Event transaksi lewat NATS JetStream | P1 | M | 211 |
| [RUPI-302](RUPI-302-budget-service.md) | Budget Service dan pengaturan budget lewat chat | P1 | M | 301 |
| [RUPI-303](RUPI-303-peringatan-budget.md) | Peringatan budget 80% dan 100% | P1 | M | 302 |
| [RUPI-304](RUPI-304-perkiraan-sisa-harian.md) | Perkiraan sisa uang harian | P1 | S | 302 |
| [RUPI-305](RUPI-305-pengingat-harian.md) | Pengingat harian opsional | P1 | S | 306 |
| [RUPI-306](RUPI-306-laporan-berkala-cronjob.md) | Laporan mingguan dan bulanan otomatis | P1 | M | 206, 207 |
| [RUPI-307](RUPI-307-perbandingan-periode.md) | Perbandingan dengan periode sebelumnya | P1 | S | 306 |
| [RUPI-308](RUPI-308-kategori-kustom.md) | Kategori kustom dan aturan kata kunci | P1 | M | 203, 204 |
| [RUPI-309](RUPI-309-query-lanjutan.md) | Query lanjutan: jumlah, terbesar, rincian | P0 | M | 206 |

## Fase 4 · Siap produksi
Gerbang: lolos load test tanpa pesan hilang atau ganda (NFR-06, NFR-07).

| Tiket | Judul | Prioritas | Ukuran | Dependensi |
| --- | --- | --- | --- | --- |
| [RUPI-401](RUPI-401-observability.md) | Observability: metrik, log, dan tracing | P1 | L | 209 |
| [RUPI-402](RUPI-402-hpa-pdb-networkpolicy.md) | HPA, PodDisruptionBudget, NetworkPolicy | P1 | M | 209 |
| [RUPI-403](RUPI-403-load-test.md) | Load test 50 pesan/detik | P1 | M | 401, 402 |
| [RUPI-404](RUPI-404-chaos-dan-zero-downtime.md) | Chaos sederhana dan rolling update | P1 | M | 403 |
| [RUPI-405](RUPI-405-backup-dan-enkripsi.md) | Backup harian dan enkripsi at rest | P1 | M | 209 |

## Fase 5 · Fitur lanjutan
Urutan dipilih berdasarkan masukan pengguna uji.

| Tiket | Judul | Prioritas | Ukuran | Dependensi |
| --- | --- | --- | --- | --- |
| [RUPI-501](RUPI-501-web-chat-dan-dashboard.md) | Web chat dan dashboard grafik | P2 | L | gerbang Fase 4 |
| [RUPI-502](RUPI-502-foto-struk-ocr.md) | Foto struk → OCR → transaksi per item | P2 | L | gerbang Fase 4 |
| [RUPI-503](RUPI-503-kanal-whatsapp.md) | Kanal WhatsApp | P2 | L | gerbang Fase 4 |
| [RUPI-504](RUPI-504-dompet-bersama.md) | Dompet bersama | P2 | L | gerbang Fase 4 |
| [RUPI-505](RUPI-505-ekspor-csv.md) | Ekspor transaksi ke CSV | P2 | S | 206 |
