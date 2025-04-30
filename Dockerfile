FROM golang:1.23-alpine as builder
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
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/wait4it
RUN chown appuser:appuser /go/bin/wait4it

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/ph4r5h4d/wait4it"
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/wait4it /go/bin/wait4it

USER appuser:appuser
ENTRYPOINT ["/go/bin/wait4it"]
