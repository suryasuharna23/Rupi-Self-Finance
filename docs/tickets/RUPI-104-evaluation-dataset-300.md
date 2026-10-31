# RUPI-104 · 300-sentence parser evaluation dataset

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | L | NFR-03, NFR-04, Phase 1 gate |

## Background
The Phase 1 gate is ≥ 95% extraction accuracy and ≥ 85% category accuracy on a labeled
dataset of at least 300 sentences. `testdata/eval/dataset.jsonl` currently contains only 20.

## Scope
- Collect real sentences from prospective test users (an open question in the requirements document).
- Spread examples evenly: all amount formats, relative times, income, multi-transaction,
  corrections, queries, small talk, typos, and deliberately ambiguous sentences.
- Add date labels to the dataset and check dates in `cmd/eval`.
- Add scores per error type (amount, type, date, category).
- Fix the parser rules until the targets are met.

## Acceptance criteria
- [ ] The dataset has ≥ 300 sentences; each built-in category has ≥ 10 examples.
- [ ] `go run ./cmd/eval` reports extraction ≥ 95%, category ≥ 85%, rule path ≥ 60%.
- [ ] The report includes a breakdown per error type.
- [ ] The dataset contains no personal data (real names, bank account numbers).

## Dependencies
None. Can be worked on in parallel with other tickets.
