# Stage 1: dependencies (cached unless go.mod/go.sum change)
FROM golang:1.26-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: build
FROM deps AS build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api

# Stage 3: runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/api .
EXPOSE 8080
ENTRYPOINT ["./api"]
