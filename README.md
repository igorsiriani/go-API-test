# GoLang API Project

This project is a GoLang API that reads a CSV file, creates a database with the information in the file, provides user CRUD and authentication endpoints, and allows searching in the text for matches or value intervals of all 52 columns available.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Search Criteria](#search-criteria)
- [License](#license)

## Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/yourusername/your-repo-name.git
    cd your-repo-name
    ```

2. **Initialize Go module:**

    ```sh
    go mod init your_project_name
    ```

3. **Install dependencies:**

    ```sh
    go get -u github.com/gin-gonic/gin
    go get -u github.com/jinzhu/gorm
    go get -u github.com/jinzhu/gorm/dialects/sqlite
    go get -u golang.org/x/crypto/bcrypt
    go get -u github.com/dgrijalva/jwt-go
    ```

4. **Run the application:**

    ```sh
    go run main.go models.go
    ```

## Usage

1. **Run the server:**

    ```sh
    go run main.go models.go
    ```

2. **Use Postman or any other API client to interact with the API.**

## API Endpoints

### User Registration

- **Endpoint:** `POST /register`
- **Description:** Register a new user.
- **Request Body:**

    ```json
    {
        "username": "testuser",
        "password": "password123"
    }
    ```

### User Login

- **Endpoint:** `POST /login`
- **Description:** Log in a user to receive a JWT token.
- **Request Body:**

    ```json
    {
        "username": "testuser",
        "password": "password123"
    }
    ```

- **Response:**

    ```json
    {
        "token": "your_jwt_token"
    }
    ```

### Get User

- **Endpoint:** `GET /user`
- **Description:** Get the authenticated user's details.
- **Headers:** `Authorization: Bearer <your_jwt_token>`

### Update User

- **Endpoint:** `PUT /user`
- **Description:** Update the authenticated user's details.
- **Headers:** `Authorization: Bearer <your_jwt_token>`
- **Request Body:**

    ```json
    {
        "username": "newusername",
        "password": "newpassword123"
    }
    ```

### Delete User

- **Endpoint:** `DELETE /user`
- **Description:** Delete the authenticated user's account.
- **Headers:** `Authorization: Bearer <your_jwt_token>`

### Search Records

- **Endpoint:** `POST /search`
- **Description:** Search records based on multiple criteria.
- **Headers:** `Authorization: Bearer <your_jwt_token>`
- **Request Body:**

    ```json
    {
        "criteria": [
            {
                "column": "Column1",
                "value1": "exampleValue"
            },
            {
                "column": "Column2",
                "value1": "10",
                "value2": "20"
            }
        ]
    }
    ```

## Search Criteria

- **column:** The column to search in.
- **value1:** The value to match or the start value for an interval.
- **value2:** The end value for an interval (optional).

## License

This project is licensed under the MIT License.
