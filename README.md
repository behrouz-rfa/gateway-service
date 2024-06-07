# Gateway Service

## Overview
The Gateway Service is a backend service that interacts with various external APIs, such as OpenAI and geolocation services, to provide functionalities like text generation and IP geolocation. It also manages user plans and caching to ensure efficient use of resources.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Error Handling](#error-handling)
- [Contributing](#contributing)
- [License](#license)

## Features
- **Caching:** Use Redis for caching responses to improve performance, with optimization for similar keys using cosine similarity.
- **Rate Limiter:** 100 request per hour for each ip but you can change it in .env
- **AI Text Generation:** Utilize OpenAI's GPT model for generating text.
- **IP Geolocation:** Retrieve geolocation information based on IP addresses.
- **User Plan Management:** Manage user plans and credits.
- **Authentication:** Secure API endpoints with JWT-based authentication.

## Installation
### Prerequisites
- Go 1.22+
- Redis
- Docker and Docker Compose

### Steps
1. Clone the repository:
    ```bash
    git clone https://github.com/behrouz-rfa/gateway-service.git
    cd gateway-service
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Setup environment variables. Create a `.env` file based on `.env.example` and configure your settings.

4. Run the service using Docker Compose:
    ```bash
    docker-compose up --build
    ```

## Configuration
Configure the service using environment variables. Refer to `.env.example` for all configurable options:
- `OPENAI_API_KEY`: Your OpenAI API key.
- `REDIS_URL`: URL for your Redis instance.
- `JWT_SECRET`: Secret key for JWT authentication.

## Usage
### Running with Docker Compose
1. Build and run the Docker containers:
    ```bash
    docker-compose up --build
    ```

### API Endpoints
- http://localhost:8080/swagger/index.html#/


## Caching with Redis
The service uses Redis for caching responses to improve performance. A feature to optimize caching for similar keys is implemented using cosine similarity.

### Cache Optimization
- **Normalization:** Text is normalized by converting to lowercase and removing punctuation.
- **Hashing:** A SHA256 hash of the normalized text is generated for caching.
- **Similarity Calculation:** Cosine similarity is calculated for text comparison.
- **Caching Mechanism:**
    - Retrieves the cached response if it exists based on proximity.
    - Stores the response in the cache with its vector.

```go
// Example of the Redis caching implementation

// NormalizeText normalizes the text by converting it to lowercase and removing punctuation
func NormalizeText(text string) string {
    // Implementation
}

// GenerateHash generates a SHA256 hash of the normalized text
func GenerateHash(text string) string {
    // Implementation
}

func CalculateCosineSimilarity(s1, s2 string) float64 {
    // Implementation
}

// GetCachedResponse retrieves the cached response if it exists based on proximity
func (r *RedisClient) GetCachedResponse(query string) ([]byte, bool) {
    // Implementation
}

// SetCachedResponse stores the response in the cache with its vector
func (r *RedisClient) SetCachedResponse(query string, response []byte) {
    // Implementation
}
