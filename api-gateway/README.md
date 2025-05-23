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