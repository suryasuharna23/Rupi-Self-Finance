# RUPI-309 · Advanced queries: count, largest, and breakdown

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 3 · Events, budgets, reports | P0 | M | FR-20, `ReportQuery.Metric` |

## Background
Only the SUM metric is currently supported. The `ReportQuery` in the requirements document also
defines COUNT, TOP_N, BREAKDOWN, and GroupBy.

## Scope
- The parser recognizes: "berapa kali jajan kopi minggu ini?" (COUNT),
  "pengeluaran terbesar bulan ini apa?" (TOP_N),
  "rincian pengeluaran bulan ini per kategori" (BREAKDOWN + GroupBy).
- The Report Service implements each metric.
- Additional time ranges: "kemarin", "minggu lalu", "dari tgl 1 sampai 15".

## Acceptance criteria
- [ ] Each metric has unit tests and examples in the evaluation dataset.
- [ ] BREAKDOWN replies are sorted from the largest category down.

## Dependencies
RUPI-206.
