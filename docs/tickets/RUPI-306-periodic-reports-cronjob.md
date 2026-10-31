# RUPI-306 · Automatic weekly and monthly reports

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 3 · Events, budgets, reports | P1 | M | FR-22, Infrastructure (CronJob) |

## Scope
- Kubernetes CronJobs: Monday 08.00 (weekly report) and the 1st of the month (monthly report), WIB.
- The job calls the Report Service for each active user and then sends the report via the Conversation Service.
- The report contains total spending/income, the top 3 categories, and the largest transaction.

## Acceptance criteria
- [ ] Re-running the job does not send duplicate reports.
- [ ] Users with no transactions in the period do not receive an empty report.
- [ ] The job completes within the 200m / 256Mi resource limit for 1,000 simulated users.

## Dependencies
RUPI-206, RUPI-207.
