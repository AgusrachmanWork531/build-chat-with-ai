# Go Secure Backend Service with Clean Architecture

This project is a secure, high-performance, and maintainable RESTful API built in Golang using the Echo framework. It demonstrates the principles of Clean Architecture, SOLID, and Clean Code, with a focus on robust authentication and separation of concerns.

## Features

-   **User Management:** Core CRUD functionality for a User entity.
-   **Secure Endpoints:**
    -   **JWT Authentication:** Protects user-specific endpoints.
    -   **Basic Authentication:** Protects administrative or internal endpoints.
-   **Clean Architecture:** A clear separation between business logic and framework-specific code.
-   **Dependency Injection:** Interfaces are used to decouple layers.
-   **Configuration Management:** Centralized configuration for application settings.
-   **Makefile:** Simplified commands for running, building, and managing the project.

## Architecture

This project follows the principles of **Clean Architecture**. The dependencies flow inwards, from the outer layers (frameworks, UI) to the inner layers (business logic, entities).

```
+-----------------------------------------------------------------+
|                            Frameworks                           |
| +-------------------------------------------------------------+ |
| |                          Handlers                           | |
| | +---------------------------------------------------------+ | |
| | |                         Usecases                        | | |
| | | +-----------------------------------------------------+ | | |
| | | |                      Entities                       | | | |
| | | +-----------------------------------------------------+ | | |
| | +---------------------------------------------------------+ | |
| +-------------------------------------------------------------+ |
+-----------------------------------------------------------------+
```

-   **Domain/Entities (`internal/user/domain.go`):** The core of the application. It contains the `User` struct and the interfaces for the `UserRepository` and `UserUsecase`. This layer has no dependencies on any other layer.
-   **Repository (`internal/user/repository.go`):** Implements the `UserRepository` interface. It is responsible for data persistence. In this project, an in-memory repository is used for demonstration purposes.
-   **Usecase (`internal/user/usecase.go`):** Implements the `UserUsecase` interface. It contains the core business logic of the application, such as user creation, retrieval, and authentication. It depends on the `UserRepository` interface.
-   **Handler (`internal/user/handler.go`):** This layer is responsible for handling HTTP requests. It uses the `UserUsecase` to perform business operations. It depends on the `UserUsecase` interface.
-   **Main (`cmd/server/`):** This is the entry point of the application. It is responsible for initializing all the components and starting the server.

## API Endpoints

| Method | Path              | Authentication | Description                  |
|--------|-------------------|----------------|------------------------------|
| `POST` | `/v1/login`       | Public         | Authenticate and get a JWT.  |
| `POST` | `/v1/users`       | Basic Auth     | Create a new user.           |
| `GET`  | `/v1/users/:id`   | JWT            | Get a user by their ID.      |

## Getting Started

### Prerequisites

-   Go 1.18 or higher
-   `make`

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/gemini-cli/portfolio-chat-ai-go.git
    cd portfolio-chat-ai-go
    ```

2.  Tidy the dependencies:
    ```bash
    make tidy
    ```

### Running the Application

To run the application, use the following command:

```bash
make run
```

The server will start on `http://localhost:8080`.

## How to Use

### 1. Create a User (Basic Auth)

Use the default credentials `admin:password` for Basic Auth to authorize the creation, and provide the new user's credentials in the URL.

```bash
# The -u flag provides the Basic Auth credentials (admin:password)
# The user's email and password for the new user are passed in the request body
curl -X POST http://localhost:8080/v1/users -u "admin:password" \
-H "Content-Type: application/json" \
-d 
    "email": "newuser@example.com",
    "password": "password123"
}'
```

### 2. Login (Public)

Use the email and password of the user you just created.

```bash
curl -X POST http://localhost:8080/v1/login \
-H "Content-Type: application/json" \
-d 
    "email": "newuser@example.com",
    "password": "password123"
}'
```

This will return a JWT token.

### 3. Get User by ID (JWT)

Copy the JWT token from the login response and the user ID from the create user response.

```bash
export TOKEN="your.jwt.token"
export USER_ID="the-user-id"

curl -X GET http://localhost:8080/v1/users/$USER_ID \
-H "Authorization: Bearer $TOKEN"
```

## Project Structure

```
.
├── .env
├── Makefile
├── README.md
├── cmd
│   └── server
│       ├── main.go
│       └── router.go
├── go.mod
├── go.sum
├── internal
│   ├── chat
│   │   ├── domain.go
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── usecase.go
│   └── user
│       ├── domain.go
│       ├── handler.go
│       ├── repository.go
│       ├── repository_mongo.go
│       └── usecase.go
└── pkg
    ├── bootstrap
    │   └── bootstrap.go
    ├── config
    │   └── config.go
    ├── database
    │   └── mongo.go
    └── middleware
        └── auth.go
```

-   **`.env`**: Menyimpan environment variables untuk development lokal.
-   **`Makefile`**: Berisi perintah untuk menjalankan, membangun, dan merapikan proyek.
-   **`cmd/server`**: Titik masuk utama aplikasi dan penyiapan router.
-   **`internal`**: Berisi semua logika bisnis inti, dipisahkan berdasarkan domain (`user`, `chat`).
    -   `domain.go`: Mendefinisikan struct dan interface inti untuk domain tersebut.
    -   `handler.go`: Lapisan handler untuk HTTP/WebSocket.
    -   `usecase.go`: Lapisan yang berisi logika bisnis inti.
    -   `repository.go`: Implementasi repository in-memory (segera dihapus).
    -   `repository_mongo.go`: Implementasi repository yang menggunakan MongoDB.
-   **`pkg`**: Berisi paket-paket yang dapat dibagikan dan digunakan di seluruh aplikasi.
    -   `bootstrap`: Memusatkan logika startup aplikasi (koneksi DB, dll.).
    -   `config`: Menangani pemuatan konfigurasi dari environment.
    -   `database`: Menyediakan helper untuk koneksi database.
    -   `middleware`: Berisi middleware Echo kustom (JWT, Basic Auth).
