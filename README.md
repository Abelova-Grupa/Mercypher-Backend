# Mercypher Backend

End-to-end encrypted chat app backend built on a scalable microservice architecture.

Mercypher ensures secure and private messaging by implementing modern encryption standards, with each message encrypted client-side before transmission. This backend is composed of decoupled microservices, enabling horizontal scalability, service isolation, and independent deployment.

## 🔐 Features

- End-to-end encryption (E2EE) using modern cryptographic libraries  
- Stateless message relay via secure transport  
- Scalable microservices for authentication, messaging, storage, and presence  

## 🚀 Technologies

- Go 
- PostgreSQL / Redis  
- Docker + Terraform + Github Actions + Azure  

## Getting started

**Prerequisites**
- Go (1.23+ recommended)
- Docker & Docker Compose
- Node.js & npm (for the frontend)
- Make (optional, for automation scripts)

### 1. The Full Docker Deployment (Recommended)

The simplest way to start the entire stack, including all microservices and infrastructure, is using Docker Compose. This ensures all environment variables and network configurations are handled automatically.

```
docker compose up --build -d
```
### 2. Manual Service Execution

If you prefer to run the services individually for debugging or development, you must first start the required infrastructure components.

```
docker compose up -d kafka
sudo make redis-up
```
Navigate to the root directory and run each service using the Go toolchain:

```
go run ./api-gateway/cmd/api-gateway
go run ./group-service/cmd/group-service
go run ./session-service/cmd/session-service
...
```

### 3. Automated Local Dev

For developers working in a terminal-centric environment, a tmux script is provided in the repository. This script automates the process of splitting windows and launching all Go services simultaneously with a single command.

### Frontend Setup

The frontend is built with React and needs to be started separately from the backend services.

- Navigate to the frontend directory.
- Install dependencies: npm install.
- Start the development server:

```
npm run dev
```

> ⚠️ **Note:** This repo contains backend logic only. The client-side encryption and UI are handled in the [Mercypher Client](#) repository.

