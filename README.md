# Bookstore API (Golang + Gin + PostgreSQL + JWT)

Mini Project sesuai brief:
- API Buku & Kategori
- Autentikasi JWT (endpoint login: `POST /api/users/login`)
- Validasi rilis tahun (1980–2024), thickness dari `total_page`
- Deployment siap untuk Railway (pakai `DATABASE_URL`)
- Migration via `sql-migrate`
- Dokumentasi endpoint di bawah

## Tech
- Golang + Gin
- PostgreSQL (conn string dari env `DATABASE_URL`)
- JWT (`JWT_SECRET`)
- Migration: `sql-migrate`

## Struktur
```
bookstore-api/
├─ main.go
├─ routes/routes.go
├─ config/config.go
├─ controllers/
├─ middlewares/
├─ models/
├─ migrations/
├─ dbconfig.yml
├─ .env.example
└─ README.md
```

## Env
Buat file `.env` saat lokal:
```
DATABASE_URL=postgres://user:password@localhost:5432/bookstore?sslmode=disable
JWT_SECRET=supersecretjwt
PORT=8080
```

> **Railway**: tidak perlu `.env`, cukup set **Variables**: `DATABASE_URL` & `JWT_SECRET`.

## Migration
Install tool:
```
go install github.com/rubenv/sql-migrate/...@latest
```
Jalankan migration (pastikan `DATABASE_URL` sudah ada di env):
```
sql-migrate up -config=dbconfig.yml -env=production
```
Rollback:
```
sql-migrate down -config=dbconfig.yml -env=production
```

## Seed User (opsional)
Daftar user dengan endpoint `POST /api/users/register`:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

## Run
```
go mod tidy
go run main.go
```
Server jalan di `http://localhost:${PORT}` (default 8080).

## Auth Flow
1. `POST /api/users/register` (opsional untuk buat user)
2. `POST /api/users/login` → dapatkan `token`
3. Set `Authorization: Bearer <token>` pada request endpoint kategori/buku

## Endpoints

### Auth
- `POST /api/users/register` – registrasi user (opsional)
- `POST /api/users/login` – login, response JWT

### Categories (JWT protected)
- `GET /api/categories` – list semua kategori
- `POST /api/categories` – tambah kategori
- `GET /api/categories/:id` – detail kategori
- `PUT /api/categories/:id` – update kategori (pesan error jika id tidak ada)
- `DELETE /api/categories/:id` – hapus kategori (pesan error jika id tidak ada)
- `GET /api/categories/:id/books` – list buku berdasar kategori

### Books (JWT protected)
- `GET /api/books` – list semua buku
- `POST /api/books` – tambah buku (server set `thickness`)
- `GET /api/books/:id` – detail buku
- `PUT /api/books/:id` – update buku (error jika id tidak ada)
- `DELETE /api/books/:id` – hapus buku (error jika id tidak ada)

### Validasi Buku
- `release_year` wajib antara **1980** dan **2024**
- `thickness` otomatis:
  - `total_page > 100` → `"tebal"`
  - `total_page < 100` → `"tipis"`
  - `total_page == 100` → `"tebal"` (asumsi >=100 tebal)

## Deploy ke Railway
1. Push repo ke GitHub.
2. Di Railway: **New Project** → **Deploy from GitHub repo**.
3. Tambahkan Variables:
   - `DATABASE_URL` (sudah otomatis jika pakai Postgres plugin)
   - `JWT_SECRET`
   - `PORT` (opsional, default 8080)
4. Jalankan migration:
   - Buka deploy **Shell** →
     ```
     sql-migrate up -config=dbconfig.yml -env=production
     ```

## Lisensi
MIT
