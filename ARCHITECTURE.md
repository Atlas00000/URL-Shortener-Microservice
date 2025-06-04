# URL Shortener Architecture Documentation

## System Overview
The URL shortener is a microservice-based application designed for high performance, scalability, and reliability. It uses a modern stack of technologies and follows best practices for distributed systems.

## Architecture Components

### 1. Web Interface
- **Technology**: HTML, CSS, JavaScript
- **Location**: Served directly from the Go server
- **Features**:
  - Modern, responsive design
  - Real-time URL shortening
  - Analytics visualization
  - Error handling and user feedback

### 2. API Server
- **Technology**: Go with Gin framework
- **Features**:
  - RESTful API endpoints
  - Rate limiting
  - Request validation
  - Error handling
  - CORS support
  - Health checks

### 3. Database Layer
- **Primary Database**: SQLite
- **Schema**:
  ```sql
  CREATE TABLE urls (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      short_id TEXT UNIQUE NOT NULL,
      long_url TEXT NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      expires_at DATETIME NOT NULL
  );

  CREATE TABLE clicks (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      url_id INTEGER NOT NULL,
      country TEXT NOT NULL,
      device_type TEXT NOT NULL,
      clicked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (url_id) REFERENCES urls(id)
  );
  ```

### 4. Cache Layer
- **Technology**: Redis
- **Purpose**: 
  - URL mapping cache
  - Rate limiting
  - Analytics data caching
- **Configuration**:
  - Default TTL: 24 hours
  - Max memory: 100MB
  - Eviction policy: LRU

### 5. Analytics System
- **Components**:
  - Click tracking
  - Geographic data collection
  - Device detection
  - Expiration tracking
- **Storage**: SQLite database
- **Processing**: Real-time aggregation

## Data Flow

### URL Shortening Flow
1. Client sends URL to `/shorten` endpoint
2. Server validates URL format
3. Generates unique short ID
4. Stores mapping in Redis cache
5. Persists data in SQLite
6. Returns shortened URL to client

### URL Redirection Flow
1. Client requests shortened URL
2. Server checks Redis cache for mapping
3. If not in cache, queries SQLite
4. Updates cache with mapping
5. Records click analytics
6. Redirects to original URL

### Analytics Flow
1. Click event received
2. Geographic data extracted
3. Device type detected
4. Data stored in SQLite
5. Cache updated if necessary
6. Real-time aggregation performed

## Security Measures

### 1. Input Validation
- URL format validation
- SQL injection prevention
- XSS protection
- Input sanitization

### 2. Rate Limiting
- IP-based rate limiting
- Redis-backed counter
- Configurable limits
- Rate limit headers

### 3. CORS Configuration
- Allowed origins configuration
- Method restrictions
- Header restrictions
- Credential handling

### 4. Error Handling
- Structured error responses
- Logging and monitoring
- Graceful degradation
- Security headers

## Performance Considerations

### 1. Caching Strategy
- Two-level caching (Redis + SQLite)
- Cache invalidation policies
- TTL management
- Memory optimization

### 2. Database Optimization
- Indexed queries
- Connection pooling
- Efficient schema design
- Query optimization

### 3. Load Handling
- Rate limiting
- Request queuing
- Resource management
- Horizontal scaling capability

## Deployment Architecture

### Docker Configuration
```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - REDIS_URL=redis:6379
      - BASE_URL=http://localhost:8080
    depends_on:
      - redis

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

### Environment Variables
- `REDIS_URL`: Redis connection string
- `BASE_URL`: Base URL for shortened links
- `RATE_LIMIT`: Requests per minute
- `CACHE_TTL`: Cache time-to-live

## Monitoring and Logging

### 1. Health Checks
- API endpoint: `/health`
- Database connectivity
- Redis connectivity
- System metrics

### 2. Logging
- Request logging
- Error logging
- Performance metrics
- Security events

### 3. Metrics
- Response times
- Error rates
- Cache hit rates
- Resource usage

## Future Improvements

### 1. Scalability
- Horizontal scaling
- Load balancing
- Database sharding
- Cache distribution

### 2. Features
- Custom short URLs
- User authentication
- API key management
- Advanced analytics

### 3. Infrastructure
- Kubernetes deployment
- Service mesh
- Automated scaling
- Multi-region support

## Development Guidelines

### 1. Code Organization
```
.
├── main.go
├── config/
├── handlers/
├── models/
├── services/
├── utils/
└── static/
```

### 2. Testing Strategy
- Unit tests
- Integration tests
- Load tests
- Security tests

### 3. Documentation
- API documentation
- Architecture documentation
- Deployment guides
- Development guides 