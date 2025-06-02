# Backend Golang Coding Test

A simple RESTful API in Golang that manages a list of users. Use MongoDB for persistence, JWT for authentication, and follow clean code practices (Hexagonal Architecture).

---

## Project Setup

1. Clone the repository into your pc
2. Run a MongoDB instance (docker, MongoDB Atlas)
3. Prepare the .env in root directory [The env sample is below]
4. Run the project
5. You can use the postman collection i provide in folder `postman`

## Project structure

```
cmd
└── rest
    └── main.go

internal
├── adapters
│   └── config
│       └── config.go
├── handlers
│   ├── middleware_test.go
│   ├── middleware.go
│   ├── requestValidator.go
│   ├── userHandler_test.go
│   └── userHandler.go
└── helpers
    ├── hash_test.go
    ├── hash.go
    ├── jwt_test.go
    └── jwt.go

storages
└── mongo
    └── repositories
        ├── userRepositoryImpl.go
        └── db.go

core
├── domain
│   ├── user_test.go
│   └── user.go
├── ports
│   └── userRepository.go
└── services
    ├── userService.go
    ├── userServiceImpl_test.go
    └── userServiceImpl.go

postman
└── be-chall-api.postman_collection.json

.env
.gitignore
go.mod
go.sum
README.md

```

## Environment Sample

```
JWT_SECRET_KEY=JWTSECRETKRUB
MONGODB_URI=mongodb://localhost:27017
```

## Run instructions

locate the root directory and run with this command

```
go run .\cmd\rest\main.go
```

## JWT usage guide

1. You can only get jwt token from response of `register endpoint` only
2. Use the token from response and attach to other endpoints before requesting, <br> e.g. `getUserByID`,`getAllUsers`, `updateUserNameAndEmail`, `deleteUser`
3. The token have 1 hour to live

## Endpoints

### Register

for register a new user and get jwt in return

`METHOD POST /register`

#### User Field

| Field    | Type   | Description          | Validation |
| -------- | ------ | -------------------- | ---------- |
| name     | string | name of the user     | required   |
| email    | string | email of the user    | required   |
| password | string | password of the user | required   |

#### Request Body Example

```json
{
  "name": "one1",
  "email": "test@gmail.com",
  "password": "1234abc"
}
```

#### Response

```json
{
  "jwToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjQ1NWRiODMzLTI4NTEtNDhkZi05M2ZmLWM4YjczNDQ0NDcxOCIsIm5hbWUiOiJvbmUxIiwiZW1haWwiOiJ0ZXN0QGdtYWlsLmNvbSIsImV4cCI6MTc0ODg5MDkzMiwiaWF0IjoxNzQ4ODg3MzMyfQ.NW-WW8eanjf3Aj2sdzQup-A2O-jNc-G7CUHj0__aViQ"
}
```

#### Request Body Example (Missing field)

```json
{
  "email": "test@gmail.com",
  "password": "1234abc"
}
```

#### Response

```json
{
  "error": "Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

#### Request Body Example (existing email)

```json
{
  "name": "one1",
  "email": "test@gmail.com",
  "password": "1234abc"
}
```

#### Response

```json
{
  "error": "email already exist"
}
```

### Get User by ID

for fetching user data by ID

`METHOD GET /user/{id}`

- NOTE : you must look up the ID from database

#### Headers

- `Authorization: Bearer <jwtoken>`

#### Response

```json
{
  "id": "455db833-2851-48df-93ff-c8b734444718",
  "name": "one1",
  "email": "test@gmail.com",
  "password": "$2a$10$Rlx6CFM57Oq.5woHDqow6.i96LK6Cm86NIobkh.RSskXGKtiM92g2",
  "created_at": "2025-06-02T18:02:12.065Z"
}
```

### No existing user

#### Response

```json
{
  "error": "mongo: no documents in result"
}
```

### Get all users

for fetching every users in database

`METHOD GET /user`

#### Headers

- `Authorization: Bearer <jwtoken>`

#### Response

```json
[
  {
    "id": "455db833-2851-48df-93ff-c8b734444718",
    "name": "one1",
    "email": "test@gmail.com",
    "password": "$2a$10$Rlx6CFM57Oq.5woHDqow6.i96LK6Cm86NIobkh.RSskXGKtiM92g2",
    "created_at": "2025-06-02T18:02:12.065Z"
  },
  {
    "id": "3e85678a-db5a-4504-a4e0-d3d0e0846fde",
    "name": "one1",
    "email": "test2@gmail.com",
    "password": "$2a$10$/uX2BOOiH904KtATyKD5ZesCSCR3.FFrgmy.HGjREuLpeXTzmMkG6",
    "created_at": "2025-06-02T18:09:07.658Z"
  }
]
```

### No existing user

#### Response

```json
{
  "error": "mongo: no documents in result"
}
```

### Update user's email or name

for update user's email or name

`METHOD PATCH /user/{id}`

- NOTE : you must look up the ID from database

#### Headers

- `Authorization: Bearer <jwtoken>`

#### User Field

| Field | Type   | Description       | Validation                                 |
| ----- | ------ | ----------------- | ------------------------------------------ |
| name  | string | name of the user  | optinal, required if email is not provided |
| email | string | email of the user | optinal, required of name is not provided  |

#### Request Body Example

```json
{
  "email": "test@gmail.com",
  "name": "test"
}
```

#### Response

```json
{
  "message": "User updated successfully"
}
```

#### Request Body Example (Invalid email format)

```json
{
  "email": "testgmail.com",
  "password": "test"
}
```

#### Response

```json
{
  "error": "Key: 'EditUser.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

#### Request Body Example (not provide body)

```json
{}
```

#### Response

```json
{
  "error": "name and email cannot be empty"
}
```

#### Request Body Example (existing email)

```json
{
  "email": "testgmail.com",
  "password": "test"
}
```

#### Response

```json
{
  "error": "email already exist"
}
```

### Delete user by ID

for deleting user from database

`METHOD DELETE /user/{id}`

- NOTE : you must look up the ID from database

#### Headers

- `Authorization: Bearer <jwtoken>`

### Valid user ID

#### Response

```json
{
  "message": "User deleted successfully"
}
```

### No existing user

#### Response

```json
{
  "error": "mongo: no documents in result"
}
```

## Submission Guidelines

- Submit a GitHub repo or zip file.
- Include a `README.md` with:
  - Project setup and run instructions
  - JWT token usage guide
  - Sample API requests/responses
  - Any assumptions or decisions made

---

## Evaluation Criteria
