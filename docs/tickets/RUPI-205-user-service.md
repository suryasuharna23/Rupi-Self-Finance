# RUPI-205 · Pecah User Service

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | S | `UserService`, FR-01 – FR-03 |

## Cakupan
- `cmd/user-service` dengan RPC `GetOrCreateUser`, `UpdateProfile`, `DeleteUser`.
- Hanya service ini yang mengakses schema `user_svc`.
- Penghapusan akun mengoordinasikan penghapusan data di Transaction Service.

## Kriteria penerimaan
- [ ] `GetOrCreateUser` aman dipanggil bersamaan untuk ID Telegram yang sama.
- [ ] `/hapusakun` tetap menghapus semua data pengguna setelah dipecah.

## Dependensi
RUPI-201, RUPI-202.
