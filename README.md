# UploadPass

A simple Go web app that lets an admin create password-protected text links.

## How it works

1. Admin logs in at `/admin` with password: `password`
2. Admin enters text content and sets an access password → gets a short link
3. Anyone with the link enters the password → sees the original text

## Run

```bash
go build -o uploadpass .
./uploadpass
```

Server starts at **http://localhost:8080**

## Tech Stack

- Go (net/http, html/template)
- SQLite (github.com/mattn/go-sqlite3)
- bcrypt for password hashing
