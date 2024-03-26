# Todo Api

## Usage

### Setting up the API

#### Build and run Docker image locally

```shell
docker build -t todo-app:latest .
```
```shell
docker run -p 127.0.0.1:8080:8080 todo-app
```

#### Pull and run Docker image from ghcr

- requires the user to have the corresponding permissions in this Github Repository

```shell
docker run -p 127.0.0.1:8080:8080 ghcr.io/oborchardt/todo:main
```

### Get users (not part of the exercise just for convenience)

```shell
curl --location 'localhost:8080/users'
```

### Create a user

```shell
curl --location 'localhost:8080/users' \
--header 'Content-Type: application/json' \
--data '{
    "name": "john",
    "password": "doe"
}'
```

### Login

```shell
curl --location 'localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{
    "name": "john",
    "password": "doe"
}'
```

### Post a Todo

- the Bearer Token in the `Authorization` header must be replaced with the token returned from the `login` route. Each token is valid for 5 minutes.

```shell
curl --location 'localhost:8080/todos' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer c90bce179e3df5820dc5782da71139214fbadae21ffca472330be049997d2b36' \
--data '{
    "title": "cool title",
    "text": "cool text",
    "isDone": false
}'
```

### Get Todos for user

- if the `shared` query parameter is set, all todos are returned including those that are shared with the user
- otherwise only the ones owned by the user are returned

```shell
curl --location 'localhost:8080/todos?shared' \
--header 'Authorization: Bearer 396122dfb31ca5d71cf649b248f272ce5cfad3b9d1fb889a81fba19aa0feaad6'
```

- the number in the path is the ID of the Todo that should be returned

```shell
curl --location 'localhost:8080/todos/1' \
--header 'Authorization: Bearer 396122dfb31ca5d71cf649b248f272ce5cfad3b9d1fb889a81fba19aa0feaad6'
```

### Delete a Todo

```shell
curl --location --request DELETE 'localhost:8080/todos/4' \
--header 'Authorization: Bearer ddc88c08db653805642f90abd95883dcc49a106babcb3c601c568e1a12e648b5'
```

### Update a Todo

- the fields specified in the body are `isDone`, `title` and `text`
- by setting `isDone = true` a Todo can be marked as done

```shell
curl --location --request PATCH 'localhost:8080/todos/5' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer ddc88c08db653805642f90abd95883dcc49a106babcb3c601c568e1a12e648b5' \
--data '{
    "isDone": false
}'
```

### Share a Todo

- `userId` in the request body is the ID of the user the Todo should be shared with

```shell
curl --location 'localhost:8080/todos/1/share' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer d9e144206e170f937d20b7a417164837a5d2d6dd2d2afe907097b562d97e3615' \
--data '{
    "userId": 2
}'
```

### Revoke share of a Todo

```shell
curl --location --request DELETE 'localhost:8080/todos/1/share' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer d9e144206e170f937d20b7a417164837a5d2d6dd2d2afe907097b562d97e3615' \
--data '{
    "userId": 2
}'
```