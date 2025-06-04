# **Product Requirements Document: Microservice URL Shortener**

## **1\. Overview**

### **1.1 Purpose**

The Microservice URL Shortener is a backend-focused portfolio project designed to demonstrate core backend development skills, including building a scalable, reliable, and maintainable microservice. It provides functionality to shorten URLs, track clicks, and provide basic geo/device analytics, using Go or Python, Redis, and PostgreSQL, while adhering to industry best practices for simplicity and performance under high traffic.

### **1.2 Scope**

* Create and manage short URLs.  
* Track clicks on short URLs.  
* Provide basic analytics (geo-location and device type).  
* Handle high traffic with efficient caching and database design.  
* Focus on simplicity, avoiding over-engineering, suitable for a portfolio project.

### **1.3 Target Audience**

* Portfolio reviewers (hiring managers, technical recruiters).  
* Developers learning microservices, caching, and analytics.

## **2\. Functional Requirements**

### **2.1 URL Shortening**

* **Feature**: Users can submit a long URL and receive a shortened URL.  
* **Details**:  
  * Input: Valid URL (e.g., [https://example.com/very/long/url](https://example.com/very/long/url)).  
  * Output: Short URL (e.g., [https://short.url/abcd123](https://short.url/abcd123)).  
  * Short URLs are unique, random, 7-character alphanumeric strings (e.g., abcd123).  
  * Short URLs expire after a configurable period (default: 30 days).  
  * API endpoint: `POST /shorten` with JSON payload `{ "url": "https://example.com" }`.  
  * Response: JSON with short URL `{ "short_url": "https://short.url/abcd123" }`.

### **2.2 URL Redirection**

* **Feature**: Users accessing a short URL are redirected to the original URL.  
* **Details**:  
  * Input: Short URL (e.g., [https://short.url/abcd123](https://short.url/abcd123)).  
  * Output: HTTP 301 redirect to the original URL.  
  * API endpoint: `GET /:short_id` (e.g., `/abcd123`).  
  * Invalid or expired URLs return a 404 error.

### **2.3 Click Tracking**

* **Feature**: Track each click on a short URL.  
* **Details**:  
  * Record timestamp, IP address, and user-agent for each click.  
  * Store click data in PostgreSQL for persistence.  
  * API endpoint: `GET /analytics/:short_id` to retrieve click count and details.

### **2.4 Geo/Device Analytics**

* **Feature**: Provide basic analytics on clicks (geo-location and device type).  
* **Details**:  
  * Geo-location: Derive country/region from IP address using a lightweight geo-IP database (e.g., MaxMind GeoLite2 free tier).  
  * Device type: Parse user-agent to identify device (e.g., mobile, desktop, tablet).  
  * Analytics available via `GET /analytics/:short_id` with aggregated data (e.g., clicks by country, device type).  
  * Focus on simplicity: basic aggregation (e.g., total clicks, top 3 countries, device breakdown).

### **2.5 High-Traffic Handling**

* **Feature**: Ensure the service handles high traffic efficiently.  
* **Details**:  
  * Use Redis for caching short URL mappings to reduce database load.  
  * Cache TTL matches URL expiration (default: 30 days).  
  * PostgreSQL for persistent storage of URLs and click data.  
  * Implement rate limiting on API endpoints to prevent abuse (e.g., 100 requests/min per IP).

## **3\. Non-Functional Requirements**

### **3.1 Performance**

* Redirect latency: \< 100ms for cached URLs (Redis).  
* Short URL creation: \< 500ms.  
* Analytics response: \< 1s for up to 10,000 clicks.  
* Target throughput: Handle 1,000 requests/second with caching.

### **3.2 Scalability**

* Microservice architecture allows horizontal scaling (e.g., multiple instances behind a load balancer).  
* Redis and PostgreSQL can be scaled independently (e.g., Redis cluster, PostgreSQL read replicas if needed).  
* Keep design simple, avoiding premature optimization.

### **3.3 Reliability**

* Ensure 99.9% uptime for URL redirection.  
* Handle database failures gracefully by relying on Redis cache.  
* Log errors for debugging (e.g., invalid URLs, database errors).

### **3.4 Security**

* Validate input URLs (e.g., ensure valid protocol, no malicious scripts).  
* Sanitize user-agent data to prevent injection attacks.  
* Implement rate limiting to mitigate DDoS risks.  
* Use HTTPS for all API endpoints.

### **3.5 Maintainability**

* Clean, modular code following Go/Python best practices (e.g., Go modules, Python PEP 8).  
* Comprehensive unit tests for core functionality (e.g., URL shortening, redirection, analytics).  
* Basic logging for monitoring and debugging.  
* Use environment variables for configuration (e.g., database credentials, Redis host).

## **4\. Technical Stack**

* **Language**: Go or Python (Go for performance, Python for readability; choose based on portfolio goals).  
* **Database**:  
  * **Redis**: Cache short URL mappings and session data.  
  * **PostgreSQL**: Persistent storage for URLs and click data.  
* **Geo-IP**: MaxMind GeoLite2 (free tier) for geo-location.  
* **Libraries/Frameworks**:  
  * Go: `gorilla/mux` (routing), `go-redis` (Redis client), `lib/pq` (PostgreSQL).  
  * Python: `FastAPI` or `Flask` (API), `redis-py` (Redis), `psycopg2` (PostgreSQL).  
* **Deployment**: Docker for containerization, deployable on a single server (e.g., AWS EC2, DigitalOcean).  
* **Testing**: Go: `testing` package; Python: `pytest`.

## **5\. System Architecture**

### **5.1 Components**

* **API Service**: Handles HTTP requests (shorten, redirect, analytics).  
* **Redis**: Caches short URL mappings and session data.  
* **PostgreSQL**: Stores URLs and click data.  
* **Geo-IP Service**: Lightweight library to map IPs to locations.

### **5.2 Data Flow**

1. **Shorten URL**:  
   * Client sends `POST /shorten` with long URL.  
   * Service validates URL, generates unique short ID, stores in PostgreSQL, caches in Redis.  
   * Returns short URL.  
2. **Redirect**:  
   * Client accesses `GET /:short_id`.  
   * Service checks Redis cache; if not found, queries PostgreSQL.  
   * Logs click (IP, user-agent, timestamp) to PostgreSQL asynchronously.  
   * Returns 301 redirect.  
3. **Analytics**:  
   * Client requests `GET /analytics/:short_id`.  
   * Service queries PostgreSQL for click data, aggregates by geo/device.  
   * Returns JSON with analytics.

### **5.3 Database Schema**

* **Table: urls**  
  * `id`: Primary key (auto-increment).  
  * `short_id`: Unique 7-char alphanumeric (indexed).  
  * `long_url`: Original URL (varchar).  
  * `created_at`: Timestamp.  
  * `expires_at`: Timestamp (default: 30 days from creation).  
* **Table: clicks**  
  * `id`: Primary key (auto-increment).  
  * `short_id`: Foreign key to `urls` (indexed).  
  * `timestamp`: Click timestamp.  
  * `ip_address`: Client IP (varchar).  
  * `user_agent`: Client user-agent (varchar).  
  * `country`: Derived from IP (varchar, nullable).  
  * `device_type`: Derived from user-agent (varchar, nullable).

## **6\. API Endpoints**

* **POST /shorten**  
  * Request: `{ "url": "https://example.com" }`  
  * Response: `{ "short_url": "https://short.url/abcd123" }`  
  * Errors: 400 (invalid URL), 429 (rate limit).  
* **GET /:short\_id**  
  * Response: 301 redirect to long URL.  
  * Errors: 404 (invalid/expired URL), 429 (rate limit).  
* **GET /analytics/:short\_id**  
  * Response: `{ "short_id": "abcd123", "total_clicks": 100, "by_country": { "US": 60, "UK": 30, "Other": 10 }, "by_device": { "mobile": 70, "desktop": 20, "tablet": 10 } }`  
  * Errors: 404 (invalid/expired URL), 429 (rate limit).

## **7\. Constraints and Assumptions**

* **Constraints**:  
  * Limited to basic analytics (country, device type) to avoid complexity.  
  * No user authentication to keep the scope manageable.  
  * Single-region deployment for simplicity.  
* **Assumptions**:  
  * Traffic is moderate (up to 1,000 requests/second).  
  * Geo-IP data is approximate (free tier MaxMind GeoLite2).  
  * No real-time analytics updates; batch processing is sufficient.

## **8\. Future Enhancements (Out of Scope)**

* User authentication for private URLs.  
* Custom short URLs.  
* Real-time analytics with streaming.  
* Multi-region deployment for lower latency.

## **9\. Deliverables**

* Source code in Go or Python (GitHub repository).  
* Dockerized application with `docker-compose` for Redis and PostgreSQL.  
* Unit tests covering core functionality.  
* README with setup, deployment, and API usage instructions.  
* Basic Postman collection for API testing.

## **10\. Success Criteria**

* Successfully shorten and redirect URLs with \< 100ms latency (cached).  
* Track clicks and provide accurate analytics.  
* Handle 1,000 requests/second without failure.  
* Clean, documented code adhering to best practices.  
* Deployable on a single server with minimal configuration.

