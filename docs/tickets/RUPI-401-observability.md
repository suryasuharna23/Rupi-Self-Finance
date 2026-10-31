# RUPI-401 · Observability: metrik, log, dan tracing

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 4 · Siap produksi | P1 | L | NFR-09, NFR-11 |

## Cakupan
- Namespace `monitoring` dengan Prometheus, Grafana, Loki, dan OpenTelemetry Collector.
- Instrumentasi OTel di interceptor gRPC; trace ID sama dengan `x-request-id`.
- Metrik: latensi per RPC, rasio error, porsi jalur aturan vs LLM, token dan biaya LLM.
- Dashboard Grafana: kesehatan service, biaya LLM, akurasi parser (dari `parser_logs`).
- Alert dasar: error rate tinggi, p95 latensi melewati NFR-01/NFR-02.

## Kriteria penerimaan
- [ ] Satu pesan bisa ditelusuri dari Channel Adapter sampai PostgreSQL dalam satu trace.
- [ ] Dashboard menampilkan p95 latensi jalur aturan dan jalur LLM.
- [ ] Dashboard menampilkan biaya LLM per pengguna aktif.

## Dependensi
RUPI-209.
