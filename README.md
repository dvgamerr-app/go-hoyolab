<div align="center">

# 🎮 go-hoyolab

**Daily HoYoLAB check-in for every HoYoverse game — set it up once, forget it forever.**

[![Build](https://img.shields.io/github/actions/workflow/status/dvgamerr-app/go-hoyolab/build.yml?label=Build&style=flat-square)](https://github.com/dvgamerr-app/go-hoyolab/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/dvgamerr-app/go-hoyolab?style=flat-square)](https://goreportcard.com/report/dvgamerr-app/go-hoyolab)
[![Release](https://img.shields.io/github/v/release/dvgamerr-app/go-hoyolab?style=flat-square)](https://github.com/dvgamerr-app/go-hoyolab/releases/latest)
[![License](https://img.shields.io/github/license/dvgamerr-app/go-hoyolab?style=flat-square)](LICENSE)
[![PRs welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](https://github.com/dvgamerr-app/go-hoyolab/pulls)

[English](./README.md) · [ภาษาไทย](./README_TH.md) · [Website](https://dvgamerr-app.github.io/go-hoyolab/) · [Discord](https://discord.gg/QDccF497Mw)

![check-in log](./docs/example-logs.png)

</div>

## ✨ Why

- **Zero setup for login** — it reuses the HoYoLAB session you already have in Firefox. No passwords, no tokens to copy.
- **One small binary** — a single `hoyolab` executable for Windows, macOS, Linux and FreeBSD. No runtime, no Docker.
- **Every game in one run** — all your enabled games are checked in together, for as many accounts as you have profiles.
- **Tells you what you got** — a Discord (or LINE) message lists today's reward for each game.
- **Updates itself** — picks up new releases from GitHub automatically.

## 🕹️ Supported games

| Game | Key in `hoyolab.yaml` | Default |
| --- | --- | :---: |
| Genshin Impact | `GSI` | ✅ |
| Honkai: Star Rail | `HSR` | ✅ |
| Zenless Zone Zero | `ZZZ` | ✅ |
| Honkai Impact 3rd | `HI3` | — |

## 🚀 Quick start

1. Sign in to <https://www.hoyolab.com/home> in **Firefox** (or LibreWolf, Waterfox, Zen, Floorp).
2. Download the archive for your OS from the [latest release](https://github.com/dvgamerr-app/go-hoyolab/releases/latest) and unpack it.
3. Run it:

   ```shell
   ./hoyolab          # Windows: hoyolab.exe
   ```

That is all. On the first run it writes a `hoyolab.yaml` next to the binary; edit it to add notifications or switch games on and off.

> [!IMPORTANT]
> **Chrome, Edge and Brave are no longer supported.** They encrypt cookies with App-Bound Encryption (`v20`), which cannot be read from outside the browser. If your old `hoyolab.yaml` has `profile[].browser: chrome` (or `edge` / `brave`), change it to `firefox` — those entries are ignored with a warning.

To see which browser profiles were found:

```shell
./hoyolab --list
```

## ⚙️ Configuration

```yaml
notify:
  discord-webhook: "https://discord.com/api/webhooks/..."   # optional
  line-notify: ""                                           # optional
  mini: true
profile:
  - browser: firefox
    name: []            # empty = every Firefox profile; or list profile names
games:                  # remove or comment out the games you do not play
  - GSI
  - HSR
  - ZZZ
```

Existing configs that still use the old per-game `config:` blocks keep working and are migrated automatically the next time the file is written.

## 🔔 Notifications

| Provider | How |
| --- | --- |
| **Discord** | Create a webhook in *Channel settings → Integrations* and paste the URL into `notify.discord-webhook`. |
| **LINE Notify** | Generate a token at [notify-bot.line.me](https://notify-bot.line.me/my/) and paste it into `notify.line-notify`. LINE has retired this service, so prefer Discord. |

![Discord notification](./docs/example-notify.png)

## ⏰ Run it every day

**Windows** — *Task Scheduler → Create Basic Task → Daily (e.g. 5 am) → Start a program → `hoyolab.exe`*.
Or put a shortcut to `hoyolab.exe` in `shell:startup` to check in each time the PC boots.

**macOS / Linux** — add a cron job:

```cron
0 5 * * * /path/to/hoyolab
```

## 🛠️ Development

```shell
go test ./...
go vet ./...
go build ./...
```

- Tests never touch a real browser profile, real credentials or an outbound notification service.
- Cookie reading lives in [`go-kooky/v2`](https://github.com/dvgamerr-app/go-kooky); game endpoints are in [`act/games.go`](./act/games.go); notification delivery is in [`act/notify.go`](./act/notify.go).
- If `go build` fails on VCS stamping, use `go build -buildvcs=false ./...`.

Contributions are welcome — please open an [issue](https://github.com/dvgamerr-app/go-hoyolab/issues) first for anything large.

## 💬 Help

Found a bug or want a game added? [Open an issue](https://github.com/dvgamerr-app/go-hoyolab/issues), or come say hi on Discord.

[![Join us on Discord](https://discordapp.com/api/guilds/475720106471849996/widget.png?style=banner2)](https://discord.gg/QDccF497Mw)

## 📄 License

[MIT](./LICENSE) © Kananek Thongkam
