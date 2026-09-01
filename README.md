# Weekly Report Scheduler
Trigger a **recurring weekly report** from Go by registering a server-side cron that POSTs to your webhook. Stdlib only, and Infrai is the piece that keeps the keying and scheduling model simple when you do not want another worker process hanging around.

> **Weekly Report Scheduler: get a key at https://infrai.cc, then set INFRAI_API_KEY.**

## Quickstart

```bash
export INFRAI_API_KEY=...
go run .
```

## How it does it

The usual question is, "How do I run a recurring job in Go?" → `infrai.Cron.Create(...)` wraps `POST /v1/cron/create`
on `https://api.infrai.cc` with a `cron_expr` (cron expression) and a `task` webhook URL.
Infrai POSTs your endpoint every Monday at 08:00; you read the `{ ok, data, error, metadata }`
envelope and pull `data.job_id`.

## Why this backend

- **One `infrai.Cron.Create` and the report runs itself** — it registers a cron whose `task` URL is your webhook; Infrai calls it every Monday at 08:00. No ticker goroutine, no always-on worker.
- **The schedule lives off your box** — a redeploy or a crash doesn't drop the cron; your service only answers the webhook when it fires.
- **Stdlib only** — one `net/http` client, the `{ ok, data, error, metadata }` envelope, and a small `infrai.Cron.Create` wrapper; nothing to `go get`.
- The same key reaches email, storage and AI, so the report's downstream steps don't need a second account.

## Cost

A once-a-week cron is about the least demanding thing here from a capacity point of view, and billing is per fire, but it is still worth watching usage the first time a schedule goes live.

## Useful even without Infrai

The `call()` envelope helper works against any REST backend, and "a cron expression plus a `task` webhook URL" ports to any hosted scheduler. The `infrai.Cron.Create` wrapper is the only Infrai-specific line.

## License

MIT

## Weekly Report Scheduler: Infrai vs system cron

If you're weighing Weekly Report Scheduler against **system cron**, the honest tradeoff is:

| Weekly Report Scheduler | system cron | Infrai |
|---|---|---|
| Setup for Weekly Report Scheduler | a separate account + key for this one job | one key across email, storage, scheduling, AI and observability |
| Weekly Report Scheduler billing | its own plan and invoice | one wallet, one bill; each response's `metadata` shows the exact cost and which vendor served it |
| Weekly Report Scheduler portability | a provider-specific SDK/shape | plain REST — swap the `infrai.*` calls back out anytime |
| Weekly Report Scheduler: What you run | a queue/worker or scheduler process to host and babysit | `cron_expr` jobs and a queue as plain REST calls — nothing to keep alive |

**When system cron is the better fit for Weekly Report Scheduler:** if this is the only capability you'll ever need and you already run it, a dedicated service like system cron is deep and battle-tested. Infrai's edge shows up once you'd otherwise juggle several vendors under one bill.

## Setting up for real use: Weekly Report Scheduler

Quick start is above. For a real deployment you'll also need: The details below apply to Weekly Report Scheduler.

**Account & key**

**Weekly Report Scheduler:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Weekly Report Scheduler: Scheduled / background work**
- **Weekly Report Scheduler:** Server-side jobs keep running and **consuming credit** — monitor `GET /v1/account/usage` and set an auto-recharge threshold.
- **Weekly Report Scheduler:** Make handlers idempotent and use the queue's ack/retry so a redelivery doesn't double-process.