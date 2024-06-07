# Order Management Service

## Order service

### gRPC Server

- exposes a gRPC server at port 9000

## Kitchen service

### API Server

#### Un-Authenticated Routes

- Kitchen service exposes `POST /login` to create jwt token for already registered users
  - this token will be passed in `Authorization` header for all authenticated routes
- Kitchen service exposes `POST /signup` to create user

#### Authenticated Routes

- Kitchen service exposes `POST /orders` to create order
- Kitchen service exposes `GET /orders` to list all orders for a customer ID

### gRPC Server

- communicates with the gRPC connection exposed by Order service at port 9000
