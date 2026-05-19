# FinanceTracker

Personal finance management application for tracking income, expenses, budgets, goals, and generating financial reports.

## Tech Stack

**Frontend** — Vue 3 · TypeScript · Pinia · Tailwind CSS · Vite  
**Backend** — Golang · Gin · GORM · PostgreSQL · JWT  
**Infrastructure** — Docker · Docker Compose

## Features

- JWT authentication with refresh tokens
- Transaction management (income, expense, transfer)
- Account management (checking, savings, credit, cash, investment)
- Category management with custom icons and colors
- Budget tracking with period-based alerts
- Financial goal tracking with contribution history
- Analytics dashboard with spending breakdowns
- Multi-currency support (default: IDR)
- Dark / light theme

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

### Run with Docker

```bash
docker-compose up -d
```

Frontend: http://localhost:5173  
Backend API: http://localhost:8080

### Run manually

**Backend**
```bash
cd backend
cp .env.example .env   # configure your DB credentials
go mod download
go run main.go
```

**Frontend**
```bash
cd frontend
npm install
npm run dev
```

## Project Structure

```
financetracker/
├── backend/          # Golang + Gin REST API
│   ├── controllers/
│   ├── models/
│   ├── middleware/
│   └── config/
├── frontend/         # Vue 3 + TypeScript SPA
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   ├── stores/
│   │   └── services/
└── docker-compose.yml
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/register` | Register |
| GET | `/api/transactions` | List transactions |
| POST | `/api/transactions` | Create transaction |
| GET | `/api/accounts` | List accounts |
| GET | `/api/budgets` | List budgets |
| GET | `/api/goals` | List goals |
| GET | `/api/reports/summary` | Financial summary |
