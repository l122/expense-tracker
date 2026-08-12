# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Define build arguments to receive data from the workflow
ARG APP_VERSION=unknown
ARG COMMIT_SHA=unknown
ARG BUILD_DATE=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w \
    -X 'main.AppVersion=${APP_VERSION}' \
    -X 'main.CommitSHA=${COMMIT_SHA}' \
    -X 'main.BuildDate=${BUILD_DATE}'" \
    -o app ./cmd/api

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/app .
# Copy the .env file from the builder stage to the runtime stage
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./app"]