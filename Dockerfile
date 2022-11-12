FROM golang:1.19-alpine as builder
RUN apk update && apk add --no-cache gcc git

ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/github.com/ph4r5h4d/wait4it
COPY . .
RUN go run main.go build
FROM alpine:3.13
COPY --from=build-env /app/wait4it .
USER 1001
ENTRYPOINT ["./wait4it"]
