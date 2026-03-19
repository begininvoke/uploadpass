<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/SQLite-003B57?style=flat-square&logo=sqlite&logoColor=white" alt="SQLite">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
</p>

# UploadPass

A lightweight, self-hosted Go web app for sharing text securely through password-protected short links. Zero external dependencies, single binary, dark modern UI.

---

## Why UploadPass?

Sometimes you need to share sensitive text — credentials, API keys, private notes — without leaving it exposed in chat logs or emails. UploadPass generates a short link protected by a password. The recipient opens the link, enters the password, and reads the text. Links auto-expire after 30 days.

---

## Features

- **Multi-admin support** — register multiple admin accounts, each with their own links
- **Password-protected links** — every link requires a password to unlock
- **Short URLs** — auto-generated 8-character codes (`/s/abc12345`)
- **Auto-expiry** — links are automatically deleted after 30 days
- **Per-admin isolation** — each admin only sees and manages their own links
- **Password management** — change your admin password or any link's password anytime
- **One-click copy** — click a link code to copy the full URL to clipboard
- **Dark modern UI** — clean, responsive design that works on desktop and mobile
- **Single binary** — no runtime dependencies, just run it
- **SQLite storage** — database file auto-created on first run

---

## Quick Start

### From source

```bash
git clone https://github.com/youruser/uploadpass.git
cd uploadpass
go build -o uploadpass .
./uploadpass
```

### Pre-built binaries

Download from [Releases](https://github.com/youruser/uploadpass/releases) — available for Linux, macOS, and Windows (amd64, arm64, 386).

### Run it

```bash
./uploadpass
```

Server starts at **http://localhost:9190**

---

## Default Credentials

| Field    | Value      |
|----------|------------|
| Username | `admin`    |
| Password | `password` |

> Change these immediately after first login via the Settings modal on the dashboard.

---

## How It Works

```
Admin                              Recipient
  |                                    |
  |  1. Login at /admin                |
  |  2. Enter text + set password      |
  |  3. Get short link  ------------->  |
  |                                    |  4. Open link
  |                                    |  5. Enter password
  |                                    |  6. Read the text
```

1. **Admin** logs in and creates a new link with some text and an access password
2. **Admin** shares the generated short link (e.g. `https://example.com/s/a1b2c3d4`)
3. **Recipient** opens the link, enters the password, and views the text
4. After **30 days**, the link is automatically deleted

---

## Screenshots

<details>
<summary>Login Page</summary>

Dark themed login page with username and password fields.

</details>

<details>
<summary>Dashboard</summary>

Admin dashboard showing the create form and a list of your links with expiry dates.

</details>

<details>
<summary>Unlock Page</summary>

Clean password prompt that recipients see when opening a shared link.

</details>

---

## Project Structure

```
uploadpass/
├── main.go              # Server, routes, all handlers
├── templates/
│   ├── login.html       # Admin login page
│   ├── register.html    # Admin registration page
│   ├── dashboard.html   # Admin dashboard
│   ├── created.html     # Link created confirmation
│   ├── unlock.html      # Password entry for recipients
│   ├── view.html        # Display unlocked text
│   └── notfound.html    # Expired/invalid link page
├── build.sh             # Cross-platform build script
├── go.mod
├── go.sum
└── uploadpass.db        # SQLite database (auto-created)
```

---

## API Routes

| Method | Path | Auth | Description |
|--------|------|:----:|-------------|
| `GET`  | `/admin` | No | Login page |
| `POST` | `/admin` | No | Authenticate |
| `GET`  | `/admin/register` | No | Registration page |
| `POST` | `/admin/register` | No | Create new admin account |
| `GET`  | `/admin/dashboard` | Yes | Dashboard with your links |
| `POST` | `/admin/create` | Yes | Create a new link |
| `POST` | `/admin/delete` | Yes | Delete a link |
| `POST` | `/admin/change-admin-password` | Yes | Change your password |
| `POST` | `/admin/change-link-password` | Yes | Change a link's password |
| `GET`  | `/admin/logout` | Yes | Log out |
| `GET`  | `/s/{code}` | No | Unlock form for a link |
| `POST` | `/s/{code}` | No | Submit password and view text |

---

## Cross-Platform Build

Build binaries for all platforms at once:

```bash
chmod +x build.sh
./build.sh
```

Outputs to `build/` directory:

| Platform        | Binary |
|-----------------|--------|
| macOS amd64     | `uploadpass-darwin-amd64` |
| macOS arm64     | `uploadpass-darwin-arm64` |
| Linux amd64     | `uploadpass-linux-amd64` |
| Linux arm64     | `uploadpass-linux-arm64` |
| Linux 386       | `uploadpass-linux-386` |
| Linux arm       | `uploadpass-linux-arm` |
| Windows amd64   | `uploadpass-windows-amd64.exe` |
| Windows arm64   | `uploadpass-windows-arm64.exe` |
| Windows 386     | `uploadpass-windows-386.exe` |

All binaries are statically compiled with `CGO_ENABLED=0` — no dependencies needed on the target machine.

---

## Deployment

### Systemd (Linux)

```ini
[Unit]
Description=UploadPass
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/uploadpass
ExecStart=/opt/uploadpass/uploadpass
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Docker

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o uploadpass .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/uploadpass .
COPY --from=builder /app/templates ./templates
EXPOSE 9190
CMD ["./uploadpass"]
```

### Reverse Proxy (Nginx)

```nginx
server {
    listen 443 ssl;
    server_name paste.example.com;

    location / {
        proxy_pass http://127.0.0.1:9190;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.23 (stdlib `net/http`, `html/template`) |
| Database | SQLite via [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO) |
| Passwords | [`golang.org/x/crypto/bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |
| Sessions | In-memory token store (24h expiry) |

---

## License

MIT
