# gostatus
Discord presence badges as SVGs to show discord status, music-, editor- or game-activity using [shields.io](https://shields.io/).

> [!NOTE]
> Originally inspired by [statusbadges](https://github.com/advaith1/statusbadges), rewritten in Go with improvements.

## Running

### Requirements
- [Go 1.21+](https://go.dev/)
- **Server Members** and **Presence** intents enabled at the [Dev Portal](https://discord.com/developers/home)
- Bot must be in a shared server with the users you want to track

```bash
git clone https://github.com/vmphase/gostatus
cd gostatus
go mod tidy
go run main.go
```

Configuration is done via `config.toml.example`:
- rename it to `config.toml`
- set `token` to your bot token
- optionally change the `port` field

## Endpoints

### `GET /badge/status/{discord_user_id}`
Current Discord presence status.

| Query        | Default      | Description |
| ------------ | ------------ | ----------- |
| `label`      | `currently`  | Left side text |
| `color`      | status-based | Right side background color |
| `labelColor` | `#555`       | Left side background color |
| `simple`     | —            | Set to `true` to collapse `idle`/`dnd` → `online` |

---

### `GET /badge/music/{discord_user_id}`
Track the user is currently listening to. Auto-detects supported music services (currently Spotify).

| Query        | Default        | Description |
| ------------ | -------------- | ----------- |
| `label`      | `listening to` | Left side text |
| `color`      | service-based  | Right side background color |
| `labelColor` | `#555`         | Left side background color |
| `fallback`   | `nothing`      | Text shown when not listening |
| `hideLogo`   | `false`        | Set to `true` to hide the service logo |

---

### `GET /badge/code/{discord_user_id}`
File and workspace the user is currently editing. Auto-detects supported editors (VSCode, Zed).

| Query        | Default       | Description |
| ------------ | ------------- | ----------- |
| `label`      | editor-based  | Left side text |
| `color`      | editor-based  | Right side background color |
| `labelColor` | `#1e1e2e`    | Left side background color |
| `fallback`   | `nothing`     | Text shown when not coding |
| `hideLogo`   | `false`       | Set to `true` to hide the editor logo |
| `prefer`     | first active  | Preferred editor slug (`vscode`, `zed`) when multiple are active |

---

### `GET /badge/playing/{discord_user_id}`
Game the user is currently playing (editor activities are excluded).

| Query        | Default    | Description |
| ------------ | ---------- | ----------- |
| `label`      | `playing`  | Left side text |
| `color`      | `#5865f2`  | Right side background color |
| `labelColor` | `#555`     | Left side background color |
| `fallback`   | `nothing`  | Text shown when not playing |

---

### `GET /badge/crunchyroll/{discord_user_id}`
Episode and series the user is currently watching on [Crunchyroll](https://www.crunchyroll.com/).

| Query        | Default    | Description |
| ------------ | ---------- | ----------- |
| `label`      | `watching` | Left side text |
| `color`      | `#5865f2`  | Right side background color |
| `labelColor` | `#555`     | Left side background color |
| `fallback`   | `nothing`  | Text shown when not watching |
| `hideLogo`   | `false`    | Set to `true` to hide the Crunchyroll logo |

---

### `GET /presence/{discord_user_id}`
Raw presence data as JSON. CORS-enabled.
