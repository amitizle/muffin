FROM golang:1.13-alpine3.11 AS builder
RUN apk --no-cache add git make ca-certificates && \
    update-ca-certificates && \
    adduser -D -g '' appuser && \
    mkdir /muffin
WORKDIR /muffin
COPY . .
RUN GOOS=linux GOARCH=amd64 make build BINARY=/go/bin/muffin

FROM alpine:3.11 AS dumbinit
ARG dumb_init_version=1.2.2
ADD "https://github.com/Yelp/dumb-init/releases/download/v${dumb_init_version}/dumb-init_${dumb_init_version}_amd64" /bin/dumb-init
RUN chmod +x /bin/dumb-init

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/muffin /bin/muffin
COPY --from=dumbinit /bin/dumb-init /bin/dumb-init
USER appuser
ENTRYPOINT ["/bin/dumb-init", "--", "/bin/muffin"]
