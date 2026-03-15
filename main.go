package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var (
	db       *sql.DB
	tmpl     *template.Template
	sessions = &sessionStore{tokens: make(map[string]time.Time)}
)

type sessionStore struct {
	sync.RWMutex
	tokens map[string]time.Time
}

func (s *sessionStore) set(token string) {
	s.Lock()
	defer s.Unlock()
	s.tokens[token] = time.Now().Add(24 * time.Hour)
}

func (s *sessionStore) valid(token string) bool {
	s.RLock()
	defer s.RUnlock()
	exp, ok := s.tokens[token]
	return ok && time.Now().Before(exp)
}

func (s *sessionStore) remove(token string) {
	s.Lock()
	defer s.Unlock()
	delete(s.tokens, token)
}

type Paste struct {
	ID        int
	Code      string
	Content   string
	CreatedAt string
}

func main() {
	var err error
	db, err = sql.Open("sqlite", "./uploadpass.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pastes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT UNIQUE NOT NULL,
			content TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)
	`); err != nil {
		log.Fatal(err)
	}
	initAdminPassword()

	tmpl = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/admin", handleAdminLogin)
	http.HandleFunc("/admin/dashboard", handleDashboard)
	http.HandleFunc("/admin/create", handleCreate)
	http.HandleFunc("/admin/delete", handleDelete)
	http.HandleFunc("/admin/change-admin-password", handleChangeAdminPassword)
	http.HandleFunc("/admin/change-link-password", handleChangeLinkPassword)
	http.HandleFunc("/admin/logout", handleLogout)
	http.HandleFunc("/s/", handleView)

	log.Println("Server running at http://localhost:9190")
	log.Fatal(http.ListenAndServe(":9190", nil))
}

func randomCode(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)[:n]
}

func isAuth(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	return sessions.valid(c.Value)
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if fwd := r.Header.Get("X-Forwarded-Proto"); fwd != "" {
		scheme = fwd
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func initAdminPassword() {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM settings WHERE key = 'admin_password'").Scan(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Exec("INSERT INTO settings (key, value) VALUES ('admin_password', ?)", string(hash))
	}
}

func checkAdminPassword(password string) bool {
	var hash string
	if err := db.QueryRow("SELECT value FROM settings WHERE key = 'admin_password'").Scan(&hash); err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if isAuth(r) {
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	if !checkAdminPassword(r.FormValue("password")) {
		tmpl.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Invalid password"})
		return
	}

	token := randomCode(32)
	sessions.set(token)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	})
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	rows, err := db.Query("SELECT id, code, content, created_at FROM pastes ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	var pastes []Paste
	for rows.Next() {
		var p Paste
		rows.Scan(&p.ID, &p.Code, &p.Content, &p.CreatedAt)
		pastes = append(pastes, p)
	}

	data := map[string]interface{}{
		"Pastes":  pastes,
		"BaseURL": baseURL(r),
	}
	if c, err := r.Cookie("flash_success"); err == nil {
		data["Success"] = c.Value
		http.SetCookie(w, &http.Cookie{Name: "flash_success", Path: "/", MaxAge: -1})
	}
	if c, err := r.Cookie("flash_error"); err == nil {
		data["Error"] = c.Value
		http.SetCookie(w, &http.Cookie{Name: "flash_error", Path: "/", MaxAge: -1})
	}
	tmpl.ExecuteTemplate(w, "dashboard.html", data)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !isAuth(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	password := r.FormValue("password")
	if content == "" || password == "" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	code := randomCode(8)
	if _, err = db.Exec("INSERT INTO pastes (code, content, password) VALUES (?, ?, ?)",
		code, content, string(hash)); err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	link := fmt.Sprintf("%s/s/%s", baseURL(r), code)
	tmpl.ExecuteTemplate(w, "created.html", map[string]string{"Link": link})
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !isAuth(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	db.Exec("DELETE FROM pastes WHERE code = ?", r.FormValue("code"))
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func handleChangeAdminPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !isAuth(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	current := r.FormValue("current_password")
	newPass := r.FormValue("new_password")
	confirm := r.FormValue("confirm_password")

	if !checkAdminPassword(current) {
		redirectWithMsg(w, r, "error", "Current password is incorrect")
		return
	}
	if len(newPass) < 4 {
		redirectWithMsg(w, r, "error", "New password must be at least 4 characters")
		return
	}
	if newPass != confirm {
		redirectWithMsg(w, r, "error", "New passwords do not match")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		redirectWithMsg(w, r, "error", "Server error")
		return
	}
	db.Exec("UPDATE settings SET value = ? WHERE key = 'admin_password'", string(hash))
	redirectWithMsg(w, r, "success", "Admin password changed successfully")
}

func handleChangeLinkPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !isAuth(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	code := r.FormValue("code")
	newPass := r.FormValue("new_password")
	if code == "" || newPass == "" {
		redirectWithMsg(w, r, "error", "Code and password are required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		redirectWithMsg(w, r, "error", "Server error")
		return
	}
	res, _ := db.Exec("UPDATE pastes SET password = ? WHERE code = ?", string(hash), code)
	if n, _ := res.RowsAffected(); n == 0 {
		redirectWithMsg(w, r, "error", "Link not found")
		return
	}
	redirectWithMsg(w, r, "success", "Password updated for link "+code)
}

func redirectWithMsg(w http.ResponseWriter, r *http.Request, kind, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "flash_" + kind,
		Value:  msg,
		Path:   "/",
		MaxAge: 5,
	})
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("session"); err == nil {
		sessions.remove(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func handleView(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/s/")
	if code == "" {
		http.NotFound(w, r)
		return
	}

	var content, hash string
	if err := db.QueryRow("SELECT content, password FROM pastes WHERE code = ?", code).
		Scan(&content, &hash); err != nil {
		tmpl.ExecuteTemplate(w, "notfound.html", nil)
		return
	}

	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "unlock.html", map[string]interface{}{"Code": code})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(r.FormValue("password"))); err != nil {
		tmpl.ExecuteTemplate(w, "unlock.html", map[string]interface{}{
			"Code":  code,
			"Error": "Wrong password",
		})
		return
	}

	tmpl.ExecuteTemplate(w, "view.html", map[string]string{"Content": content})
}
