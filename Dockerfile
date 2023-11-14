FROM golang:1.20.10-alpine3.18 AS builder

WORKDIR /build
ENV GOPROXY=https://goproxy.cn,direct

COPY . .
RUN go mod download
RUN go build -ldflags '-w -s' -o openapi-cli ./cmd/openapi-cli

FROM alpine:3.18.4

COPY --from=builder /build/openapi-cli openapi-cli

CMD ["/openapi-cli v"]
