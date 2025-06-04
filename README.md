# URL Shortener Microservice

A modern, scalable URL shortening service built with Go, featuring a clean web interface, analytics tracking, and Redis caching.

## 🌟 Features

- **URL Shortening**: Convert long URLs into short, manageable links
- **Custom Expiration**: Set custom expiration dates for shortened URLs
- **Analytics**: Track clicks, geographic data, and device information
- **Modern UI**: Clean, responsive web interface
- **Caching**: Redis-based caching for improved performance
- **Persistence**: SQLite database for reliable data storage
- **Docker Support**: Easy deployment with Docker and Docker Compose

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Git

### Running the Service

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/urlshortener.git
   cd urlshortener
   ```

2. Start the service:
   ```bash
   docker-compose up --build
   ```

3. Access the web interface at `http://localhost:8080`

## 🏗️ Architecture

### Components

- **Web Interface**: Single-page application for URL shortening
- **API Server**: Go-based REST API
- **Database**: SQLite for persistent storage
- **Cache**: Redis for performance optimization
- **Analytics**: Click tracking and geographic data

### Technology Stack

- **Backend**: Go (Golang)
- **Web Framework**: Gin
- **Database**: SQLite
- **Cache**: Redis
- **Frontend**: HTML, CSS, JavaScript
- **Containerization**: Docker

## 📚 API Documentation

### Endpoints

#### 1. Shorten URL
```http
POST /shorten
Content-Type: application/json

{
    "url": "https://example.com/very/long/url",
    "expiration_days": 30
}
```

#### 2. Redirect
```http
GET /{shortID}
```

#### 3. Analytics
```http
GET /analytics?short_id={shortID}
```

#### 4. Health Check
```http
GET /health
```

For detailed API documentation, see [API.md](API.md)

## 🔧 Configuration

### Environment Variables

- `BASE_URL`: Base URL for shortened links (default: http://localhost:8080)
- `REDIS_HOST`: Redis host (default: redis)
- `REDIS_PORT`: Redis port (default: 6379)
- `DB_PATH`: SQLite database path (default: ./data/urlshortener.db)

### Docker Configuration

The service is configured through `docker-compose.yml`:

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - DB_TYPE=sqlite
      - DB_PATH=/app/data/urlshortener.db
    volumes:
      - ./data:/app/data
    depends_on:
      - redis

  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
```

## 📊 Analytics Features

- Click tracking
- Geographic data
- Device information
- Expiration tracking

## 🔒 Security Features

- URL validation
- Rate limiting
- Input sanitization
- CORS configuration

## 🧪 Testing

Run tests with:
```bash
go test ./...
```

## 📈 Performance Considerations

- Redis caching for frequently accessed URLs
- Rate limiting to prevent abuse
- Efficient database queries
- Connection pooling

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📚 Documentation Recommendations

### 1. API Documentation
- Create detailed API documentation using OpenAPI/Swagger
- Include request/response examples
- Document rate limits and error codes

### 2. Architecture Documentation
- Create architecture diagrams
- Document component interactions
- Explain data flow

### 3. Deployment Guide
- Document production deployment steps
- Include scaling considerations
- Add monitoring setup instructions

### 4. Development Guide
- Document development environment setup
- Include testing procedures
- Add contribution guidelines

### 5. User Guide
- Create user documentation
- Include common use cases
- Add troubleshooting guide

## 🔮 Future Improvements

1. **Custom Short URLs**: Allow users to specify custom short IDs
2. **User Authentication**: Add user accounts and authentication
3. **API Keys**: Implement API key authentication
4. **Bulk Operations**: Support bulk URL shortening
5. **Advanced Analytics**: Add more detailed analytics features
6. **Custom Domains**: Support custom domains for shortened URLs
7. **QR Code Generation**: Add QR code generation for shortened URLs
8. **Link Preview**: Add link preview functionality
9. **Rate Limiting Dashboard**: Add a dashboard for rate limit monitoring
10. **Export Features**: Add analytics export functionality

## 📞 Support

For support, please open an issue in the GitHub repository. 