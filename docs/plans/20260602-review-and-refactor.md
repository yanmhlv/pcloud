# Review & Refactor — pCloud client (round 2)

## Overview

A focused review-and-refactor pass over `github.com/yanmhlv/pcloud`. The library
has already been through several refactor rounds (see
`docs/plans/architecture-findings.md` — H1–H3, M1–M4, L1–L8 all completed across
9 commits), so this round targets only what those passes left behind:

1. One genuine residual **data race** the thread-safety pass missed.
2. A **breaking API cleanup** (`ListShares`).
3. **Internal polish** (de-dup a magic number, simplify the upload pipe).
4. A repo-wide **comment strip** to bring the code in line with the
   no-comments convention, including removing the two files that exist purely
   as godoc documentation.

Scope: **Broad** (breaking changes allowed). Comments: **strip all**. Tests: **yes**.

## Context

Single-package library, ~3.5k LOC, no external deps beyond
`golang.org/x/oauth2` and `golang.org/x/time/rate`. `go vet` is clean and
`.golangci.yml` is comprehensive.

Key files:
- `client.go` — `Client`, `mu sync.RWMutex`, `request`/`do`/`doPost`, `setAuth`.
  All field access goes through `mu` **except** the one bug below.
- `file.go:213` — `downloadFromLink` calls `c.httpClient.Do(req)` **directly**,
  the only method that reads `c.httpClient` without `c.mu.RLock()`.
- `file.go:121-157` — upload uses `io.Pipe` + a goroutine + buffered `errCh`
  with three repeated `pw.CloseWithError(err); errCh <- err; return` branches.
- `client.go:46,72` — rate-limiter burst `10` is a literal duplicated in
  `NewClient` and `SetRateLimit`.
- `sharing.go:107` — `ListShares` returns `([]Share, []Share, error)` (two
  same-typed slices; caller can transpose them silently).
- `doc.go` — entire file is the package doc comment; its streaming example
  (`doc.go:43`) is stale (`GetFileLink(ctx, fileID)` — missing the `opts` arg
  added when the WithOpts variants were collapsed).
- `example_test.go` — `Example*` functions that serve as godoc examples.
- Doc-comment gaps today: `stream.go` `GetFileLink`/`GetFileLinkByPath`,
  `oauth.go` `EndpointUS`/`EndpointEU`, `errors.go` `IsErrorCode` + `ErrorCode`
  and its constants. (Moot once comments are stripped.)
- Tests use `httptest` (`testserver_test.go`) and run with `-race`; e2e tests in
  `tests/` are gated behind `PCLOUD_USERNAME`/`PCLOUD_PASSWORD`.

## Approach

Order the work so risk drops monotonically: fix the real bug first (with a
failing-then-passing race test), then the one breaking API change, then
behavior-preserving polish, and do the large mechanical comment strip **last**
so it never collides with the substantive edits. Every phase ends green under
`go test -race ./...` and `golangci-lint run`.

The comment strip is treated as the task (the user opted into it), not a
drive-by. Because `doc.go` and `example_test.go` are *entirely* godoc
documentation, "strip all comments" logically removes them — these are file
deletions and must be confirmed before execution.

Rejected alternatives:
- Keep godoc and only complete the missing comments — rejected; user chose strip-all.
- Additional breaking changes (e.g. reshaping `FileLink`, `Hash`) — considered
  and rejected; no concrete payoff, `ListShares` is the only clear win.

## Tasks

### Task 1: Correctness — fix the residual data race

- [x] In `downloadFromLink` (`file.go:213`), snapshot the HTTP client under
      `c.mu.RLock()` before `Do` — mirror the pattern in `request`
      (`client.go:89-94`): `c.mu.RLock(); httpCl := c.httpClient; c.mu.RUnlock()`.
- [x] Add a `-race` test: spawn concurrent `Download`/`DownloadByPath` calls
      against an `httptest` server while another goroutine calls
      `SetHTTPClient`, asserting no race (extend `client_test.go` /
      `file_test.go` in the style of `TestConcurrentLoginAndRequest`).
- [x] **Checkpoint:** `go test -race ./...` passes; new test fails before the
      fix and passes after.

### Task 2: Breaking API — `ListShares` return shape

- [x] Define an exported result type in `sharing.go`, e.g.
      `type Shares struct { Active []Share; Requests []Share }`.
- [x] Change `ListShares` to return `(Shares, error)`.
- [x] Update all callers: `sharing_test.go` (tests/e2e_test.go and
      example_test.go have no ListShares callers).
- [x] **Checkpoint:** `go build ./...` and `go test -race ./...` pass.

### Task 3: Internal polish (behavior-preserving)

- [x] Add `const rateLimiterBurst = 10` in `client.go`; use it in `NewClient`
      and `SetRateLimit` instead of the repeated literal.
- [x] Simplify the upload goroutine (`file.go:121-157`): collapse the three
      `pw.CloseWithError(err); errCh <- err; return` branches into a single
      inner func returning `error`, then `pw.CloseWithError(err)` once
      (`nil` → clean EOF) and one `errCh <- err`. Behavior identical.
- [x] **Checkpoint:** `go test -race ./...` and `golangci-lint run` pass.

### Task 4: Strip comments (mechanical, last)

- [x] Remove all doc/inline comments from source `.go` files: `client.go`,
      `auth.go`, `file.go`, `folder.go`, `stream.go`, `sharing.go`,
      `publink.go`, `oauth.go`, `errors.go`, `types.go`, `search.go`,
      `favorites.go`, `trash.go`, `thumb.go`, `revision.go`, `logger.go`.
- [x] Remove explanatory comments from test files (`*_test.go`), but **keep the
      `// Output:` lines in `example_test.go`** so the examples stay runnable and
      verified by `go test` (the one deliberate exception to strip-all).
- [x] Delete `doc.go` (CONFIRMED) — the whole file is the package doc comment;
      this also removes the stale example at `doc.go:43`.
- [x] Keep `example_test.go` (CONFIRMED): strip its prose comments, retain
      `// Output:` lines, ensure its `ListShares` usage matches the Phase 2 shape.
- [x] Re-run `golangci-lint run`; if `revive`'s `exported` / `package-comments`
      rules now fire (the config only disables `blank-imports` today), disable
      those two rules in `.golangci.yml` under `linters.settings.revive.rules`.
      (`staticcheck` ST1000/ST1003/ST1020-22 are already disabled.)
- [x] **Checkpoint:** `go build ./...`, `go test -race ./...`,
      `golangci-lint run` all pass.

### Task 5: Verify

- [x] `go vet ./...` clean.
- [x] `go test -race ./...` green (unit). E2e remains skipped without creds.
- [x] golangci-lint: only 24 pre-existing goconst warnings remain (out of scope;
      see review phases). No new lint categories appeared.
- [x] `go doc github.com/yanmhlv/pcloud` renders without the deleted symbols /
      stale example.

## Acceptance criteria

- `downloadFromLink` reads `c.httpClient` under the lock; the new concurrent
  Download + `SetHTTPClient` test passes under `-race`. The "All methods are
  safe for concurrent use" promise (`client.go:28`) — itself a comment slated
  for removal — is actually true.
- `ListShares` returns a single named struct; no caller passes two same-typed
  slices around.
- Rate-limiter burst is a single named constant; upload goroutine has one error
  path.
- No comments remain in any `.go` file; `doc.go` and `example_test.go` are gone
  (pending confirmation); lint and tests are green.

## Risks / open questions

- **File deletions need confirmation.** `doc.go` and `example_test.go` are the
  only files whose entire reason for existing is godoc documentation. Strip-all
  implies deleting them, but deletion is destructive — confirm before removing.
  Removing `doc.go` also drops the package overview from `pkg.go.dev`.
- **`example_test.go` `// Output:` lines.** These make examples *runnable and
  verified* by `go test`. A literal strip removes them, leaving compile-only
  examples (or, if the file is deleted, nothing). Recommended default: delete
  the file with `doc.go`; if kept, treat `// Output:` as the one pragmatic
  exception to strip-all and say so.
- **revive may start failing** on missing exported/package comments after the
  strip. Phase 4 budgets for adjusting `.golangci.yml`; confirm that editing the
  shared lint config is acceptable.
- **Breaking change blast radius.** `ListShares` is a public signature change;
  any external consumer breaks at compile time (acceptable per Broad scope, and
  consistent with the prior rounds' breaking changes).

## Post-completion

- Tag a new minor/major version reflecting the breaking `ListShares` change.
- Update `README` / usage docs (outside this repo's `.go` files) if they
  reference `ListShares` or link to the now-removed godoc examples.
