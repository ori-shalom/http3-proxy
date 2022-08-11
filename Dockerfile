# syntax=docker/dockerfile:1.3

############ BUILDER ############

ARG GO_VERSION=1.19
ARG ALPINE_VERSION=3.16

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

# Disable CGO.
ENV CGO_ENABLED=0

# Copy and build project.
COPY . /app

WORKDIR /app
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go build -o main cmd/main.go

############ RUNTIME ############
FROM alpine:${ALPINE_VERSION} as runtime

# Create non-root user and group.
ARG USER=default
ARG GROUP=default
# -S: Create a system user
# -D: Do not assign a password
# -H: Do not create home directory
RUN addgroup -S $GROUP && adduser -S -D -H -G $GROUP $USER

# Use non-root user.
USER $USER

# Copy app.
COPY --from=builder --chown=$USER /app/main /app/main

# Workdir.
WORKDIR /app

ENV PORT 8080
ENV TARGET_HOST 172.17.0.1

ENTRYPOINT ["/app/main"]
