FROM golang:alpine as builder

WORKDIR /amongushelper

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /amongushelper/main

FROM scratch

COPY --chown=65534:0 --from=builder /amongushelper/main /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 65534

ENTRYPOINT ["/main"]