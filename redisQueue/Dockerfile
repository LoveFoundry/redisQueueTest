FROM golang:1.25 AS build
WORKDIR /src

# подтягиваем зависимости с кешем
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/redisQueue \
    go mod download

# копируем остальной код
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/app ./cmd/queue

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/app /app/app
USER nonroot:nonroot

EXPOSE 8080


ENV REDIS_ADDR=redis:6379 \
    QUEUE_KEY=queue:messages \
    MAXLEN=50 \
    PORT=8080

ENTRYPOINT ["/app/app"]
