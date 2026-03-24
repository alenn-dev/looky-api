# Looky API

Real-time order tracking API built with Go. Demonstrates JWT auth with roles, WebSockets, Kafka, GraphQL, and Docker.

## Tech Stack

- **Go** + **Fiber v3** — HTTP framework
- **GORM** + **PostgreSQL** — database
- **JWT** — authentication with roles
- **WebSockets** — real-time order status updates
- **Kafka** — event streaming between services
- **GraphQL** (gqlgen) — order history queries
- **Docker Compose** — full stack orchestration

## Project Structure
```
main.go
internal/
  config/         # Environment variables
  database/       # DB connection
  handlers/       # HTTP + WebSocket + GraphQL handlers
  middleware/      # Auth & role middleware
  models/         # GORM models & DTOs
  routes/         # Route definitions
  services/       # Business logic
  kafka/          # Producer & consumer
  graphql/        # Schema, resolvers, generated code
  utils/          # JWT helpers
```

## Roles

| Role | Permissions |
|------|-------------|
| `customer` | Create orders, cancel pending orders, view own orders |
| `driver` | Update order to on_the_way and delivered |
| `restaurant_owner` | Manage restaurants and products, confirm/prepare/cancel orders |

## Order Flow
```
pending → confirmed → preparing → on_the_way → delivered
                ↓
           cancelled
```

## Requirements

- [Go 1.25+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) and Docker Compose

## Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/your-username/looky-api.git
cd looky-api
```

### 2. Set up environment variables
```bash
cp .env.example .env
```

Edit `.env` with your values:
```env
DB_URL="host=localhost user=admin password=admin123 dbname=looky_api port=5432 sslmode=disable"
JWT_SECRET="your_jwt_secret_here"
KAFKA_BROKERS="localhost:9092"
```

Generate a secure JWT secret:
```bash
openssl rand -hex 32
```

### 3. Start the stack
```bash
docker compose up -d
```

### 4. Run the API
```bash
go run main.go
```

The API will be available at `http://localhost:3000`.

### Or run everything with Docker
```bash
docker compose up -d
```

## API Endpoints

### Auth

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/register` | No | Register user |
| POST | `/api/auth/login` | No | Login and get token |

### Restaurants

| Method | Endpoint | Auth | Role |
|--------|----------|------|------|
| GET | `/api/restaurants` | Yes | Any |
| GET | `/api/restaurants/:id` | Yes | Any |
| POST | `/api/restaurants` | Yes | restaurant_owner |
| PATCH | `/api/restaurants/:id` | Yes | restaurant_owner |
| DELETE | `/api/restaurants/:id` | Yes | restaurant_owner |

### Products

| Method | Endpoint | Auth | Role |
|--------|----------|------|------|
| GET | `/api/restaurants/:restaurantId/products` | Yes | Any |
| GET | `/api/restaurants/:restaurantId/products/:id` | Yes | Any |
| POST | `/api/restaurants/:restaurantId/products` | Yes | restaurant_owner |
| PATCH | `/api/restaurants/:restaurantId/products/:id` | Yes | restaurant_owner |
| DELETE | `/api/restaurants/:restaurantId/products/:id` | Yes | restaurant_owner |

### Orders

| Method | Endpoint | Auth | Role |
|--------|----------|------|------|
| GET | `/api/orders` | Yes | Any |
| GET | `/api/orders/:id` | Yes | Any |
| POST | `/api/orders` | Yes | customer |
| PATCH | `/api/orders/:id/status` | Yes | Any |

### GraphQL

| Endpoint | Description |
|----------|-------------|
| POST `/graphql` | GraphQL queries |
| GET `/playground` | GraphQL playground (dev only) |

#### Example queries
```graphql
# Get all delivered orders
query {
  orders(status: delivered) {
    id
    total
    status
    items {
      name
      quantity
      unitPrice
    }
  }
}

# Order status history
query {
  orderHistory(orderId: "your-order-id") {
    status
    changedAt
  }
}
```

### WebSocket

Connect to receive real-time order updates:
```
ws://localhost:3000/ws?token=your_jwt_token
```

You will receive messages like:
```json
{"order_id": "uuid", "status": "on_the_way"}
```

## Monitoring

- **Kafka UI**: `http://localhost:8080` — view topics and messages in real time

## Stopping
```bash
docker compose down
```

Remove all data:
```bash
docker compose down -v
```