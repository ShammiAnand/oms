# Order Management Service

## Order service

### gRPC Server

- exposes a gRPC server at port 9000

## Kitchen service

### API Server

- Kitchen service exposes `POST /orders` to create order
- Kitchen service exposes `GET /orders/{customerID}` to list all orders for a customer ID
  - ideally the above two endpoints should have authentication

### gRPC Server

- communicates with the gRPC connection exposed by Order service at port 9000

TODOs:

1. add database (mongo?)
2. add jwt auth -- login / signup routes
