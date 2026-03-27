# Football API - Technical Test

REST API backend untuk manajemen tim sepak bola amatir XYZ menggunakan Golang, Gin, PostgreSQL, dan Redis (optional runtime support).

## Executive Summary
Project ini adalah backend API production minded untuk kebutuhan aplikasi Android internal perusahaan. Fokus implementasi:
- Arsitektur modular berbasis domain dengan separation of concerns yang jelas.
- Security baseline yang kuat (auth, RBAC, validation, rate limit, secure headers).
- Data integrity (soft delete, constraint DB, transaction, immutable result).
- API contract konsisten dengan response envelope standar.

## Scope Yang Dikerjakan
Sesuai requirement test, backend ini meng-cover:
1. Pengelolaan tim (`teams`): create, read, update, soft delete.
2. Pengelolaan pemain (`players`) dengan relasi 1 tim : banyak pemain.
3. Pengelolaan jadwal pertandingan (`matches`).
4. Submit hasil pertandingan + event gol.
5. Report hasil pertandingan (status akhir, top scorer, akumulasi kemenangan).
6. Authentication & authorization (`admin`, `staff`, `viewer`).
7. Security dan kualitas API untuk skenario nyata.

## Arsitektur Singkat
- Gaya arsitektur: `Feature-Based Modular Architecture` + `Clean Architecture (pragmatic)`.
- Versioning API: `/api/v1`.
- Pola domain: `controller -> service -> repository -> model`.
- Controller hanya HTTP concern.
- Service memegang business rules.
- Repository isolasi akses DB.

## Alasan Pemilihan JWT
JWT dipilih karena cocok untuk REST API Android yang stateless.

Kelebihan:
- Stateless dan mudah scale horizontal.
- Praktis untuk mobile bearer-token flow.
- Tidak butuh session storage server untuk request-level auth.

Tradeoff:
- Revocation tidak semudah session tradisional.
- Harus disiplin expiry policy (access pendek + refresh).
- Payload token tidak boleh berisi data sensitif.

## Struktur Project
```text
.
├── cmd/
│   └── api/
│       └── main.go                        # bootstrap app, middleware, central route register
├── internal/
│   ├── core/
│   │   ├── config/
│   │   │   └── config.go                  # env loader + default config
│   │   ├── database/
│   │   │   └── database.go                # DB init, pool setup, automigrate, seed
│   │   ├── errors/
│   │   │   └── app_error.go               # shared app error type
│   │   ├── logger/
│   │   │   └── logger.go                  # logger abstraction
│   │   └── middleware/
│   │       ├── auth.go                    # JWT auth + RBAC role middleware
│   │       ├── rate_limit.go              # login rate limiter
│   │       ├── recovery.go                # panic recovery + standardized error
│   │       ├── request_id.go              # request trace id middleware
│   │       └── security.go                # CORS + secure headers
│   ├── domains/
│   │   ├── auth/
│   │   │   ├── controllers/auth_controller.go
│   │   │   ├── dtos/auth_dto.go
│   │   │   ├── repositories/auth_repository.go
│   │   │   ├── routes/routes.go
│   │   │   └── services/auth_service.go
│   │   ├── users/
│   │   │   ├── controllers/user_controller.go
│   │   │   ├── dtos/user_dto.go
│   │   │   ├── models/user.go
│   │   │   ├── repositories/user_repository.go
│   │   │   ├── routes/routes.go
│   │   │   └── services/user_service.go
│   │   ├── teams/
│   │   │   ├── controllers/team_controller.go
│   │   │   ├── dtos/team_dto.go
│   │   │   ├── models/team.go
│   │   │   ├── repositories/team_repository.go
│   │   │   ├── routes/routes.go
│   │   │   └── services/team_service.go
│   │   ├── players/
│   │   │   ├── controllers/player_controller.go
│   │   │   ├── dtos/player_dto.go
│   │   │   ├── models/player.go
│   │   │   ├── repositories/player_repository.go
│   │   │   ├── routes/routes.go
│   │   │   └── services/player_service.go
│   │   ├── matches/
│   │   │   ├── controllers/match_controller.go
│   │   │   ├── dtos/match_dto.go
│   │   │   ├── models/match.go
│   │   │   ├── repositories/match_repository.go
│   │   │   ├── routes/routes.go
│   │   │   └── services/match_service.go
│   │   └── reports/
│   │       ├── controllers/report_controller.go
│   │       ├── dtos/report_dto.go
│   │       ├── repositories/report_repository.go
│   │       ├── routes/routes.go
│   │       └── services/report_service.go
│   └── shared/
│       ├── constants/enums.go             # role, status match, player position
│       ├── pagination/pagination.go       # pagination parser + helper
│       ├── response/response.go           # response envelope standard
│       └── validator/validator.go         # global request validation
├── migrations/
│   └── 001_init.sql                       # reference SQL schema + index + FK
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

## Domain Yang Diimplementasikan
- `auth`: login + refresh token.
- `users`: manajemen user untuk role management (`admin/staff/viewer`).
- `teams`: CRUD tim sepak bola.
- `players`: CRUD pemain dengan constraint nomor punggung unik per tim.
- `matches`: jadwal pertandingan + submit result final (immutable).
- `reports`: report agregasi hasil pertandingan + cache sederhana.

## API Endpoint Ringkas
Base URL: `/api/v1`

Auth:
- `POST /auth/login`
- `POST /auth/refresh`

Users (admin only):
- `GET /users`
- `POST /users`

Teams:
- `GET /teams` (admin/staff/viewer)
- `GET /teams/:id` (admin/staff/viewer)
- `POST /teams` (admin)
- `PUT /teams/:id` (admin)
- `DELETE /teams/:id` (admin, soft delete)

Players:
- `GET /players` (admin/staff/viewer)
- `GET /players/:id` (admin/staff/viewer)
- `POST /players` (admin)
- `PUT /players/:id` (admin)
- `DELETE /players/:id` (admin, soft delete)

Matches:
- `GET /matches` (admin/staff/viewer)
- `GET /matches/:id` (admin/staff/viewer)
- `POST /matches` (admin)
- `DELETE /matches/:id` (admin, soft delete)
- `POST /matches/:id/result` (admin, immutable submit)
- `POST /matches/:id/rollback` (admin, optional rollback)

Reports:
- `GET /reports/matches` (admin/staff/viewer)
- `POST /reports/revalidate` (admin)

## Response Envelope Standard
Success:
```json
{
  "success": true,
  "message": "...",
  "data": {}
}
```

Error:
```json
{
  "success": false,
  "message": "...",
  "code": "...",
  "traceId": "..."
}
```

## Data Integrity & Transaction
Implementasi integrity utama:
- Soft delete pada entitas utama (`users`, `teams`, `players`, `matches`).
- Transaction boundary untuk submit result pertandingan:
  - update skor/status match
  - insert seluruh goal events
  - atomic (gagal satu, rollback semua)
- Immutable final result:
  - hasil pertandingan yang sudah submit tidak bisa di-submit ulang.
  - rollback dibatasi endpoint admin.

## Database Design & Constraint
Constraint utama:
- `UNIQUE(team_id, jersey_number)` pada players.
- Foreign key antar entitas untuk menjaga referential integrity.

Index utama:
- `players.team_id`
- `goal_events.match_id`
- plus index status/tanggal di match untuk query list/report.

## Pagination & Filtering
Endpoint list mendukung pagination dan filtering:
- `GET /teams`: `page`, `limit`, `search`, `city`
- `GET /players`: `page`, `limit`, `search`, `position`, `team_id`
- `GET /matches`: `page`, `limit`, `status`, `home_team_id`, `away_team_id`, `date_from`, `date_to`

## Security Coverage
- JWT authentication.
- RBAC authorization middleware.
- Request validation (`validator.v10`).
- Password hashing (`bcrypt`).
- Login rate limiter.
- Login rate limiter configurable: `memory` (default) / `redis` (distributed).
- CORS whitelist.
- Security headers middleware.
- Structured JSON request logging (traceId, method, path, status, latency).
- Global panic recovery.
- Trace ID (`X-Request-ID`) di request dan response error.
- No hardcoded production secret (menggunakan `.env`).

## Threat Modeling Ringkas (STRIDE)
| Threat | Entry Point | Mitigasi |
|---|---|---|
| Spoofing | login/auth header | JWT signature + expiry + rate limit |
| Tampering | submit result | RBAC admin + immutable policy + transaction |
| Repudiation | perubahan data sensitif | trace ID + structured logging style |
| Information Disclosure | API error | generic error message + no stacktrace exposure |
| Denial of Service | auth brute force | login rate limiter |
| Elevation of Privilege | endpoint mutasi | role middleware (server-side authorization) |


## Menjalankan Project
### Opsi A - Docker (disarankan)
```bash
docker compose up --build
```

Health check:
```bash
GET http://localhost:8080/health
```

### Opsi B - Local Go Runtime
1. Copy env:
```bash
cp .env.example .env
```
2. Pastikan PostgreSQL tersedia.
3. Run:
```bash
go run ./cmd/api
```

## Seed Data Default
Auto-seed saat startup:
- Admin:
  - email: `admin@xyz.com`
  - password: `Admin123!`
- Sample teams + players untuk mempermudah testing API.

## Konfigurasi Rate Limit
- Default: `RATE_LIMIT_STORE=memory`.
- Untuk distributed limiter: set `RATE_LIMIT_STORE=redis` dan pastikan `REDIS_ADDR` reachable.

## Demo Flow 
1. `POST /auth/login` -> ambil token.
2. `POST /matches` -> buat jadwal (status `scheduled`).
3. `POST /matches/:id/result` -> submit hasil (status `finished`).
4. Ulang submit result yang sama -> harus gagal immutable.
5. `GET /reports/matches` -> lihat report agregasi hasil.


## Security/Quality Checks
Command yang bisa dijalankan:
```bash
go test ./...
govulncheck ./...
gosec ./...
```

## Dev vs Production
Development:
- Konfigurasi lokal pragmatis untuk kecepatan testing.

Production recommendation:
- HTTPS + HSTS di reverse proxy.
- Secrets via secret manager.
- CORS origins ketat.
- Monitoring + alerting terintegrasi.
- Vulnerability scanning rutin di CI/CD.
