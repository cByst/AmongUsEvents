FROM golang:alpine as builder

WORKDIR /amongusevents

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /amongusevents/main

FROM scratch

COPY --chown=65534:0 --from=builder /amongusevents/main /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 65534

ENTRYPOINT ["/main"]