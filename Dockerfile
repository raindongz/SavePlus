# build stage 
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz


# run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY db/migration ./migration
RUN apk update && apk add tzdata

EXPOSE 8080
# below line will be overrided if docker compose file have entrypoint
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
