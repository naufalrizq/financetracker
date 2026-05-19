# FinanceTracker Setup Guide

**Tech Stack: Vue.js 3 + Golang/Gin + PostgreSQL**

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Node.js** (v18 or higher) - [Download](https://nodejs.org/)
- **Go** (v1.21 or higher) - [Download](https://golang.org/dl/)
- **PostgreSQL** (v13 or higher) - [Download](https://www.postgresql.org/download/)
- **Docker & Docker Compose** (optional but recommended) - [Download](https://www.docker.com/)
- **Git** - [Download](https://git-scm.com/)

## 🚀 Quick Start (Docker - Recommended)

### Step 1: Navigate to Project Directory
```bash
cd /home/nrizq/Documents/Codes/Learning/VueGolang/financetracker
```

### Step 2: Start All Services
```bash
# Start PostgreSQL, Backend, and Frontend
docker-compose up -d

# View logs (optional)
docker-compose logs -f
```

### Step 3: Access the Application
- **Frontend**: http://localhost:8080
- **Backend API**: http://localhost:8000
- **API Documentation**: http://localhost:8000/docs

### Step 4: Test with Demo Account
- **Email**: demo@financetracker.com
- **Password**: demo123456

---

## 🛠️ Manual Setup (Development)

### Backend Setup (Golang + Gin)

#### Step 1: Navigate to Backend Directory
```bash
cd backend
```

#### Step 2: Install Dependencies
```bash
# Initialize Go modules
go mod tidy

# Install Air for hot reloading (optional)
go install github.com/cosmtrek/air@latest
```

#### Step 3: Setup Environment Variables
```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your settings
nano .env
```

**Required Environment Variables:**
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=financetracker
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production-minimum-32-characters

# Server Configuration
PORT=8000
GIN_MODE=debug
```

#### Step 4: Setup PostgreSQL Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE financetracker;

# Create user (optional)
CREATE USER financetracker_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE financetracker TO financetracker_user;

# Exit PostgreSQL
\q
```

#### Step 5: Run Backend Server
```bash
# Option 1: With Air (hot reloading)
air

# Option 2: Direct Go run
go run main.go

# Option 3: Build and run
go build -o financetracker
./financetracker
```

**Backend should be running on**: http://localhost:8000

### Frontend Setup (Vue.js 3)

#### Step 1: Navigate to Frontend Directory
```bash
cd ../frontend
```

#### Step 2: Install Dependencies
```bash
# Install npm packages
npm install

# Or using yarn
yarn install
```

#### Step 3: Setup Environment Variables
```bash
# Create environment file
touch .env.local

# Add API URL
echo "VITE_API_URL=http://localhost:8000/api" > .env.local
```

#### Step 4: Run Frontend Development Server
```bash
# Start development server
npm run dev

# Or using yarn
yarn dev
```

**Frontend should be running on**: http://localhost:8080

---

## 🧪 Testing

### Backend Tests
```bash
cd backend

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./controllers -v
```

### Frontend Tests
```bash
cd frontend

# Run unit tests
npm run test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch
```

---

## 📦 Production Build

### Backend Production Build
```bash
cd backend

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Or using Docker
docker build -t financetracker-backend .
```

### Frontend Production Build
```bash
cd frontend

# Build for production
npm run build

# Preview production build
npm run preview

# Or using Docker
docker build -t financetracker-frontend .
```

---

## 🐳 Docker Commands

### Development
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f [service_name]

# Stop all services
docker-compose down

# Rebuild services
docker-compose up --build
```

### Production
```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Start production services
docker-compose -f docker-compose.prod.yml up -d
```

---

## 🔧 Troubleshooting

### Common Issues

#### 1. Database Connection Error
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Start PostgreSQL
sudo systemctl start postgresql

# Check database exists
psql -U postgres -l
```

#### 2. Port Already in Use
```bash
# Check what's using port 8000
lsof -i :8000

# Kill process using port
kill -9 <PID>
```

#### 3. Go Module Issues
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
go mod tidy
```

#### 4. Node.js Issues
```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Environment-Specific Issues

#### Development
- Ensure `.env` file exists in backend directory
- Check if all required environment variables are set
- Verify database connection settings

#### Production
- Use production environment variables
- Ensure proper SSL certificates
- Configure reverse proxy (nginx/apache)

---

## 📚 API Documentation

### Authentication Endpoints
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh JWT token

### Transaction Endpoints
- `GET /api/transactions` - Get transactions with filters
- `POST /api/transactions` - Create transaction
- `PUT /api/transactions/:id` - Update transaction
- `DELETE /api/transactions/:id` - Delete transaction

### Account Endpoints
- `GET /api/accounts` - Get user accounts
- `POST /api/accounts` - Create account
- `PUT /api/accounts/:id` - Update account
- `DELETE /api/accounts/:id` - Delete account

### Full API documentation available at: http://localhost:8000/docs

---

## 🎯 Next Steps

1. **Customize the application** to your needs
2. **Add more features** like recurring transactions, data export
3. **Deploy to production** using Docker or cloud services
4. **Set up monitoring** and logging
5. **Configure backup** for your database

---

## 🆘 Support

If you encounter any issues:

1. Check the troubleshooting section above
2. Review the logs: `docker-compose logs -f`
3. Ensure all prerequisites are installed
4. Verify environment variables are set correctly

**Project Status**: ✅ Production Ready
**Last Updated**: December 2024