# Personal Financial Management System

A robust backend API for managing personal finances, tracking income and expenses, and setting budgets. Built with Go, Gin, and PostgreSQL.

## 🚀 Features

- **User Management**: Secure user registration, login, and password updates using JWT authentication.
- **Financial Records**: Full CRUD operations for financial entries (Income/Expense).
- **Summaries**: Get financial summaries by current month, year, or custom ranges.
- **Budgeting**: Set and track budgets, view budget history, and check usage status.
- **Containerized**: Fully dockerized application and database setup.

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Framework**: [Gin](https://gin-gonic.com/)
- **Database**: PostgreSQL
- **Database Driver**: pgx
- **SQL Generation**: [sqlc](https://sqlc.dev/)
- **Migration**: Golang Migrate
- **Configuration**: Viper
- **Authentication**: JWT (JSON Web Tokens)
- **Testing**: Testify, Mockgen

## 📂 Project Structure

- `api/`: HTTP handlers and route definitions.
- `db/`: Database migrations (`migration/`) and generated SQLC code (`sqlc/`, `query/`).
- `token/`: JWT token generation and validation.
- `util/`: Configuration and utility functions.
- `main.go`: Application entry point.

## 🔧 Prerequisites

- **Go** 1.24+
- **Docker** & **Docker Compose**
- **Make** (optional, for running convenience scripts)

## ⚙️ Configuration

The application expects an `app.env` file in the root directory (or environment variables) with the following keys:

```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5433/personal_financial?sslmode=disable
SERVER_PORT=8088
```

## 📦 Getting Started

### Using Docker Compose (Recommended)

The easiest way to run the entire system (Database + API + Migrations) is via Docker Compose.

```bash
docker-compose up --build -d
```

This will start:

- PostgreSQL on port `5433`
- The API server on port `8088`

### Local Development

If you prefer to run the Go application locally:

1.  **Start the Database**:
    You can use the provided `docker-compose` purely for the DB, or run a standalone container.

    ```bash
    docker run --name postgres -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15
    ```

2.  **Create Database**:

    ```bash
    make createdb
    ```

3.  **Run Migrations**:

    ```bash
    make migrateup
    ```

4.  **Install Dependencies**:

    ```bash
    go mod tidy
    ```

5.  **Run the Server**:
    ```bash
    make server
    # OR
    go run main.go
    ```

## 📡 API Endpoints

### Auth

- `POST /create-user`: Register a new user.
- `POST /login-user`: Login and receive access token.
- `PUT /update-password`: Change user password.

### Financials

- `POST /new-financial`: Add a new income/expense record.
- `GET /my-financial`: List all financial records.
- `GET /financial/get/:id`: Get a specific record.
- `PUT /financial/update/:id`: Update a record.
- `DELETE /financial/delete/:id`: Delete a record.

### Summary

- `GET /summary/current-month`: Summary for the current month.
- `GET /summary/current-year`: Summary for the current year.
- `GET /summary/each-year`: Summary broken down by year.
- `GET /summary/month`: Summary by specific month/year.

### Budget

- `POST /budget/`: Set a new budget.
- `GET /budget/`: Get current budget.
- `GET /budget/check`: Check if budget is exceeded.
- `GET /budget/history`: View budget history.

## 🧪 Testing

Run internal tests using:

```bash
make test
```

## 📜 License

[MIT](LICENSE)
