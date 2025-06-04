# URL Shortener API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Currently, the API does not require authentication. Future versions will implement API key authentication.

## Rate Limiting
- 60 requests per minute per IP address
- Rate limit headers are included in responses:
  - `X-RateLimit-Limit`: Maximum requests per window
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Time until rate limit resets

## Endpoints

### 1. Shorten URL
Creates a shortened version of a long URL.

**Endpoint:** `POST /shorten`

**Request Body:**
```json
{
    "url": "https://example.com/very/long/url",
    "expiration_days": 30  // Optional, defaults to 30 days
}
```

**Response:**
```json
{
    "short_url": "http://localhost:8080/YtHDX-8",
    "long_url": "https://example.com/very/long/url",
    "expires_at": "2024-07-03T13:28:20.59Z"
}
```

**Status Codes:**
- `201 Created`: URL successfully shortened
- `400 Bad Request`: Invalid URL format or request body
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

### 2. Redirect to Original URL
Redirects to the original URL using the shortened ID.

**Endpoint:** `GET /{shortID}`

**Parameters:**
- `shortID` (path parameter): The shortened URL identifier

**Response:**
- Redirects to the original URL with status code 301 (Moved Permanently)

**Status Codes:**
- `301 Moved Permanently`: Successful redirect
- `404 Not Found`: URL not found or expired
- `400 Bad Request`: Invalid short ID

### 3. Get Analytics
Retrieves analytics data for a specific URL.

**Endpoint:** `GET /analytics`

**Query Parameters:**
- `short_id` (required): The short ID of the URL to get analytics for

**Response:**
```json
{
    "total_clicks": 3,
    "clicks_by_country": {
        "US": 1,
        "AU": 1,
        "RU": 1
    },
    "clicks_by_device": {
        "mobile": 1,
        "tablet": 1,
        "desktop": 1
    },
    "last_click": "2024-06-03T13:28:58.7Z"
}
```

**Status Codes:**
- `200 OK`: Analytics retrieved successfully
- `400 Bad Request`: Missing or invalid short_id parameter
- `404 Not Found`: URL not found
- `500 Internal Server Error`: Server error

### 4. Record Click
Records a click event for a URL.

**Endpoint:** `POST /analytics/click`

**Query Parameters:**
- `short_id` (required): The short ID of the URL to record the click for

**Response:**
```json
{
    "status": "success"
}
```

**Status Codes:**
- `200 OK`: Click recorded successfully
- `400 Bad Request`: Missing or invalid short_id parameter
- `404 Not Found`: URL not found
- `500 Internal Server Error`: Server error

### 5. Health Check
Checks if the service is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
    "status": "healthy",
    "time": "2024-06-03T13:28:20.59Z"
}
```

**Status Codes:**
- `200 OK`: Service is healthy

## Error Responses
All error responses follow this format:
```json
{
    "error": "Error message description"
}
```

## Examples

### Shortening a URL
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url", "expiration_days": 30}'
```

### Getting Analytics
```bash
curl http://localhost:8080/analytics?short_id=YtHDX-8
```

### Recording a Click
```bash
curl -X POST http://localhost:8080/analytics/click?short_id=YtHDX-8
```

## Notes
- All timestamps are in UTC
- URLs must include scheme (http/https) and host
- Shortened URLs expire after the specified number of days
- Analytics data is stored in the database and can be retrieved at any time 