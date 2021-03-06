# ---------------------  dev (build) image --------------------- #

FROM golang:1.12-alpine as builder

RUN apk add git
RUN apk add make

RUN mkdir -p /opt/gobetween
WORKDIR /opt/gobetween

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make build-static

# --------------------- final image --------------------- #

FROM scratch

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /opt/gobetween/bin/gobetween  .

CMD ["/gobetween", "-c", "/etc/gobetween/conf/gobetween.toml"]

LABEL org.label-schema.vendor="gobetween" \
      org.label-schema.url="http://gobetween.io" \
      org.label-schema.name="gobetween" \
      org.label-schema.description="Modern & minimalistic load balancer for the Сloud era"
