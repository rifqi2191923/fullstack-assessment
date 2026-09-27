# Fullstack Engineer Assessment

Task Management application assessment using the requested stack:

- Backend: Go + Gin
- Database: MySQL
- Cache: Redis
- Frontend: React Native + TypeScript (Expo)

## Requirements covered
- Task filtering: status, keyword, assignee, page, limit, sort
- PUT `/api/tasks/:id`
- Soft DELETE `/api/tasks/:id`
- Consistent JSON error responses
- Redis cache for GET `/api/tasks` for 60 seconds
- Cache invalidation after create/update/delete
- Search, status filter, pagination, edit modal, loading state
- Duplicate title returns HTTP 409
- Refresh after update
- Soft-deleted tasks hidden
- Backend tests for update/search/cache invalidation
- Frontend component test for search/task list

## Run infrastructure

```bash
docker compose up -d mysql redis
```

MySQL is exposed on `3306`, Redis on `6379`.

If you already run MySQL locally, use the database created by `backend/migrations/001_init.sql` and configure the backend `.env` accordingly.

## Backend

```bash
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/server
```

API: `http://localhost:8080`

## Frontend

```bash
cd frontend
npm install
npx expo start
```

Set `EXPO_PUBLIC_API_URL` if the API is not on the default host.

## API examples

```text
GET    /api/tasks?status=todo&keyword=web&assignee=1&page=1&limit=10&sort=created_at_desc
POST   /api/tasks
PUT    /api/tasks/:id
DELETE /api/tasks/:id
```

## Notes

The database is MySQL as requested. Redis is still included because it is an explicit part of the assessment.
