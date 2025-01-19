# go-rateLimit: Token Bucket Rate Limiting for Multiple Clients

This repository implements a **Token Bucket** rate limiting algorithm in Go to handle **rate-limiting** for multiple clients. It allows for a specific number of requests to be made within a given time period (e.g., 3 requests every 30 seconds), with each client having independent rate limits based on unique identifiers (such as mobile numbers or IP addresses).

## Features
- **Token Bucket Rate Limiting**: Uses the token bucket algorithm for rate limiting. Clients are allowed a certain number of requests, and once the limit is reached, they need to wait until tokens are refilled.
- **Multiple Clients Support**: Rate limits are applied individually per client, identified by their unique identifier (e.g., mobile number).
- **JSON Request Handling**: Incoming requests are parsed as JSON, and the data (such as mobile number) is used to identify clients for rate limiting.
- **Customizable Rate Limits**: Configure the maximum number of requests and the refill rate for the token bucket.
- **Gin Framework**: Built on top of the **Gin** web framework for fast and efficient HTTP server handling.
- **Middleware for Rate Limiting**: The middleware ensures that rate-limiting logic is applied to HTTP requests.

## Use Cases
- **API Rate Limiting**: Implement rate limiting on your APIs, controlling how many requests each client can make in a specific time frame.
- **Client-Specific Rate Limits**: Apply rate limiting on a per-client basis, using unique identifiers such as phone numbers or IP addresses.
- **Scalable and Flexible**: Handles multiple clients with independent rate limits, making it scalable for high-traffic systems.

## How to Use

### Prerequisites

- **Go (Golang)** version 1.16 or above.
- **Gin Framework** (automatically handled via Go modules).

### 1. Clone the Repository
Clone this repository to your local machine using the following command:

```bash
git clone https://github.com/ArjunDev17/go-rateLimit.git
cd go-rateLimit
```

### 2. Install Dependencies
Install the required dependencies using Go modules:

```bash
go mod tidy
```

### 3. Run the Server
Run the server to start accepting API requests:

```bash
go run main.go
```

The server will run on `localhost:8081`.

### 4. Make API Requests
You can make a POST request to `/api/v1/onboard` with a JSON payload. For example:

```json
{
  "name": "John",
  "mobileNumber": "1234567890"
}
```

### Example

Assume the rate-limiting parameters are configured to allow **3 requests per 30 seconds**. If a client (identified by their mobile number) exceeds this limit, they will receive a **429 Too Many Requests** response.

#### Example Request:
```bash
curl -X POST http://localhost:8081/api/v1/onboard -d '{"name": "John", "mobileNumber": "1234567890"}' -H "Content-Type: application/json"
```

### Rate Limiting Behavior
- If a client (identified by mobile number) sends more than the allowed number of requests in the configured time period (e.g., 3 requests every 30 seconds), the server will respond with a **429 Too Many Requests** error.
- If the client is within the limit, the server will proceed and respond with a success message.

#### Response on Success (within limit):
```json
{
  "message": "User onboarded successfully!",
  "mobile": "1234567890"
}
```

#### Response on Rate Limit Exceeded:
```json
{
  "error": "Rate limit exceeded. Try again later."
}
```

## Configuration

You can customize the rate-limiting parameters by adjusting the configuration values in the `main.go` file:
- **Maximum Tokens**: The maximum number of tokens allowed in the bucket (e.g., 3 tokens).
- **Refill Rate**: The time interval (e.g., 30 seconds) in which tokens are refilled.
- **Refill Count**: The number of tokens to add each time the bucket is refilled.

For example:
```go
r.Use(rateLimiter.LimitMiddleware(3, 30*time.Second, 1)) // 3 tokens max, 1 token every 30 seconds
```

## How Token Bucket Works

The **Token Bucket** algorithm works as follows:
- The bucket holds a fixed number of tokens, representing the number of requests a client can make.
- Each time a request is made, one token is consumed.
- Tokens are refilled at a defined rate over time. If the client exceeds the number of available tokens, the request is rejected with a rate-limiting error (HTTP 429).
- Once the client waits long enough, the tokens refill, allowing them to make more requests.

## Example Folder Structure
```
go-rateLimit/
│
├── main.go            # Main application file
├── middleware/        # Custom middleware for rate limiting
│   └── rate_limiter.go
├── go.mod             # Go module definition
├── go.sum             # Go module checksum
└── README.md          # This file
```

## Technologies Used
- **Go (Golang)**: The programming language used to implement the rate-limiting logic.
- **Gin**: A fast web framework used for routing and handling HTTP requests.
- **Token Bucket Algorithm**: A rate-limiting algorithm to manage the flow of requests from clients.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request for improvements, bug fixes, or new features.

### License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

### Notes:
- Replace `yourusername` in the `git clone` command with your actual GitHub username.
- Adjust the configuration in the `main.go` file as needed based on your requirements.
- Ensure your server is running before testing API requests with tools like `curl` or Postman.

---

This `README.md` provides a clear explanation of the repository’s functionality, setup instructions, and example use cases, making it easy for anyone to understand and use your rate-limiting implementation.
