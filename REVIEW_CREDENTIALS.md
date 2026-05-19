# FinanceTracker — Review Credentials

## Demo Account (works without backend)

| Field    | Value                          |
|----------|--------------------------------|
| Email    | `demo@financetracker.com`      |
| Password | `demo123456`                   |

## How to Login

1. Open the app at `http://localhost:8080`
2. Click **"Use Demo"** button on the login page — it fills the credentials automatically
3. If the Golang backend is not running, the app will automatically fall back to Demo Mode on any login attempt

## Seeded Backend Users (requires `docker-compose up`)

| Role  | Email                        | Password       |
|-------|------------------------------|----------------|
| User  | `demo@financetracker.com`    | `demo123456`   |
| Admin | `admin@financetracker.com`   | `Admin123!`    |

> **Note:** The frontend works standalone in demo mode. Backend users only work when the Golang + PostgreSQL containers are running.
