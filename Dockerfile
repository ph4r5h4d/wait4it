FROM golang:1.16-alpine3.13 as build-env
RUN apk add git gcc
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go run main.go build
FROM alpine:3.13
COPY --from=build-env /app/wait4it .
USER 1001
ENTRYPOINT ["./wait4it"]
