FROM golang:1.19

# Copy and build project.
COPY . /app

WORKDIR /app
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go build -o main cmd/main.go

ENV PORT 8080
ENV TARGET_HOST 172.17.0.1

ENTRYPOINT ["/app/main"]
