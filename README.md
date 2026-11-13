# Product Management API

A RESTful API built with Go and PostgreSQL for managing products.

## Technologies

- Go
- PostgreSQL
- pgx/v5
- net/http
- godotenv

## Features

- Create a product
- Get all products
- Get a product by ID
- Update a product
- Delete a product

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/products` | Create a product |
| GET | `/products` | Get all products |
| GET | `/products/{id}` | Get product by ID |
| PUT | `/products/{id}` | Update a product |
| DELETE | `/products/{id}` | Delete a product |

## Database

Database name: `product_management`

Table: `products`

Fields:
- id
- name
- price
- quantity
- created_at

## Installation

1. Clone the repository:

```bash
git clone YOUR_GITHUB_REPOSITORY_URL
cd product-management-api
```

2. Install dependencies:

```bash
go mod tidy
```

3. Create a `.env` file:

```env
DATABASE_URL=postgres://YOUR_USERNAME:YOUR_PASSWORD@localhost:5432/product_management?sslmode=disable
```

4. Run the project:

```bash
go run .
```

Server runs on:

```text
http://localhost:8080
```

## Author

Go Backend Developer