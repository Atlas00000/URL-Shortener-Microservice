# URL Shortener Development Guide

## Prerequisites

### Required Software
- Go 1.21 or later
- Docker and Docker Compose
- Git
- SQLite 3
- Redis 7.0 or later

### Development Tools
- VS Code or preferred IDE
- Postman or similar API testing tool
- Redis CLI
- SQLite CLI

## Development Environment Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd url-shortener
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Environment Configuration
Create a `.env` file in the project root:
```env
REDIS_URL=localhost:6379
BASE_URL=http://localhost:8080
RATE_LIMIT=60
CACHE_TTL=86400
```

### 4. Start Development Services
```bash
docker-compose up -d redis
```

## Project Structure

```
.
├── main.go                 # Application entry point
├── config/                 # Configuration management
│   └── config.go
├── handlers/              # HTTP request handlers
│   ├── url.go
│   └── analytics.go
├── models/                # Data models
│   ├── url.go
│   └── click.go
├── services/             # Business logic
│   ├── shortener.go
│   └── analytics.go
├── utils/                # Utility functions
│   ├── cache.go
│   └── validator.go
├── static/              # Static assets
│   ├── index.html
│   ├── styles.css
│   └── script.js
├── tests/               # Test files
│   ├── unit/
│   └── integration/
└── docs/               # Documentation
    ├── API.md
    └── ARCHITECTURE.md
```

## Development Workflow

### 1. Running the Application
```bash
go run main.go
```

### 2. Running Tests
```bash
# Run all tests
go test ./...

# Run specific test
go test ./handlers -v

# Run with coverage
go test ./... -cover
```

### 3. Code Style
- Follow Go standard formatting:
  ```bash
  go fmt ./...
  ```
- Use `golint` for code quality:
  ```bash
  golint ./...
  ```

### 4. Database Management

#### SQLite
```bash
# Open SQLite console
sqlite3 urls.db

# Common commands
.tables
.schema urls
.schema clicks
```

#### Redis
```bash
# Open Redis CLI
redis-cli

# Common commands
KEYS *
GET <key>
TTL <key>
```

## API Development

### 1. Adding New Endpoints
1. Create handler in `handlers/` directory
2. Register route in `main.go`
3. Add tests in `tests/` directory
4. Update API documentation

### 2. Testing Endpoints
```bash
# Using curl
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Using Postman
# Import the provided Postman collection
```

## Frontend Development

### 1. HTML/CSS/JavaScript
- Edit files in `static/` directory
- Use modern ES6+ features
- Follow responsive design principles
- Test in multiple browsers

### 2. Asset Management
- Keep assets in `static/` directory
- Use relative paths
- Optimize images and scripts
- Minify production assets

## Testing

### 1. Unit Tests
- Test individual components
- Mock external dependencies
- Use table-driven tests
- Example:
  ```go
  func TestURLShortener(t *testing.T) {
      tests := []struct {
          name     string
          input    string
          expected string
      }{
          // Test cases
      }
      // Test implementation
  }
  ```

### 2. Integration Tests
- Test component interactions
- Use test database
- Clean up after tests
- Example:
  ```go
  func TestURLShorteningFlow(t *testing.T) {
      // Setup
      // Test implementation
      // Cleanup
  }
  ```

### 3. Load Testing
- Use `hey` or similar tools
- Test rate limiting
- Monitor performance
- Example:
  ```bash
  hey -n 1000 -c 50 http://localhost:8080/shorten
  ```

## Debugging

### 1. Logging
- Use structured logging
- Include request IDs
- Log errors with context
- Example:
  ```go
  log.Printf("[%s] Error shortening URL: %v", requestID, err)
  ```

### 2. Debugging Tools
- Use Delve for debugging
- Monitor Redis with Redis Commander
- Use SQLite Browser for database inspection

## Deployment

### 1. Building
```bash
# Build binary
go build -o url-shortener

# Build Docker image
docker build -t url-shortener .
```

### 2. Running in Production
```bash
# Using Docker Compose
docker-compose up -d

# Using binary
./url-shortener
```

## Best Practices

### 1. Code Organization
- Keep functions small and focused
- Use meaningful names
- Add comments for complex logic
- Follow Go idioms

### 2. Error Handling
- Use custom error types
- Add context to errors
- Log errors appropriately
- Return meaningful error messages

### 3. Performance
- Use connection pooling
- Implement caching
- Optimize database queries
- Monitor resource usage

### 4. Security
- Validate all input
- Use prepared statements
- Implement rate limiting
- Follow security headers

## Contributing

### 1. Branch Strategy
- `main`: Production code
- `develop`: Development branch
- Feature branches: `feature/name`
- Bug fix branches: `fix/name`

### 2. Pull Request Process
1. Create feature branch
2. Write tests
3. Update documentation
4. Submit PR
5. Address review comments

### 3. Code Review
- Check code style
- Verify tests
- Review documentation
- Test functionality

## Troubleshooting

### Common Issues

#### 1. Redis Connection
```bash
# Check Redis status
docker-compose ps redis

# Check Redis logs
docker-compose logs redis
```

#### 2. Database Issues
```bash
# Check database file
ls -l urls.db

# Verify database integrity
sqlite3 urls.db "PRAGMA integrity_check;"
```

#### 3. Port Conflicts
```bash
# Check port usage
netstat -an | grep 8080

# Change port in .env
BASE_URL=http://localhost:8081
```

## Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [SQLite Documentation](https://www.sqlite.org/docs.html)

### Tools
- [Postman](https://www.postman.com/)
- [Redis Commander](https://github.com/joeferner/redis-commander)
- [SQLite Browser](https://sqlitebrowser.org/)
- [Delve](https://github.com/go-delve/delve) 