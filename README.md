# Task Management Service

A task management service built with Go and MongoDB, providing functionality for adding, updating, deleting, and retrieving tasks. This service offers a RESTful API for managing tasks and demonstrates the use of MongoDB for data storage.

## Project Structure

Here's the updated project structure, grouped by layers and presented in a table format:

### Project Structure

### **API**

```
api
├── controllers
│   ├── auth
│   │   ├── dto
│   │   │   ├── auth_request.go
│   │   │   └── auth_response.go
│   │   ├── controller.go
│   │   └── controller_test.go
│   ├── base
│   │   └── controller.go
│   ├── task
│   │   ├── dto
│   │   │   ├── add_request.go
│   │   │   └── get_response.go
│   │   ├── task_controller.go
│   │   └── task_controller_test.go
│   └── user
│       ├── controller.go
│       └── controller_test.go
├── errors
│   └── error.go
├── middleware
│   └── auth
│       ├── role.go
│       └── role_test.go
├── router
│   └── router.go
└── i_controller.go
```

### **App**

```
app
├── common
│   ├── cqrs
│   │   ├── command
│   │   │   ├── mocks
│   │   │   │   └── command_handler_mock.go
│   │   │   └── command_handler.go
│   │   └── query
│   │       ├── mocks
│   │       │   └── query_handler_mock.go
│   │       └── query_handler.go
│   ├── i_jwt
│   │   ├── mock
│   │   │   └── service_mock.go
│   │   └── service.go
│   └── i_repo
│       ├── mocks
│       │   ├── task_mock.go
│       │   └── user_mock.go
│       ├── task.go
│       └── user.go
├── task
│   ├── command
│   │   ├── add
│   │   │   ├── command.go
│   │   │   ├── handler.go
│   │   │   └── handler_test.go
│   │   ├── delete
│   │   │   ├── handler.go
│   │   │   └── handler_test.go
│   │   └── update
│   │       ├── command.go
│   │       ├── handler.go
│   │       └── handler_test.go
│   └── query
│       ├── get
│       │   ├── handler.go
│       │   └── handle_test.go
│       └── get_all
│           ├── handler.go
│           └── handler_test.go
└── user
    ├── admin_status
    │   └── command
    │       ├── command.go
    │       ├── handler.go
    │       └── handler_test.go
    └── auth
        ├── command
        │   ├── command.go
        │   ├── handler.go
        │   └── handler_test.go
        ├── common
        │   └── result.go
        └── query
            ├── handler.go
            ├── handler_test.go
            └── query.go
```

### **Domain**

```
domain
├── errors
│   ├── common.go
│   ├── task.go
│   └── user.go
├── i_hash
│   ├── mocks
│   │   └── Service.go
│   └── service.go
└── models
    ├── task
    │   ├── task.go
    │   └── task_test.go
    └── user
        ├── user.go
        └── user_test.go
```

### **Infrastructure**

```
infrastructure
├── db
│   └── db.go
├── hash
│   ├── service.go
│   └── service_test.go
├── jwt
│   ├── service.go
│   └── service_test.go
└── repo
    ├── task
    │   ├── repo.go
    │   └── repo_test.go
    └── user
        └── repo.go
```

### **Others**

```
bin
└── task-manager
config
└── env.go
docs
└── api_definition.md
scripts
└── filter_test_packages.sh
coverage.out
example.env
go.mod
go.sum
main.go
Makefile
README.md
```

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/beka-birhanu/task-manager-authentication.git
   ```

````

2. Change to the project directory:

   ```bash
   cd task-manager-authentication
   ```

3. Install the Go dependencies:

   ```bash
   go mod tidy
   ```

## Configuration

Before running the application, ensure you have a MongoDB instance running and update the configuration in the `.env` file with your specific details:

1. **Clone the provided example environment file**:

   ```bash
   cp example.env .env
   ```

2. **Update the `.env` file** with the following configurations:

   ```plaintext
   PUBLIC_HOST=http://localhost                # The public URL for the API.
   PORT=8080                                   # The port on which the server will run.
   DB_CONNECTION_STRING=<your-mongodb-connection-string> # MongoDB connection string.
   DB_NAME=taskdb                              # The name of the MongoDB database.
   JWT_SECRET=<your-jwt-secret>                # The secret key for signing JWT tokens.
   JWT_EXPIRATION_IN_SECONDS=86400             # JWT expiration time in seconds (24 hours).
   ```

   Replace `<your-mongodb-connection-string>` and `<your-jwt-secret>` with your MongoDB connection string and a secure JWT secret, respectively.

## Running the Application

To run the application, use:

```bash
make run
```

The application will start a server on port `8080`. You can access the API at `http://localhost:8080/api/v1`.

## Testing

To run tests with coverage, use:

```bash
make test
```

To generate and view a detailed test coverage report, use:

```bash
make coverage
```

## API Endpoints

- **User Authentication**
  - **Register**: `POST /api/v1/auth/register`
  - **Login**: `POST /api/v1/auth/login`
  - **Logout**: `POST /api/v1/auth/logOut`
- **Task Management**
  - **Add Task**: `POST /api/v1/tasks`
  - **Get All Tasks**: `GET /api/v1/tasks`
  - **Get Task by ID**: `GET /api/v1/tasks/{id}`
  - **Update Task**: `PUT /api/v1/tasks/{id}`
  - **Delete Task**: `DELETE /api/v1/tasks/{id}`
- **User Management**
  - **Promote User**: `PATCH /api/v1/users/{username}/promot`

Refer to `docs/api_definition.md` for detailed API usage and request/response formats.
````
