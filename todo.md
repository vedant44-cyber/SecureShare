# Development Checklist

### Phase 1: Project Setup (Infrastructure)

* [ ] Initialize Git Repository: Create a new repo `secure-share` with `.gitignore`.
* [ ] Directory Structure: Create folders:

  * `backend/` (Go)
  * `frontend/` (React)
  * `infrastructure/` (Optional, or root for Docker)
* [ ] Docker Compose: Create `docker-compose.yml` defining services:

  * MinIO (S3 compatible storage)
  * Redis (Metadata & Atomic counters)
* [ ] Run Infrastructure: `docker-compose up -d` and verify:

  * MinIO Console accessible at localhost:9001
  * Redis connection available at localhost:6379
* [ ] MinIO Setup: Create bucket `secure-share` using MinIO console.

---

### Phase 2: Backend Development (Go)

* [ ] Module Init: run `go mod init` and install dependencies (`gin`, `aws-sdk-go-v2`, `go-redis`).
* [ ] Config Layer: create `internal/config/config.go` to load Env Variables (Redis Addr, S3 Keys, etc.).

#### Storage Layer (S3)

* [ ] Implement `NewS3Store()` to connect to MinIO/S3.
* [ ] Implement `Upload(key, stream)` using `PutObject`.
* [ ] Implement `Download(key)` using `GetObject`.
* [ ] Implement `Delete(key)`.

#### Database Layer (Redis)

* [ ] Implement `NewRedisStore()`.
* [ ] Implement `SaveFileMetadata(id, filename, mime, limit, ttl)`.

  * Tip: Use Redis Pipeline to set the Hash and the Limit Key atomically.
* [ ] Implement `DecrementDownloadCount(id)` using DECR.
* [ ] Implement `GetMetadata(id)`.

#### API Handlers

* [ ] Create `POST /api/upload`:

  * Accept multipart file stream.
  * Generate UUID.
  * Stream to S3.
  * Save metadata to Redis.
  * Return `{ "id": "uuid" }`.

* [ ] Create `GET /api/download/:id`:

  * Check Redis DECR. If <0, return 404.
  * Fetch Metadata (Filename).
  * Stream S3 Body to Response Writer.
  * Critical: Trigger BurnFile goroutine if limit reaches 0.

* [ ] Server Entry: Wire everything inside `main.go` and start Gin server.

* [ ] Test with Postman: Upload & download dummy file.

---

### Phase 3: Frontend Development (React)

* [ ] Scaffold App using Vite (`npm create vite@latest`).
* [ ] UI Layout: Simple centered card with Upload & Download views.

#### Crypto Library (`lib/crypto.js`)

* [ ] Implement `generateKey()` (AES-GCM 256).
* [ ] Implement `encryptFile(file, key)` with IV generation and blob output.
* [ ] Implement `decryptFile(blob, key)` reversing above.
* [ ] Implement `exportKey/importKey` (JWK format).

#### Upload Feature

* [ ] File input + TTL + Limit inputs.
* [ ] On submit: encrypt → POST blob → receive id.
* [ ] Generate share URL `/#id={id}&key={keyString}`.
* [ ] Show "Copy Link" button.

#### Download Feature

* [ ] Parse URL hash params.
* [ ] GET blob via axios (`arraybuffer`).
* [ ] Decrypt and trigger browser save.

---

### Phase 4: Polish & Security

* [ ] Configure CORS for frontend domain only.
* [ ] Add file size limit check in Go.
* [ ] Burn-on-Read test: limit=1 should fail second access.
* [ ] Cleanup code, remove key logs, ensure comments explain Zero Knowledge.

---

This checklist covers the entire Full Stack implementation. Good luck! 🚀
