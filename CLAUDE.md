# Repository notes

## Project conventions

- This repository is a Go CLI for HoYoLAB daily check-ins and LINE/Discord notifications.
- Keep existing CLI flags, YAML keys, API endpoints, notification payloads, and error prefixes backward compatible.
- Treat browser profile discovery and live HoYoLAB calls as external integrations; unit tests must not require a real browser profile, credentials, or outbound notification.
- Keep provider-specific notification payload construction in `NotifyMessage` and shared delivery/status handling in `sendNotification`.
- Run `gofmt` only on Go files changed for the current task. Before handoff, run `go test ./...`, `go vet ./...`, `go build ./...`, and `git diff --check`.
- Do not run dependency cleanup or upgrades unless explicitly requested because `go.mod` and `go.sum` may contain in-progress user work.

## Optimization record — 2026-08-02

- LINE Notify and Discord now reuse one Resty client during a notification batch instead of constructing a client for each provider.
- Empty messages and configurations without notification providers return before allocating an HTTP client.
- Shared request execution and HTTP status validation moved to `sendNotification`, removing duplicated provider error handling while preserving `line::` and `Discord::` error prefixes.
- The LINE authorization header now uses `SetHeader` directly instead of allocating a one-entry map.
- Added local regression tests for form delivery and non-200 status propagation without contacting external notification services.
- Updated `README.md` with development checks, notification architecture, test isolation, and the correct license link.

## Baseline record — 2026-08-02

- Pre-existing modified files: `act/act.go`, `act/api.go`, `go.mod`, `go.sum`, and `main.go`.
- Those changes include browser/profile discovery, user-agent selection, logging/CLI work, and dependency updates. They belong to the user and must not be reverted or rewritten.
- Baseline `go test ./...`: pass; the repository initially had no test files.
- Baseline `go vet ./...`: pass.
- Baseline standard `go build ./...`: blocked by VCS stamping because Go's Git subprocess does not inherit the command-level `safe.directory`.
- Baseline `go build -buildvcs=false ./...`: pass.
- `golangci-lint`: not installed in this environment.

## Verification record — 2026-08-02

- Final `gofmt -d act/notify.go act/notify_test.go`: pass with no diff.
- Final `go test -count=1 ./act`: pass.
- Final `go test -count=1 ./...`: pass.
- Final `go vet ./...`: pass.
- Final `go build -buildvcs=false ./...`: pass; standard `go build ./...` remains blocked only by the baseline VCS-stamping limitation.
- Final `git diff --check`: pass.

## Refactor record — 2026-10-06

- Game endpoints now live in `act/games.go` (GSI, HSR, ZZZ, HI3); `hoyolab.yaml` only lists `games: [labels]`. Legacy `config:` entries are still read, merged with built-in defaults, and migrated on write. The unused `profiles` key and per-game endpoint blocks are no longer written.
- Zenless Zone Zero (`ZZZ`) is enabled by default and sends the `x-rpc-signgame: zzz` header via the `signgame` field.
- `act/api.go` request/decoding logic is shared through generic `execute`; user-agent map key `"firefox "` typo fixed.
- `main.go` split into `findHoyolabStores`/`runProfile`/`formatAward`; fixed shadowed `resAcc` (nil deref in notification), cookie reuse across stores, award index bounds, and nil log writer in non-dev mode.
- Live run on this machine panics inside go-kooky Chrome cookie decryption (third-party, not changed).

## Browser support change — 2026-10-08

- Chrome, Edge and Brave are no longer read: their cookies use App-Bound Encryption (`v20`), which go-kooky intentionally does not implement. Cookies come from Firefox and forks (LibreWolf, Waterfox, Zen, Floorp); Opera is attempted best-effort and skipped with a warning when it is `v20`.
- `main.go` no longer imports `github.com/dvgamerr-app/go-kooky` root / `browser/all`. It uses the rewritten sibling library (`E:\.dvgamerr\go-kooky`): `browser.GetProfiles` + `CookieStore.CookiesForHost`, wrapped by `storeJar` so `act.IsCookieToken(http.CookieJar)` is unchanged. The `--list` base64 key format (`browser:profile`) is unchanged.
- Default config profile is now `firefox` with a Firefox user agent. Existing `profile[].browser: chrome|edge|brave` entries log a warning and match no store. A profile entry with a `browser` now only matches stores of that browser even when `name` is empty (previously empty `name` matched every browser).
- go-kooky was rewritten as module `github.com/dvgamerr-app/go-kooky/v2` (new API, incompatible with v1). `go.mod` now requires `.../go-kooky/v2 v2.0.0` and `main.go` imports `.../v2/browser`, and v2.0.0 is published (tag `v2.0.0` on branch `v2` of dvgamerr-app/go-kooky, commit 772bf33; pushed over ssh because the https/gh token lacked the `workflow` scope). `go.mod` was tidied and `go.sum` filled with no `go.work`. `github.com/go-ini/ini` is pinned to `v1.67.0` (indirect) because v1.67.3 declares its module path as `gopkg.in/ini.v1` and breaks the graph; do not bump it. `go.sum` for go-kooky/v2 was fetched with `GOPROXY=direct GONOSUMDB=github.com/dvgamerr-app/go-kooky` because proxy.golang.org / sum.golang.org had negatively cached the tag; re-verify with a normal `go mod verify` once those caches expire.
- Verified with `go.work`: `go vet ./...`, `go test -count=1 ./...`, `go build -buildvcs=false ./...` pass; live `--list` on this machine (Chrome/Edge/Brave only, no Firefox) exits with the new guidance message and no panic. Reading a real Firefox profile is unverified — none is installed here.

## Docs and site — 2026-10-10

- `README.md` / `README_TH.md` rewritten for the Firefox-only flow; keep the two in sync.
- GitHub Pages now serves the static page in `site/` (no Jekyll; `_config.yml` removed). `.github/workflows/gh-pages.yml` uploads `./site` and only runs when `site/**` or the workflow changes. Images in `site/img/` are copies of `docs/*.png`.
- The repository Pages source must be set to "GitHub Actions" (Settings → Pages) for the deploy job to publish.
