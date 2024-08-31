# Room reservation system

Room reservation system is a API service where you can [create reservations](#create), [list reservations for a room](#list), [get reservation](#get), [delete reservation](#delete), [update reservation](#update). Additionally it has a feature when creating a new reservation, that checks for overlapping reservations for a room, that is if starting and ending time of both reservations intersect.

# Usage

***Requirements***
- [Go 1.20+](https://go.dev/)
- [Docker](https://www.docker.com/)

For accessing and running the API locally follow these steps:


- Clone the repository and cd:
```
	git clone https://github.com/erazr/reservation-system
	cd reservation-system
```

Rename the [.env.example](./.env.example) to .env and change the variables if desired.

- Run the application:

```
	make up
```
OR
```
	docker compose up -d
```

# Tests

For running the tests:
```
	go test ./...
```

Additionally there is a [github action for running tests](./.github/workflows/test.yaml) for each push in repository

# Libraries

- [swaggo](https://github.com/swaggo/swag) to automatically generate RESTful API documentation with Swagger 2.0. 
- [dockertest](https://github.com/ory/dockertest) lets spin up docker containers for integration tests.
- [go-chi](https://github.com/go-chi/chi) lightweight, idiomatic and composable router for building Go HTTP services.
- [pgx](https://github.com/jackc/pgx) simple PostgreSQL driver.
- [golang-migrate](https://github.com/golang-migrate/migrate) for migrating the database in tests.
- [zerolog](https://github.com/rs/zerolog) for structured logging.

# Features

For testing the API and viewing additional documentation visit http://localhost:8080/swagger/index.html

![swagger](https://github.com/user-attachments/assets/c45194d6-2705-451e-a877-eb265313abc4)

## Create

- URL: http://localhost:8080/api/v1/reservations
- Method: POST
- Request Body:

```
	{
		"room_id": "1",
  		"start_time": "29-08-2024 13:00"
  		"end_time": "29-08-2024 14:00",
	}
```

## List

List reservations for a room.

- URL: http://localhost:8080/api/v1/reservations/room/{roomID}
- Method: GET
- Successfull Response:

```
	{
  		"success": true,
  		"data": [
    		{
				"id": "946e2eb89bdc",
      			"room_id": "1",
      			"start_time": "29-08-2024 13:00",
      			"end_time": "29-08-2024 14:00"
    		}
  		]
	}
```

## Get

Get individual reservation

- URL: http://localhost:8080/api/v1/reservations/{ID}
- Method: GET
- Successfull Response:

```
	{
		"success": true,
  		"data": {
			"id": "946e2eb89bdc",
      		"room_id": "1",
      		"start_time": "29-08-2024 13:00",
      		"end_time": "29-08-2024 14:00"
    	}
	}
```

## Delete

- URL: http://localhost:8080/api/v1/reservations/{ID}
- Method: DELETE
- Successfull Response:

```
	204	No Content
```

## Update

- URL: http://localhost:8080/api/v1/reservations/{ID}
- Method: PATCH
- Successfull Response:

```
	204	No Content
```
