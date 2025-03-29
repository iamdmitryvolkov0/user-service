# User Service

## Overview

User Service is a microservice responsible for managing user data, authentication, and related functionalities. It is
built using Go and follows a clean architecture approach. The service exposes RESTful APIs and provides Swagger
documentation for easy integration and testing.

## Features

- User registration and authentication
- Password hashing and secure storage
- Database migrations and seeding
- RESTful API with structured responses
- Swagger UI for API documentation

## Technologies Used

- **Go** (Golang)
- **Chi Router** - lightweight and fast router for handling HTTP requests
- **PostgreSQL** - relational database
- **Swaggo** - Swagger documentation generator for Go
- **Docker** - containerization

## Setup and Installation

### Prerequisites

Make sure you have the following installed on your system:

- [Go](https://go.dev/doc/install) (version 1.24 or later)
- [Docker](https://www.docker.com/get-started) (optional, for containerization)
- [PostgreSQL](https://www.postgresql.org/download/)

### Clone the Repository

```sh
git clone https://github.com/iamdmitryvolkov0/user-service.git
cd user-service
```

### Environment Configuration

Copy template using command

```
cp .env.example .env
```

or create a `.env` file in the root directory and add the necessary environment variables:

```env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=user_db

JWT_SECRET=hash

COMPOSE_BAKE=1
```

### Running the Application

#### Running with Docker

1. Build and run the Docker container:
    ```sh
    make up
    ```

## API Documentation

Swagger UI is available here:
[click me](http://localhost:8080/swagger/index.html#/)

### Generate Swagger Docs

If you modify API definitions, regenerate the documentation with:

```sh
make docs
```

## Database Migrations

Applying automatically every time container starts.
<br>Also, if users table is empty, seeder will create some.

## Contact

For any issues or suggestions, feel free to open an issue on GitHub or contact the maintainers.

