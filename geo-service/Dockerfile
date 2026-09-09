# Stage 1: dependencies (cached unless go.mod/go.sum change)
FROM golang:1.26-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: build both binaries
FROM deps AS build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/worker ./cmd/worker

# Stage 3: runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/api /out/worker ./
EXPOSE 8081
# Default command is the API; docker-compose overrides to /app/worker for the worker container.
CMD ["./api"]
