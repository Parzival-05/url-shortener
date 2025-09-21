FROM golang:1.24.6-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk update && apk add make
RUN make build

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}

ENTRYPOINT ["/app/main", "--storage", "postgres"]


