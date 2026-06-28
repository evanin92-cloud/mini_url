package main

import (
"encoding/json"
"html/template"
"log"
"net/http"
"sync"
"time"
)

type Link struct {
ID          string    `json:"id"`
OriginalURL string    `json:"original_url"`
ShortID     string    `json:"short_id"`
CreatedAt   time.Time `json:"created_at"`
Clicks      int       `json:"clicks"`
}

var (
links   = make(map[string]*Link)
mu      sync.RWMutex
counter = 1000
)

func main() {
http.HandleFunc("/", indexHandler)
http.HandleFunc("/shorten", shortenHandler)
http.HandleFunc("/api/shorten", apiShortenHandler)
http.HandleFunc("/links", linksHandler)
http.HandleFunc("/{shortID}", redirectHandler)

log.Println("🚀 Сервер запущен на http://localhost:8080")
log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
tmpl := ` + "`" + `<!DOCTYPE html>
<html>
<head><title>Mini URL</title></head>
<body style="font-family:Arial;max-width:600px;margin:50px auto;padding:20px;">
<h1>🔗 Mini URL</h1>
<form method="POST" action="/shorten">
<input type="url" name="url" placeholder="Введите URL..." style="width:70%;padding:10px;">
<button type="submit" style="padding:10px 20px;">Сократить</button>
</form>
{{range .}}
<div style="border:1px solid #ddd;padding:10px;margin:10px 0;">
<a href="/{{.ShortID}}" target="_blank">http://localhost:8080/{{.ShortID}}</a>
<span style="color:gray;">→ {{.OriginalURL}}</span>
<span style="float:right;">👁 {{.Clicks}}</span>
</div>
{{else}}
<p>Нет ссылок</p>
{{end}}
</body>
</html>` + "`" + `
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
if r.Method != "POST" {
http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
return
}
var req struct{ URL string }
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
shortID := r.URL.Path[1:]
if shortID == "" || shortID == "favicon.ico" {
http.NotFound(w, r)
return
}
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
