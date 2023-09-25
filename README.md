# Booking

This is a REST API for booking hotel rooms along with registration, authentication and canceling bookings. It also implement adminstration, where admin can delete users along with their bookings.

## Note

- There is only one administrator, which is created at initialization of the database.
  Login : **007@gmail.com**
  Password : **Cheburek**

## Technologies

Project was built using Fiber as a web framework and MongoDB as a database.

## Installation

To run this project, you need to have Docker installed on your machine. Then follow these steps:

1. Clone this repository.
2. Navigate to the project directory.
3. Run the following command to build the Docker image:

```bash
make build
```

4. Run the following command to start the Docker container:

```bash
make start
```

5. Open your web browser and go to http://localhost:{port number}.

If you want to run without Docker:

```bash
go run .
```
