FROM golang:1.15-alpine3.12 as build-env
RUN apk add git gcc
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o wait4it
FROM alpine:3.11
COPY --from=build-env /app/wait4it .
USER 1001
ENTRYPOINT ["./wait4it"]