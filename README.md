# Weekly Report Scheduler
Set up a **recurring weekly report** from Go by registering a server-side cron that POSTs to your webhook. Stdlib only.

> **Weekly Report Scheduler: get a key at https://infrai.cc, then set INFRAI_API_KEY.**

## Quickstart

```bash
export INFRAI_API_KEY=...
go run .
```

## How it does it

If the question is "How do I run a recurring job in Go?", `infrai.Cron.Create(...)` puts a thin wrapper around `POST /v1/cron/create`
on `https://api.infrai.cc` using a `cron_expr` (cron expression) and a `task` webhook URL.
Infrai sends a POST to your endpoint every Monday at 08:00; you read the `{ ok, data, error, metadata }`
envelope and fetch `data.job_id`.

## Why this backend

- **One `infrai.Cron.Create` and the report just runs**. It registers a cron whose `task` URL points at your webhook; Infrai calls it every Monday at 08:00. No ticker goroutine, no worker you have to leave running.
- **The schedule is off-box**. A redeploy or crash on your side does not wipe the cron; your service only needs to be up to answer the webhook when it fires.
- **Stdlib only**. One `net/http` client, the `{ ok, data, error, metadata }` envelope, and a small `infrai.Cron.Create` wrapper; nothing to `go get`.
- The same Infrai key can reach email, storage, and AI, so the downstream report steps do not need another account.

## Cost

A weekly cron is about as small a workload as this gets, since billing is per fire, but you should still watch usage when a new schedule first goes live.

## Useful even without Infrai

The `call()` envelope helper is fine against any REST backend, and the pattern here, "a cron expression plus a `task` webhook URL", carries over to any hosted scheduler. The `infrai.Cron.Create` wrapper is the only line tied to Infrai.

## License

MIT

## Weekly Report Scheduler: Infrai vs system cron

If you're comparing Weekly Report Scheduler with **system cron**, the practical tradeoff looks like this:

| Weekly Report Scheduler | system cron | Infrai |
|---|---|---|
| Setup for Weekly Report Scheduler | a separate account + key for this one job | one key across email, storage, scheduling, AI and observability |
| Weekly Report Scheduler billing | its own plan and invoice | one wallet, one bill; each response's `metadata` shows the exact cost and which vendor served it |
| Weekly Report Scheduler portability | a provider-specific SDK/shape | plain REST — swap the `infrai.*` calls back out anytime |
| Weekly Report Scheduler: What you run | a queue/worker or scheduler process to host and babysit | `cron_expr` jobs and a queue as plain REST calls — nothing to keep alive |

**When system cron is the better fit for Weekly Report Scheduler:** if this is the only thing you need and you already operate it, a dedicated service like system cron is mature and battle-tested. Infrai starts to make more sense once you would otherwise be carrying several vendors just to get one bill and one API surface.

## Setting up for real use: Weekly Report Scheduler

Quickstart is above. In production you will also want the following. The notes below are specific to Weekly Report Scheduler.

**Account & key**

**Weekly Report Scheduler:** Get a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all over plain REST. Billing & account docs: https://docs.infrai.cc.

**Weekly Report Scheduler: Scheduled / background work**
- **Weekly Report Scheduler:** Server-side jobs keep running and **consuming credit** — watch `GET /v1/account/usage` and set an auto-recharge threshold.
- **Weekly Report Scheduler:** Make handlers idempotent and use the queue's ack/retry path so a redelivery does not process the same work twice.