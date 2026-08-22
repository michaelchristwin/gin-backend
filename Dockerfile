# --- build stage ---
FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd

# --- final stage ---
FROM gcr.io/distroless/static-debian13

WORKDIR /app
COPY --from=builder /out/server /app/server

# where the sqlite file will live — mount a volume here
VOLUME ["/data"]
ENV DB_PATH=/data/app.db

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]