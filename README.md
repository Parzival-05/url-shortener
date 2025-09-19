# url-shortener
One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

Run with default values:
```
docker-compose -f docker-compose.yml up -d
```
Or fill the .env file if needed.

# Env example

```
PORT=8080
APP_ENV=local

URLSHORTENER_DB_HOST=localhost
# ^^^ localhost for local running & urlshortener for docker ^^^ 
URLSHORTENER_DB_PORT=5432
URLSHORTENER_DB_DATABASE=urlshortener
URLSHORTENER_DB_USERNAME=user
URLSHORTENER_DB_PASSWORD=password1234
URLSHORTENER_DB_SCHEMA=public

SECRET_ALPHABET = P5DRriUYXyL7tHujbQn6lTC2VcKpBf8Zm4vM0EhWzOSFJN1sa3Gdgq9kIxe_Aow
```