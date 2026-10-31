# RUPI-308 · Custom categories and keyword rules

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 3 · Events, budgets, reports | P1 | M | FR-17 |

## Scope
- Sentences such as "bikin kategori Kopi" and "kalau ada kata 'kopken' masukkan ke Jajan".
- Store them in `transaction_svc.categories` (the `user_id` and `keywords` columns already exist).
- The parser receives the user's list of custom categories and gives their keywords
  priority over the built-in dictionary.
- A `/kategori` command to view and delete custom categories.

## Acceptance criteria
- [ ] Once the rule is created, "kopken 25rb" goes into that custom category.
- [ ] User A's custom categories are not visible to user B.
- [ ] Built-in category names cannot be overridden.

## Dependencies
RUPI-203, RUPI-204.
