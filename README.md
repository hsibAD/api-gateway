# API Gateway

This service acts as the API Gateway for the Online Grocery Store microservices, providing a unified interface for clients to interact with the Order and Payment services.

## Features

- Single entry point for all client requests
- Request routing to appropriate microservices
- JWT Authentication
- Rate limiting
- Request/Response transformation
- Error handling
- Logging and monitoring
- Swagger documentation

## Tech Stack

- Go 1.21+
- Gin Web Framework
- JWT Authentication
- Redis (Rate limiting)
- gRPC (Communication with microservices)
- Prometheus metrics
- Swagger/OpenAPI

## Project Structure

```
api-gateway/
├── cmd/                    # Application entry points
├── internal/              
│   ├── auth/              # Authentication middleware
│   ├── config/            # Configuration
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # Custom middleware
│   ├── proxy/             # gRPC client proxies
│   └── server/            # Server setup
├── pkg/                   # Public packages
├── docs/                  # Documentation
├── config/               # Configuration files
└── test/                 # Integration tests
```

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables:
```bash
cp .env.example .env
```

3. Start the service:
```bash
make run
```

## API Documentation

See `/docs/swagger.yaml` for the complete API specification.

## Environment Variables

- `PORT` - Server port (default: 8080)
- `ORDER_SERVICE_URL` - Order service gRPC URL
- `PAYMENT_SERVICE_URL` - Payment service gRPC URL
- `JWT_SECRET` - JWT signing secret
- `REDIS_URL` - Redis URL for rate limiting
- `RATE_LIMIT` - Requests per minute per IP

## License

MIT 