# UploadPass

A lightweight Go web app for sharing text securely through password-protected short links.

## Features

- **Admin panel** — login, create links, manage everything from a clean dashboard
- **Password-protected links** — each link requires a password to view the content
- **Short URLs** — auto-generated 8-character codes (`/s/abc12345`)
- **Change passwords** — update admin password or any link's password anytime
- **Copy to clipboard** — one-click copy for generated links
- **Delete links** — remove links you no longer need
- **Dark modern UI** — responsive design, works on desktop and mobile

## Quick Start

```bash
go build -o uploadpass .
./uploadpass
```

Server starts at **http://localhost:9190**

## Default Admin Password

```
password
```

You can change it from the dashboard after logging in.

## How It Works

| Step | Who | What |
|------|-----|------|
| 1 | Admin | Log in at `/admin` |
| 2 | Admin | Enter text + set access password → get a short link |
| 3 | Anyone | Open the short link → enter password → view the text |

## Project Structure

```
uploadpass/
├── main.go              # Server, routes, handlers
├── templates/
│   ├── login.html       # Admin login page
│   ├── dashboard.html   # Admin dashboard
│   ├── created.html     # Link created confirmation
│   ├── unlock.html      # Password entry for viewers
│   ├── view.html        # Display unlocked text
│   └── notfound.html    # Invalid/deleted link page
├── go.mod
├── go.sum
└── uploadpass.db        # SQLite database (auto-created)
```

## Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/admin` | Admin login page |
| POST | `/admin` | Authenticate admin |
| GET | `/admin/dashboard` | Dashboard (requires auth) |
| POST | `/admin/create` | Create new link |
| POST | `/admin/delete` | Delete a link |
| POST | `/admin/change-admin-password` | Change admin password |
| POST | `/admin/change-link-password` | Change a link's password |
| GET | `/s/{code}` | Enter password to view |
| POST | `/s/{code}` | Submit password, view text |

## Tech Stack

- **Go** — standard library (`net/http`, `html/template`)
- **SQLite** — `github.com/mattn/go-sqlite3`
- **bcrypt** — `golang.org/x/crypto/bcrypt` for password hashing
