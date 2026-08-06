# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Detailed project conventions and an Encore framework reference live in [context/project-conventions.md](context/project-conventions.md) — read it when working on anything non-trivial. When a new architectural rule or repeated pattern is agreed with the user, add it there.

## What this is

A B2B/B2C car rental brokerage that aggregates inventory from external suppliers (Flex, Hertz). Two monorepo halves:

- [backend/](backend/) — Go on the [Encore](https://encore.dev) framework (app id `ai-booking-travel-bo22`). Postgres via `sqlc`.
- [frontend/](frontend/) — Next.js App Router + React Query over a generated Encore client, next-intl (Hebrew/English), Shadcn/Tailwind, with **Payload CMS embedded in the same Next app** (see [frontend/CMS/](frontend/CMS/) and the `app/(payload)` route group, backed by its own Postgres — Neon in deployed envs).

Actors and their consequences on business logic:

- **Agents (B2B)** belong to an office inside an organization. Lower markup, deferred billing — reservations are invoiced monthly to the office/organization. *Organic* organizations pay for all their offices; *inorganic* organizations' offices pay for themselves.
- **Customers (B2C)** must pay by credit card upfront before a booking is placed. No auth until checkout; a user is then created and logs in via phone + OTP.
- Roles are `customer | agent | admin | accountant` ([backend/services/accounts/db/models.go](backend/services/accounts/db/models.go)).

Order flow: location search → availability search (locations + dates) → choose plan → book → reservation created → (agents: voucher, then monthly invoice).

## Commands

Backend (`cd backend`, Docker must be running):

```bash
encore run              # run the whole backend; local Postgres/pubsub auto-provisioned
encore check            # compile-time check
encore test ./...       # ALWAYS run backend tests this way, never bare `go test`
encore test ./services/booking/... -run TestSearchAvailability   # single package / test
make gen                # regenerate sqlc bindings + gomock mocks after ANY query/migration change
encore db shell bookings           # psql into a local service DB
encore secret set --type local <Name>
```

Frontend (`cd frontend`):

```bash
pnpm dev      # Next dev (turbopack); Payload CMS starts inside the same process
pnpm build    # next build --webpack
pnpm lint
pnpm gen      # regenerate shared/client.ts from the LOCAL running backend
```

`pnpm gen` requires `encore run` to be live — regenerate it whenever a backend endpoint signature changes, or the frontend API modules will not typecheck.

PDF generation needs a Gotenberg container: `docker run --rm -p 8080:3000 gotenberg/gotenberg:8-chromium`.

## Backend architecture

Five Encore services under [backend/services/](backend/services/): `accounts`, `booking`, `reservation`, `billing`, `notifications`. Shared, non-API code lives in [backend/internal/](backend/internal/) (it defines no endpoints).

**Service layout convention.** A service package holds only endpoint declarations (thin) plus `<service>_service.go` with the `//encore:service` struct, `sqldb.NewDatabase`, and `initService`. Real logic lives in `handlers/<entity>/` subpackages, which receive `db.Querier` and caches by injection. Example: [services/booking/booking.go](backend/services/booking/booking.go) declares `SearchAvailability` and delegates to `handlers/availability`.

**Endpoint conventions.** Pointer receivers on the service struct; request type named `…Params`, response `…Response`, both pointers:
`func (s *Service) DoSomething(ctx context.Context, p *DoSomethingParams) (*DoSomethingResponse, error)`

**Authorization** is tag-driven, not per-handler. Global middlewares in [internal/middleware/require_role.go](backend/internal/middleware/require_role.go) map Encore tags to roles — annotate an endpoint `tag:admin`, `tag:agent`, `tag:customer`, or `tag:accountant` rather than checking `auth.Data()` inline. The auth handler is [services/accounts/auth_handler.go](backend/services/accounts/auth_handler.go); session cookie is `gr_session` (see [backend/encore.app](backend/encore.app)).

**Errors.** Never return raw errors from endpoints. Use the [internal/api_errors](backend/internal/api_errors/) builders; every error carries a `Details` struct with a stable `Code` string. The frontend maps that code straight to a translation key, so *adding a new error means adding a code in `codes_<service>.go` plus keys in both `messages/en.json` and `messages/he.json` under `ApiErrors`*.

**Persistence.** Per-service `db/` package: `query/<entity>_query.sql` (one file per entity), `migrations/NN_description.up.sql`, `sqlc.yaml` generating pgx/v5 code with `emit_interface` (so handlers depend on `db.Querier` and tests use the generated mocks in `services/*/mocks/`). Edit the `.sql` files and migrations, then `make gen` — never hand-edit `*.sql.go`, `models.go`, or the mocks.

**Broker abstraction.** [internal/broker/](backend/internal/broker/) hides supplier differences behind narrow interfaces — `LocationSearcher`, `AvailabilitySearcher`, `Booker`, `Canceler`, `VoucherProvider` — implemented per supplier in `flex_*.go` / `hertz_*.go` (Flex speaks SOAP/XML, Hertz speaks OTA XML). Supplier-specific quirks belong in those files; everything above the interface stays broker-agnostic.

**Availability → booking is snapshot-based.** A search computes a full price breakdown per plan (`PlanPriceDetails` in [handlers/availability/availability_snapshot.go](backend/services/booking/handlers/availability/availability_snapshot.go)), serializes all plans as JSON into `available_plans_snapshots`, and returns the snapshot id + plan id. Booking re-reads the snapshot rather than re-pricing, so **any new priced field must be added to `PlanPriceDetails` and threaded through booking/reservation**, not just to the search response.

**Pricing** is centralized in [internal/pricing/pricing.go](backend/internal/pricing/pricing.go) (markup → ERP → discount → totals). Markup rates are DB-driven per agent/office/organization; fallbacks and per-broker ERP day charges are Encore config in each service's `config.cue`.

**Async work.** Pub/Sub topics for emails ([services/notifications/events/](backend/services/notifications/events/)), OTP requests, and booking cancellations; cron jobs for open-reservation alerts and monthly billing. Invoicing integrates iCount ([internal/icount/](backend/internal/icount/)).

## Frontend architecture

Route groups under [frontend/app/](frontend/app/) separate the audiences: `(app)/[lang]/(withNavbar)/(booking)` is the booking funnel (`results` → `plans` → `order`), `(app)/admin/*` the admin portal, `(app)/accounting/*` the accountant portal, `(payload)/*` the CMS. `[lang]` drives next-intl; messages live in [frontend/messages/](frontend/messages/).

**API access is funneled** — never call the generated client directly. [shared/api/\_api.ts](frontend/shared/api/_api.ts) wraps every call: builds the client from `NEXT_PUBLIC_API_BASE_URL`/`NEXT_PUBLIC_ENCORE_ENV`, attaches the `X-Lang` header and bearer token, and converts Encore `APIError`s into `AppError` carrying the backend error `Code`. Per-entity modules (`booking-api.ts`, `reservations-api.ts`, …) sit on top; React Query hooks consume those. `shared/client.ts` is generated — do not edit.

**Errors in UI** go through `useTranslatedError` ([shared/hooks/useTranslatedError.ts](frontend/shared/hooks/useTranslatedError.ts)), which resolves `AppError.code` against the `ApiErrors` namespace.

**Components** are colocated: put a component in a `_components` folder at the nearest route segment; only genuinely reused ones go to [frontend/shared/components/](frontend/shared/components/). Auth state is Zustand ([shared/auth/authStore.ts](frontend/shared/auth/authStore.ts)). Theme tokens are defined in `globals.css`; use those Tailwind tokens with Shadcn components rather than raw colors.

## Working notes

- Both `//encore:api` and `// encore:api` (with a space) appear in the tree; match the surrounding file.
- Branches: `dev` deploys to stage, `main` to production. Backend deploys to Encore Cloud, frontend to Vercel.
- An Encore MCP server is configured in `.vscode/mcp.json` (`encore mcp run --app=ai-booking-travel-bo22`) for querying local endpoints, DBs, and API specs.
