# Email Validator API

A stateless REST API service that validates email addresses in real-time using external validation services.

## Features

- Real-time email validation
- Basic syntax validation
- Disposable email detection
- Free email provider detection
- RESTful API endpoints
- CORS enabled
- No data storage (completely stateless)

## Prerequisites

- Go 1.21 or higher
- Mailboxlayer API key (sign up at https://mailboxlayer.com)

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd tempmailblock
```

2. Copy the environment file and add your API keys:
```bash
cp .env.example .env
```

3. Edit `.env` and add your API keys:
```
MAILBOXLAYER_API_KEY=your_key_here
PORT=8080
```

4. Install dependencies:
```bash
go mod download
```

5. Run the server:
```bash
go run main.go
```

## API Endpoints

### Validate Email
```
GET /validate?email=test@example.com
```

#### Response
```json
{
    "email": "test@example.com",
    "is_valid": true,
    "provider": "example.com",
    "is_free": false,
    "is_disposable": false
}
```

### Health Check
```
GET /health
```

#### Response
```json
{
    "status": "healthy"
}
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- 200: Successful validation
- 400: Missing or invalid parameters
- 500: Internal server error or external API failure

## Security Considerations

- API keys are stored in environment variables
- CORS is enabled but can be configured as needed
- Rate limiting should be implemented in production
- Use HTTPS in production

## License

MIT License 