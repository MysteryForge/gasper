
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git gcc musl-dev make
WORKDIR /build
# Install xk6 v1.9.3
RUN go install go.k6.io/xk6/cmd/xk6@a9915b8e1519a26cbfbbafb93cd4159ff0e617e8
COPY . .
RUN make build

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/bin/gasper /usr/local/bin/gasper
WORKDIR /scripts
ENTRYPOINT ["gasper"]
CMD ["run", "--help"]

