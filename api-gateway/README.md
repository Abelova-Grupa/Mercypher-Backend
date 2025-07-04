# Mercypher API Gateway

This is the API Gateway service for the **Mercypher** chat application. It acts as the central point of communication between clients and the internal microservices, handling HTTP, WebSocket, and gRPC traffic.

---

## 🔧 Features

- User registration and login
- Authentication middleware for WebSocket connections
- WebSocket support for chat messaging
- gRPC server for receiving messages and status updates from internal services

---

## 🚀 HTTP Endpoints

| Method | Path         | Description              |
|--------|--------------|--------------------------|
| POST   | `/login`     | Login user with email and password |
| POST   | `/register`  | Register a new user      |
| GET    | `/logout`    | Logout authenticated user |
| GET    | `/ws`        | WebSocket endpoint for chat and status updates (auth required) |

### 📝 `/register` and `/login` format
For registration, users are required to provide a username, email address, and password.

```json
{
    "username":"exampleUsername",
    "email":"example@email.xyz",
    "password":"examplePassword123" 
}
```

For logging in, users must submit their username and password, along with an optional authentication token.

```json
{
    "username":"exampleUsername",
    "password":"examplePassword123",
    "token":"exampleToken"  // Token is optional
}
```

Those request formats are defined in `./internal/servers/http_server.go`:
```go
type LoginRequest struct {
	Username 	string `json:"username" binding:"required"`
    Password 	string `json:"password" binding:"required"`
	Token		string `json:"token"`
}

type RegisterRequest struct {
	Username 	string `json:"username" binding:"required"`
	Email	 	string `json:"email" binding:"required"`
    Password 	string `json:"password" binding:"required"`
}
```

### 🔒 Authentication

- The `/ws` route is protected by `AuthMiddleware()`.
- Clients must send a valid token to establish a WebSocket connection.

---

## 🌐 WebSocket Communication

Once connected to `/ws`, the client sends and receives messages using an `Envelope` format.

### 📦 Envelope Format

```json
{
  "type": "message", // or "search", "status"
  "payload": { ... } // content varies by type
}
```

> For more details on payload formats, refer to `./internal/domain` directory.

---

## 🔌 GRPC Communication

Here are example gRPC messages that gateway handles as a server.

```
{
    "chat_message": {
        "body": "Hello World!",
        "message_id": "MSG1",
        "recipient_id": "USR1",
        "sender_id": "USR554",
        "timestamp": "49831638"
    },
    "message_status": {
        "message_id": "MSG1",
        "recipient_id": "USR1",
        "status": "SEEN",
        "timestamp": "49833413"
    }
}
```
## 🧑‍🎓 Developer note
> **Note:** The majority of the API Gateway implementation was developed by @jelisavac-l. Should any issues or unexpected behavior arise within this component, please do not hesitate to direct any questions, concerns, or constructive criticism my way. I take full responsibility for its current state and welcome feedback for future improvements.
