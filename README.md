# Fullstack Engineer Assessment

Task Management Application untuk assessment Fullstack Engineer.

## Tech Stack

### Backend
- Go
- Gin
- MySQL
- Redis

### Frontend
- React Native
- Expo
- TypeScript
- Jest
- React Native Testing Library

### Infrastructure
- Docker
- Docker Compose

---

## Project Structure

```text
fullstack-assessment/
├── backend/
│   ├── cmd/server/
│   ├── internal/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   ├── migrations/
│   ├── .env.example
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── __tests__/
│   │   ├── components/
│   │   └── services/
│   ├── App.tsx
│   ├── package.json
│   └── jest.config.js
│
├── docker-compose.yml
├── Makefile
└── README.md