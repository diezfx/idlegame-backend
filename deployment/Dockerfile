# syntax=docker/dockerfile:1.6

FROM golang:1.21 AS builder
ARG APP_NAME=idlegame-backend
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
ENV GOCACHE=/root/.cache/go-build
# Build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/out/app /app/cmd/$APP_NAME/main.go


FROM debian:bookworm-slim as runner
COPY --from=builder /app/out/app /app
COPY ./db/migrations /migrations 
COPY ./config /config
# Run
CMD ["/app"]