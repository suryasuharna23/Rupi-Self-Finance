# RUPI-403 · Load test 50 pesan/detik

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 4 · Siap produksi | P1 | M | NFR-01, NFR-06, gerbang Fase 4 |

## Cakupan
- Skrip k6 (atau ghz) yang mengirim webhook tiruan dari 1.000 pengguna simulasi.
- 50 pesan/detik selama 10 menit, campuran kalimat jalur aturan dan LLM (di-mock).
- Laporan hasil disimpan di `docs/`.

## Kriteria penerimaan
- [ ] p95 jalur aturan < 500 ms.
- [ ] Tidak ada pesan hilang dan tidak ada transaksi ganda (dicocokkan dengan database).
- [ ] HPA Parser Service naik dan turun selama tes.

## Dependensi
RUPI-401, RUPI-402.
