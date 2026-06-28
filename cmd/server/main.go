package main

import (
"encoding/json"
"html/template"
"log"
"net/http"
"sync"
"time"

"github.com/go-chi/chi/v5"
"github.com/go-chi/chi/v5/middleware"
)

type Link struct {
ID          string     + "json:\"id\"" + 
OriginalURL string     + "json:\"original_url\"" + 
ShortID     string     + "json:\"short_id\"" + 
CreatedAt   time.Time  + "json:\"created_at\"" + 
Clicks      int        + "json:\"clicks\"" + 
}

var (
links   = make(map[string]*Link)
mu      sync.RWMutex
counter = 1000
)

func main() {
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)

r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

r.Get("/", indexHandler)
r.Post("/shorten", shortenHandler)
r.Post("/api/v1/shorten", apiShortenHandler)
r.Get("/{shortID}", redirectHandler)
r.Get("/links", linksHandler)

log.Println("🚀 Сервер запущен на http://localhost:8080")
log.Println("📍 Откройте в браузере: http://localhost:8080")
log.Fatal(http.ListenAndServe(":8080", r))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
tmpl :=  + "" + 
<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Mini URL</title>
<style>
* { margin:0; padding:0; box-sizing:border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; color: #1a1a2e; }
.container { max-width: 900px; margin: 0 auto; padding: 20px; }
.header { background: #ffffff; padding: 16px 0; border-bottom: 1px solid #e1e5eb; }
.header .container { display: flex; justify-content: space-between; align-items: center; }
.logo { font-size: 24px; font-weight: 700; color: #4a6cf7; }
.nav a { color: #4a5568; text-decoration: none; margin-left: 20px; padding: 6px 12px; border-radius: 6px; }
.nav a:hover { color: #4a6cf7; background: #eef2ff; }
.hero { text-align: center; padding: 40px 0 30px; }
.hero h1 { font-size: 38px; font-weight: 800; background: linear-gradient(135deg, #4a6cf7, #6c63ff); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
.hero p { font-size: 18px; color: #718096; }
.card { background: #ffffff; border-radius: 16px; padding: 32px; box-shadow: 0 4px 20px rgba(0,0,0,0.06); border: 1px solid #edf2f7; max-width: 700px; margin: 0 auto 30px; }
.form-control { width: 100%; padding: 14px 18px; border: 2px solid #e2e8f0; border-radius: 12px; font-size: 16px; background: #f7fafc; }
.form-control:focus { border-color: #4a6cf7; outline: none; }
.btn { display: inline-block; padding: 12px 28px; border: none; border-radius: 12px; font-size: 16px; font-weight: 600; cursor: pointer; }
.btn-primary { background: linear-gradient(135deg, #4a6cf7, #6c63ff); color: #fff; }
.btn-primary:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(74,108,247,0.35); }
.link-card { background: #fff; border-radius: 12px; padding: 16px 20px; border: 1px solid #edf2f7; margin-bottom: 12px; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px; }
.link-card:hover { border-color: #4a6cf7; }
.link-short { color: #4a6cf7; font-weight: 600; }
.footer { text-align: center; padding: 24px 0; color: #718096; font-size: 14px; border-top: 1px solid #e1e5eb; margin-top: 40px; }
</style>
</head>
<body>
<header class="header">
<div class="container">
<div class="logo">🔗 Mini URL</div>
<nav class="nav">
<a href="/">Главная</a>
<a href="/links">Мои ссылки</a>
</nav>
</div>
</header>
<div class="container">
<section class="hero">
<h1>Сократите ваши ссылки</h1>
<p>Быстрое и удобное создание коротких ссылок</p>
</section>
<div class="card">
<h2>Создать короткую ссылку</h2>
<form id="shorten-form" method="POST" action="/shorten">
<input type="url" id="original-url" name="url" class="form-control" placeholder="Введите длинный URL..." required>
<button type="submit" class="btn btn-primary" style="margin-top:16px;">Сократить</button>
</form>
</div>
<div class="links-section">
<h2>📋 Последние ссылки</h2>
{{range .}}
<div class="link-card">
<span>{{.OriginalURL}}</span>
<a href="/{{.ShortID}}" target="_blank" class="link-short">http://localhost:8080/{{.ShortID}}</a>
<span>👁 {{.Clicks}} переходов</span>
</div>
{{else}}
<p>Нет ссылок. Создайте свою первую!</p>
{{end}}
</div>
<footer class="footer"><p>&copy; 2026 Mini URL</p></footer>
</div>
</body>
</html>
 + "" + 
t, _ := template.New("index").Parse(tmpl)
mu.RLock()
var list []*Link
for _, l := range links {
list = append(list, l)
}
mu.RUnlock()
t.Execute(w, list)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
url := r.FormValue("url")
if url == "" {
http.Error(w, "URL обязателен", http.StatusBadRequest)
return
}
mu.Lock()
counter++
shortID := generateShortID(counter)
links[shortID] = &Link{ID: shortID, OriginalURL: url, ShortID: shortID, CreatedAt: time.Now()}
mu.Unlock()
http.Redirect(w, r, "/", http.StatusSeeOther)
}

func apiShortenHandler(w http.ResponseWriter, r *http.Request) {
var req struct{ URL string  + "json:\"url\"" +  }
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Неверный формат", http.StatusBadRequest)
return
}
if req.URL == "" {
http.Error(w, "URL обязателен", http.StatusBadRequest)
return
}
mu.Lock()
counter++
shortID := generateShortID(counter)
links[shortID] = &Link{ID: shortID, OriginalURL: req.URL, ShortID: shortID, CreatedAt: time.Now()}
mu.Unlock()
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"short_url": "http://localhost:8080/" + shortID})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
shortID := chi.URLParam(r, "shortID")
mu.RLock()
link, ok := links[shortID]
mu.RUnlock()
if !ok {
http.NotFound(w, r)
return
}
mu.Lock()
link.Clicks++
mu.Unlock()
http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func linksHandler(w http.ResponseWriter, r *http.Request) {
mu.RLock()
var list []*Link
for _, l := range links {
list = append(list, l)
}
mu.RUnlock()
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(list)
}

func generateShortID(n int) string {
chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
result := ""
for n > 0 {
result = string(chars[n%62]) + result
n = n / 62
}
for len(result) < 6 {
result = "a" + result
}
return result
}
