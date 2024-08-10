## Task Management API

This API allows users to manage tasks and handle user authentication. It provides endpoints to create, update, delete, and retrieve tasks, along with user registration, login, and logout.

### API Endpoints

#### **Task Management**

- **Create Task**: `POST /api/v1/tasks`

  - **Request Body**:

    ```json
    {
      "title": "string",
      "description": "string",
      "dueDate": "string (ISO 8601 format)",
      "status": "string"
    }
    ```

  - **Response**:
    - `201 Created`
    - **Headers**: `Location: /api/v1/tasks/{id}`

- **Update Task**: `PUT /api/v1/tasks/{id}`

  - **Path Parameters**: `{id}` (UUID)
  - **Request Body**:

    ```json
    {
      "title": "string",
      "description": "string",
      "dueDate": "string (ISO 8601 format)",
      "status": "string"
    }
    ```

  - **Response**: `200 OK`

- **Delete Task**: `DELETE /api/v1/tasks/{id}`

  - **Path Parameters**: `{id}` (UUID)
  - **Response**: `200 OK`

- **Get All Tasks**: `GET /api/v1/tasks`

  - **Response**:
    ```json
    [
      {
        "id": "uuid",
        "title": "string",
        "description": "string",
        "dueDate": "string (ISO 8601 format)",
        "status": "string"
      }
    ]
    ```

- **Get Single Task**: `GET /api/v1/tasks/{id}`
  - **Path Parameters**: `{id}` (UUID)
  - **Response**:
    ```json
    {
      "id": "uuid",
      "title": "string",
      "description": "string",
      "dueDate": "string (ISO 8601 format)",
      "status": "string"
    }
    ```

#### **User Management**

- **Create User**: `POST /api/v1/users`

  - **Request Body**:
    ```json
    {
      "username": "beka_birhanu",
      "password": "************",
      "isAdmin": true
    }
    ```
  - **Response**: `201 Created`

    ```json
    {
      "id": "00000000-0000-0000-0000-000000000000",
      "username": "beka_birhanu",
      "isAdmin": true
    }
    ```

    **Headers**: `Set-Cookie: token=<token_value>; HttpOnly; Secure`

#### **Authentication**

- **Sign in**: `POST /api/v1/auth/login`

  - **Request Body**:
    ```json
    {
      "username": "beka_birhanu",
      "password": "************"
    }
    ```
  - **Response**: `200 OK`

    ```json
    {
      "id": "00000000-0000-0000-0000-000000000000",
      "username": "beka_birhanu",
      "isAdmin": true
    }
    ```

    **Headers**: `Set-Cookie: token=<token_value>; HttpOnly; Secure`

- **Sign out**: `POST /api/v1/auth/logOut`
  - **Headers**: `Cookie: token=<token_value>`
  - **Response**: `204 No Content`
    **Headers**: `Set-Cookie: token=; HttpOnly; Secure; Max-Age=0`

---

This documentation reflects the updated API structure, including authentication and user management consistent with the other project you mentioned.
