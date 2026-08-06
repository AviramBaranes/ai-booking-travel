---
name: test-endpoint
description: Write Go tests for an Encore API endpoint in this app (validation, domain errors, happy path) following the repo's test conventions. Use when the user says "write tests for <endpoint>", "test the X endpoint", "add tests for this handler", or runs /test-endpoint.
---

# test-endpoint

Write tests for an `//encore:api` endpoint, following this repo's conventions.

Test files live next to the thin endpoint file in the service package (e.g. `inspection_centers.go` →
`inspection_centers_test.go` in `package get_car`), **not** in the `handlers/<area>` package.

## Workflow

### 1. Read the code under test

- The thin endpoint file (`services/get_car/<area>.go`) — note the `//encore:api` line: access level
  (`public`/`auth`), method, path, and any `tag:`.
- The handler (`services/get_car/handlers/<area>/<file>.go`) — params + `validate:` tags, response type,
  every early-return error, which querier methods it calls, and **whether it reads `auth.Data()`**.
- The domain error vars (`handlers/<area>/<area>_service.go`) and `internal/api_errors/codes_*.go`.
- An existing test for the pattern: `services/get_car/inspection_centers_test.go`.

### 2. Enumerate the cases

Cover, in this order:

1. **Validation** — one case per `validate:` rule, table-driven, plus one "valid params" case.
2. **Domain errors** — every non-happy return in the handler (already-exists, not-found, ...). Set up real
   state to trigger them; don't mock the querier.
3. **Happy path** — assert the response *and* read the row back from the DB to verify what was persisted.

### 3. Write the test

**Call the endpoint directly — never `initService()`.** Encore generates a package-level function in
`encore.gen.go` for each service-struct method; call that. Encore handles the rest.

```go
resp, err := CreateInspectionCenter(ctx, p)     // ✅
resp, err := svc.CreateInspectionCenter(ctx, p) // ❌ never
```

**Auth: usually not needed.** Calling an endpoint directly does not run the auth handler, and `auth` access
level / `tag:role` gating does not reject the call. Only authenticate when the **handler itself reads the
auth data** (`auth.Data().(*AuthData)` / `GetAuthData()`) — then the data must be there or the handler
misbehaves. Authenticate by putting the auth data on the **context**, not with `et.OverrideAuthInfo`:

```go
// services/get_car/test_utils_test.go
func authContext(ctx context.Context, userID int64, role UserRole) context.Context {
    uid := auth.UID(strconv.FormatInt(userID, 10))
    return auth.WithContext(ctx, uid, &AuthData{UserID: userID, Role: role})
}

// in the test
ctx := authContext(context.Background(), 1, UserRoleAdmin)
resp, err := ListAdmins(ctx)
```

Role-gating itself is covered by the middleware tests in `internal/middleware/require_role_test.go` — don't
re-test it per endpoint.

**All DB access goes through sqlc.** Never write a raw SQL string in a test. To seed, read back, or clean
up, add the query to `services/get_car/db/query/*_query.sql`, run `make gen` from `backend/`, and call the
generated `db.Querier` method via `testQuerier()`.

**Clean up every row you create**, with `t.Cleanup` so it runs even on failure:

```go
createCenter := func(t *testing.T, p inspection.CreateInspectionCenterParams) (*inspection.CreateInspectionCenterResponse, error) {
    t.Helper()
    resp, err := CreateInspectionCenter(ctx, p)
    if resp != nil {
        t.Cleanup(func() {
            if err := query.DeleteInspectionCenter(ctx, resp.ID); err != nil {
                t.Errorf("failed to delete test inspection center %d: %v", resp.ID, err)
            }
        })
    }
    return resp, err
}
```

**Validation cases** are table-driven over a mutation of a known-good params factory, asserting the exact
field:

```go
cases := []struct {
    name   string
    field  string
    mutate func(p *inspection.CreateInspectionCenterParams)
}{
    {name: "Missing name", field: "name", mutate: func(p *inspection.CreateInspectionCenterParams) { p.Name = "" }},
    {name: "Latitude too high", field: "latitude", mutate: func(p *inspection.CreateInspectionCenterParams) { p.Latitude = 90.1 }},
}

for _, c := range cases {
    c := c
    t.Run(c.name, func(t *testing.T) {
        p := validCenterParams()
        c.mutate(&p)

        expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
            Code:  api_errors.CodeInvalidValue,
            Field: c.field,
        })
        api_errors.AssertApiError(t, expectedErr, p.Validate())
    })
}
```

**Assert errors with `api_errors.AssertApiError(t, want, got)`** — it compares code, message and details.
Never assert on `err.Error()` strings.

**Helpers go in a test utils file, not inline**:

- Generic and cross-package (random names/emails, ...) → `internal/testutils` (`testutils.RandomName(prefix)`).
  Keep them parameterized, never one-off per test.
- Package-scoped (`testQuerier()`, `authContext()`, row factories like `createTestUser`) →
  `services/get_car/test_utils_test.go`.
- Only genuinely single-use helpers stay in the test file itself.

Use unique names/emails for every created row (`testutils.RandomName("center")`) so runs don't collide.

### 4. Run and report

```bash
encore test ./services/get_car/... -run TestCreateInspectionCenter -v
```

Then run the full suite (`encore test ./...`) before reporting — endpoint tests share a database and can
surface breakage elsewhere. Report failures with the real output; never claim green without a passing run.

## Checklist

- [ ] Endpoint called via the package-level function, no `initService()`
- [ ] Auth only where the handler reads it, via `authContext` on the ctx
- [ ] One validation case per `validate:` rule + a valid-params case
- [ ] Every handler error branch covered
- [ ] Happy path verifies the persisted row, not just the response
- [ ] Zero raw SQL — new queries added to `query/*.sql` + `make gen`
- [ ] Created rows cleaned up in `t.Cleanup`
- [ ] Reusable helpers in `internal/testutils` or `test_utils_test.go`
- [ ] `encore test ./...` passes
