
<div align="center">
  <h1>⚡ TRACE</h1>
  <p>
    <strong>Secure. Ephemeral. Cyberpunk.</strong>
  </p>
  <p>
    A modern, high-performance pastebin alternative designed for developers who care about aesthetics and privacy.
  </p>

  <p>
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/Backend-Go-00ADD8?style=flat-square&logo=go" alt="Go"></a>
    <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/Frontend-Vue_3-4FC08D?style=flat-square&logo=vue.js" alt="Vue"></a>
    <a href="https://www.postgresql.org/"><img src="https://img.shields.io/badge/Database-PostgreSQL-316192?style=flat-square&logo=postgresql" alt="Postgres"></a>
    <a href="https://www.docker.com/"><img src="https://img.shields.io/badge/Deploy-Docker-2496ED?style=flat-square&logo=docker" alt="Docker"></a>
    <img src="https://img.shields.io/badge/License-MIT-yellow?style=flat-square" alt="License">
  </p>

  <img src="https://placehold.co/1200x600/09090b/6366f1?text=App+Screenshot+Here" alt="Trace App Preview" width="100%" style="border-radius: 10px; border: 1px solid #333;">
</div>

<br>

## ✨ Features

Trace is not just a notepad. It's a fully-featured tool for secure code sharing.

- **💎 Glassmorphism UI** — A stunning dark-mode interface with ambient glow and blur effects.
- **⚡ Smart Detection** — Automatically detects the programming language using `flourite`.
- **🔒 Secure by Design** — Optional password protection (bcrypt) and encrypted transport.
- **🔥 Burn After Reading** — Create self-destructing links that vanish after the first view.
- **⌨️ Keyboard First** — Built-in Command Palette (`Cmd+K`) for power users.
- **👁️ Markdown Preview** — Render documentation and notes instantly with XSS sanitization.
- **🚀 Production Ready** — Dockerized, Rate-limited, and optimized for high loads.

## 🛠️ Tech Stack

This project follows **Clean Architecture** principles.

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Backend** | **Go (Golang)** | Echo framework, Clean Arch (Domain/UseCase/Repo). |
| **Frontend** | **Vue 3** | Composition API, TypeScript, Vite. |
| **Editor** | **CodeMirror 6** | High-performance code editor. |
| **Database** | **PostgreSQL** | Reliable storage with `pgx` driver. |
| **Styling** | **CSS3** | Custom properties, Glassmorphism, Responsive. |
| **Infra** | **Docker** | Multi-stage builds, Nginx reverse proxy. |

## 🚀 Getting Started

### Prerequisites

- Docker & Docker Compose
- *Or* Go 1.24+ and Node.js 20+ (for local dev)

### 🐳 Fast Run (Docker)

The easiest way to run Trace is using Docker Compose.

```bash
# 1. Clone the repository
git clone https://github.com/YOUR_USERNAME/trace.git
cd trace

# 2. Run in Production mode (Nginx + App + DB)
docker compose up -d --build

# 3. Visit http://localhost:3000
```

### 💻 Local Development

If you want to contribute or modify the code:

```bash
# 1. Start Infrastructure (DB + Backend)
docker compose -f docker-compose.dev.yml up --build

# 2. Start Frontend (Hot Reload)
cd web
npm install
npm run dev

# 3. Visit http://localhost:5173
```

## 🔌 API Reference

Trace exposes a simple REST API.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/paste` | Create a new paste. |
| `GET` | `/api/v1/paste/:id` | Retrieve a paste (JSON). |
| `GET` | `/raw/:id` | Retrieve raw content (Plain Text). |

**Example using cURL:**

```bash
# Create a paste
curl -X POST http://localhost:8080/api/v1/paste \
  -H "Content-Type: application/json" \
  -d '{"content": "fmt.Println(\"Hello\")", "language": "go"}'

# Get raw content
curl http://localhost:8080/raw/<PASTE_ID>
```

## 🎹 Keyboard Shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Cmd + K` / `Ctrl + K` | Open Command Palette |
| `Cmd + S` / `Ctrl + S` | Save Paste |
| `Esc` | Close Modals |

## 🛡️ Security

- **Rate Limiting:** Built-in protection against DDoS/Spam (Token Bucket algorithm).
- **Sanitization:** DOMPurify is used for Markdown rendering to prevent XSS.
- **Isolation:** Backend is not exposed directly in production; traffic goes through Nginx.


---

<div align="center">
  <sub>Built with ❤️ by <a href="https://github.com/shadowkyst">ShadowKyst</a></sub>
</div>
