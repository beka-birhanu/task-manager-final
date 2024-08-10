# Task Management Service

A task management service built with Go and MongoDB, providing functionality for adding, updating, deleting, and retrieving tasks. This service offers a RESTful API for managing tasks and demonstrates the use of MongoDB for data storage.

## Project Structure

Here's the updated project structure, grouped by layers and presented in a table format:

### Project Structure

| **Layer**          | **Directory/File**              | **Description**                                    |
| ------------------ | ------------------------------- | -------------------------------------------------- |
| **API**            | `api/controllers/auth`          | Authentication controllers and DTOs.               |
|                    | `api/controllers/auth/dto`      | Data Transfer Objects for authentication.          |
|                    | `api/controllers/base`          | Base controller for shared logic.                  |
|                    | `api/controllers/task`          | Task controllers and DTOs.                         |
|                    | `api/controllers/task/dto`      | Data Transfer Objects for tasks.                   |
|                    | `api/controllers/user`          | User controllers.                                  |
|                    | `api/errors`                    | API error handling and definitions.                |
|                    | `api/middleware/auth`           | Middleware for role-based authentication.          |
|                    | `api/router`                    | Router configuration and initialization.           |
|                    | `api/i_controller.go`           | Interface for controllers.                         |
| **Domain**         | `domain/errors`                 | Domain-specific error definitions.                 |
|                    | `domain/i_hash`                 | Interface for hashing service.                     |
|                    | `domain/models/task`            | Task model definition.                             |
|                    | `domain/models/user`            | User model definition.                             |
| **App**            | `app/common/cqrs/command`       | CQRS command handlers and interfaces.              |
|                    | `app/common/cqrs/query`         | CQRS query handlers and interfaces.                |
|                    | `app/common/i_jwt`              | Interface for JWT service.                         |
|                    | `app/common/i_repo`             | Interface for repositories (task, user).           |
|                    | `app/task/command/add`          | Command for adding a task.                         |
|                    | `app/task/command/delete`       | Command for deleting a task.                       |
|                    | `app/task/command/update`       | Command for updating a task.                       |
|                    | `app/task/query/get`            | Query for retrieving a single task.                |
|                    | `app/task/query/get_all`        | Query for retrieving all tasks.                    |
|                    | `app/user/admin_status/command` | Command for updating user admin status.            |
|                    | `app/user/auth/command`         | Command for user authentication (register, login). |
|                    | `app/user/auth/common`          | Common authentication-related logic and results.   |
|                    | `app/user/auth/query`           | Query for user authentication.                     |
| **Infrastructure** | `infrastructure/db`             | Database connection and setup.                     |
|                    | `infrastructure/hash`           | Hashing service implementation.                    |
|                    | `infrastructure/jwt`            | JWT service implementation.                        |
|                    | `infrastructure/repo/task`      | Task repository implementation.                    |
|                    | `infrastructure/repo/user`      | User repository implementation.                    |
| **Config**         | `config/env.go`                 | Environment variable configuration.                |
| **Documentation**  | `docs/api_definition.md`        | API definition and usage documentation.            |
| **Root**           | `example.env`                   | Example environment file.                          |
|                    | `main.go`                       | Entry point for the application.                   |
|                    | `README.md`                     | Project documentation.                             |
|                    | `go.mod`                        | Go module file for dependency management.          |
|                    | `go.sum`                        | Go module checksums for dependencies.              |

This table groups the files by layers, providing a clear overview of the project structure.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/beka-birhanu/task-manager-authentication.git
   .git
   ```

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

3. **Running the Application**:

   After updating the configuration, you can run the application using:

   ```bash
   go run main.go
   ```

The application will start a server on port `8080`. You can access the API at `http://localhost:8080/api/v1`.

---

This section ensures that users know how to configure the environment variables required to run the application.## Running the Application

To run the application, use:

```bash
go run main.go
```

The application will start a server on port `8080`. You can access the API at `http://localhost:8080/api/v1`.

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
