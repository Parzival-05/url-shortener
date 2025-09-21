# url-shortener

A simple URL shortening service written in Go. It provides a clean HTTP API to create unique, non-sequential short links and retrieve the original URL from a shortened code.

## How It Works: The Alphabet is Your Secret Key
This service uses the Sqids library to generate short codes. It's crucial to understand that Sqids does not use a traditional "secret key" or "salt" parameter.

Instead, your entire configuration-primarily the character order of your alphabet-acts as your secret.

The Sqids algorithm deterministically shuffles the alphabet based on its exact structure. Even a tiny change to the alphabet order will produce a completely different output for the same input number. To reverse-engineer your IDs, an attacker would need to guess the exact permutation of your 63-character alphabet.


## Getting Started

1. Clone
   ```bash
   git clone https://github.com/Parzival-05/url-shortener
   cd url-shortener
   ```

2. Configure your environment:
   
   Create a `.env` file in the project root (you can copy `example.env`).

3. Install dependencies:
   ```
   go mod tidy
   ```
4. Run the service (with postgres as storage):
   ```
    make docker-run
   ```
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
make run ARGS="--storage postgres"
``` 
or 
```bash
make run ARGS="--storage inmemory"
``` 

Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
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


# API Reference

Check out Swagger UI to explore the API: http://localhost:8080/swagger/index.html
(use your actual port instead of 8080)