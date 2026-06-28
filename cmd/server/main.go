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

// Структура ссылки
type Link struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortID     string    `json:"short_id"`
	CreatedAt   time.Time `json:"created_at"`
	Clicks      int       `json:"clicks"`
}

// Глобальное хранилище
var (
	links   = make(map[string]*Link) // ключ - короткий ID
	mu      sync.RWMutex
	counter = 1000 // стартовое значение для генерации ID
)

func main() {
	// Настройка роутера
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Статические файлы
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// Главная страница
	r.Get("/", indexHandler)

	// Создание ссылки (HTML-форма)
	r.Post("/shorten", shortenHandler)

	// Создание ссылки (API)
	r.Post("/api/v1/shorten", apiShortenHandler)

	// Редирект по короткой ссылке
	r.Get("/{shortID}", redirectHandler)

	// Страница со всеми ссылками
	r.Get("/links", linksHandler)

	// Запуск сервера
	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Println("📍 Откройте в браузере: http://localhost:8080")
	log.Println("📋 Для создания ссылки используйте форму на главной странице")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// indexHandler — главная страница
// indexHandler — главная страница со списком ссылок
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Mini URL — Сократитель ссылок</title>
		<style>
			* { margin: 0; padding: 0; box-sizing: border-box; }
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
				background: #f5f7fa;
				color: #1a1a2e;
				min-height: 100vh;
			}
			.container { max-width: 900px; margin: 0 auto; padding: 20px; }
			.header {
				background: #ffffff;
				padding: 16px 0;
				border-bottom: 1px solid #e1e5eb;
				box-shadow: 0 2px 4px rgba(0,0,0,0.05);
			}
			.header .container { display: flex; justify-content: space-between; align-items: center; }
			.logo { font-size: 24px; font-weight: 700; color: #4a6cf7; }
			.nav a {
				color: #4a5568;
				text-decoration: none;
				margin-left: 20px;
				font-weight: 500;
				padding: 6px 12px;
				border-radius: 6px;
			}
			.nav a:hover { color: #4a6cf7; background: #eef2ff; }
			.hero { text-align: center; padding: 40px 0 30px; }
			.hero h1 { font-size: 38px; font-weight: 800; background: linear-gradient(135deg, #4a6cf7, #6c63ff); -webkit-background-clip: text; -webkit-text-fill-color: transparent; margin-bottom: 12px; }
			.hero p { font-size: 18px; color: #718096; }
			.card {
				background: #ffffff;
				border-radius: 16px;
				padding: 32px;
				box-shadow: 0 4px 20px rgba(0,0,0,0.06);
				border: 1px solid #edf2f7;
				max-width: 700px;
				margin: 0 auto 30px;
			}
			.card h2 { margin-bottom: 20px; }
			.form-group { margin-bottom: 16px; }
			.form-control {
				width: 100%;
				padding: 14px 18px;
				border: 2px solid #e2e8f0;
				border-radius: 12px;
				font-size: 16px;
				background: #f7fafc;
				transition: border-color 0.2s;
			}
			.form-control:focus { border-color: #4a6cf7; outline: none; background: #ffffff; }
			.btn {
				display: inline-block;
				padding: 12px 28px;
				border: none;
				border-radius: 12px;
				font-size: 16px;
				font-weight: 600;
				cursor: pointer;
				transition: all 0.2s;
			}
			.btn-primary { background: linear-gradient(135deg, #4a6cf7, #6c63ff); color: #ffffff; }
			.btn-primary:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(74,108,247,0.35); }
			.btn-secondary { background: #edf2f7; color: #2d3748; }
			.btn-secondary:hover { background: #e2e8f0; }
			.result {
				margin-top: 20px;
				padding: 20px;
				background: #f0fff4;
				border: 2px solid #68d391;
				border-radius: 12px;
				display: none;
			}
			.result.show { display: block; }
			.short-url {
				font-size: 18px;
				font-weight: 600;
				color: #4a6cf7;
				word-break: break-all;
			}
			.links-section { margin-top: 30px; }
			.links-section h2 { margin-bottom: 16px; }
			.link-card {
				background: #ffffff;
				border-radius: 12px;
				padding: 16px 20px;
				border: 1px solid #edf2f7;
				margin-bottom: 12px;
				transition: all 0.2s;
				display: flex;
				justify-content: space-between;
				align-items: center;
				flex-wrap: wrap;
				gap: 12px;
			}
			.link-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.08); border-color: #4a6cf7; }
			.link-original { color: #2d3748; font-size: 15px; max-width: 40%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
			.link-short { color: #4a6cf7; font-weight: 600; }
			.link-clicks { color: #718096; font-size: 14px; }
			.empty { text-align: center; color: #718096; font-size: 16px; padding: 20px 0; }
			.info-text { text-align: center; color: #718096; font-size: 14px; margin-top: 10px; }
			.footer { text-align: center; padding: 24px 0; color: #718096; font-size: 14px; border-top: 1px solid #e1e5eb; margin-top: 40px; }
		</style>
	</head>
	<body>
		<header class="header">
			<div class="container">
				<div class="logo">🔗 Mini URL</div>
				<nav class="nav">
					<a href="/" style="color:#4a6cf7; background:#eef2ff;">Главная</a>
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
					<div class="form-group">
						<input type="url" id="original-url" name="url" class="form-control" placeholder="Введите длинный URL..." required>
					</div>
					<button type="submit" class="btn btn-primary">Сократить</button>
				</form>
				<div id="result" class="result">
					<p>✅ Ваша короткая ссылка:</p>
					<div style="display:flex; gap:12px; align-items:center; margin-top:8px;">
						<span id="short-url" class="short-url"></span>
						<button onclick="copyToClipboard()" class="btn btn-secondary">📋 Копировать</button>
					</div>
				</div>
			</div>

			<div class="links-section">
				<h2>📋 Последние ссылки</h2>
				{{if .}}
					{{range .}}
					<div class="link-card">
						<span class="link-original" title="{{.OriginalURL}}">{{.OriginalURL}}</span>
						<a href="/{{.ShortID}}" target="_blank" class="link-short">http://localhost:8080/{{.ShortID}}</a>
						<span class="link-clicks">👁 {{.Clicks}} переходов</span>
					</div>
					{{end}}
				{{else}}
					<div class="empty">
						<p>😕 Нет ссылок. Создайте свою первую ссылку!</p>
					</div>
				{{end}}
			</div>

			<footer class="footer">
				<p>&copy; 2026 Mini URL. Все права защищены.</p>
			</footer>
		</div>

		<script>
			// Отображение результата после создания
			document.getElementById('shorten-form').addEventListener('submit', function(e) {
				e.preventDefault();
				const url = document.getElementById('original-url').value;
				fetch('/api/v1/shorten', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ url: url })
				})
				.then(res => res.json())
				.then(data => {
					if (data.short_url) {
						document.getElementById('short-url').textContent = data.short_url;
						document.getElementById('result').classList.add('show');
						document.getElementById('original-url').value = '';
						// Обновить список ссылок через 1 секунду
						setTimeout(() => location.reload(), 1000);
					} else {
						alert('Ошибка: ' + (data.error || 'Неизвестная ошибка'));
					}
				})
				.catch(err => alert('Ошибка: ' + err.message));
			});

			function copyToClipboard() {
				const text = document.getElementById('short-url').textContent;
				navigator.clipboard.writeText(text).then(() => {
					const btn = event.target;
					btn.textContent = '✅ Скопировано!';
					setTimeout(() => { btn.textContent = '📋 Копировать'; }, 2000);
				});
			}
		</script>
	</body>
	</html>
	`
	t, _ := template.New("index").Parse(tmpl)
	mu.RLock()
	var list []*Link
	for _, l := range links {
		list = append(list, l)
	}
	mu.RUnlock()
	t.Execute(w, list)
}

// shortenHandler — создание ссылки через HTML-форму
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL обязателен", http.StatusBadRequest)
		return
	}

	mu.Lock()
	counter++
	shortID := generateShortID(counter)
	link := &Link{
		ID:          shortID,
		OriginalURL: url,
		ShortID:     shortID,
		CreatedAt:   time.Now(),
		Clicks:      0,
	}
	links[shortID] = link
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// apiShortenHandler — создание ссылки через API (JSON)
func apiShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL обязателен", http.StatusBadRequest)
		return
	}

	mu.Lock()
	counter++
	shortID := generateShortID(counter)
	link := &Link{
		ID:          shortID,
		OriginalURL: req.URL,
		ShortID:     shortID,
		CreatedAt:   time.Now(),
		Clicks:      0,
	}
	links[shortID] = link
	mu.Unlock()

	response := map[string]string{
		"short_url": "http://localhost:8080/" + shortID,
		"original":  req.URL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// redirectHandler — редирект по короткой ссылке
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

// linksHandler — возвращает все ссылки в формате JSON
// linksHandler — показывает все ссылки в виде HTML-страницы
func linksHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	var list []*Link
	for _, l := range links {
		list = append(list, l)
	}
	mu.RUnlock()

	// Создаём HTML-шаблон для отображения ссылок
	tmpl := `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Мои ссылки — Mini URL</title>
		<style>
			* { margin: 0; padding: 0; box-sizing: border-box; }
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
				background: #f5f7fa;
				color: #1a1a2e;
				min-height: 100vh;
			}
			.container { max-width: 900px; margin: 0 auto; padding: 20px; }
			.header {
				background: #ffffff;
				padding: 16px 0;
				border-bottom: 1px solid #e1e5eb;
				box-shadow: 0 2px 4px rgba(0,0,0,0.05);
			}
			.header .container { display: flex; justify-content: space-between; align-items: center; }
			.logo { font-size: 24px; font-weight: 700; color: #4a6cf7; }
			.nav a {
				color: #4a5568;
				text-decoration: none;
				margin-left: 20px;
				font-weight: 500;
				padding: 6px 12px;
				border-radius: 6px;
			}
			.nav a:hover { color: #4a6cf7; background: #eef2ff; }
			h1 { margin: 30px 0 20px; }
			.link-card {
				background: #ffffff;
				border-radius: 12px;
				padding: 16px 20px;
				border: 1px solid #edf2f7;
				margin-bottom: 12px;
				transition: all 0.2s;
				display: flex;
				justify-content: space-between;
				align-items: center;
				flex-wrap: wrap;
				gap: 12px;
			}
			.link-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.08); border-color: #4a6cf7; }
			.link-original { color: #2d3748; font-size: 15px; max-width: 40%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
			.link-short { color: #4a6cf7; font-weight: 600; }
			.link-clicks { color: #718096; font-size: 14px; }
			.link-date { color: #a0aec0; font-size: 13px; }
			.empty { text-align: center; color: #718096; font-size: 18px; padding: 40px 0; }
			.btn {
				display: inline-block;
				padding: 10px 24px;
				border: none;
				border-radius: 12px;
				font-size: 14px;
				font-weight: 600;
				cursor: pointer;
				text-decoration: none;
				transition: all 0.2s;
			}
			.btn-primary { background: linear-gradient(135deg, #4a6cf7, #6c63ff); color: #ffffff; }
			.btn-primary:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(74,108,247,0.35); }
			.footer { text-align: center; padding: 24px 0; color: #718096; font-size: 14px; border-top: 1px solid #e1e5eb; margin-top: 40px; }
		</style>
	</head>
	<body>
		<header class="header">
			<div class="container">
				<div class="logo">🔗 Mini URL</div>
				<nav class="nav">
					<a href="/">Главная</a>
					<a href="/links" style="color:#4a6cf7; background:#eef2ff;">Мои ссылки</a>
				</nav>
			</div>
		</header>

		<div class="container">
			<h1>📋 Мои ссылки</h1>
			<a href="/" class="btn btn-primary" style="margin-bottom:20px;">+ Создать новую</a>

			{{if .}}
				{{range .}}
				<div class="link-card">
					<span class="link-original" title="{{.OriginalURL}}">{{.OriginalURL}}</span>
					<a href="/{{.ShortID}}" target="_blank" class="link-short">http://localhost:8080/{{.ShortID}}</a>
					<span class="link-clicks">👁 {{.Clicks}} переходов</span>
					<span class="link-date">📅 {{.CreatedAt.Format "02.01.2006 15:04"}}</span>
				</div>
				{{end}}
			{{else}}
				<div class="empty">
					<p>😕 У вас пока нет ссылок</p>
					<p style="font-size:14px; margin-top:10px;">Создайте свою первую ссылку на <a href="/" style="color:#4a6cf7;">главной странице</a></p>
				</div>
			{{end}}
		</div>

		<footer class="footer">
			<p>&copy; 2026 Mini URL. Все права защищены.</p>
		</footer>
	</body>
	</html>
	`
	t, _ := template.New("links").Parse(tmpl)
	t.Execute(w, list)
}

// generateShortID — генерирует короткий ID на основе числа
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