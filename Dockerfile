FROM golang:1.22.0-alpine3.19 as builder

WORKDIR /build

RUN apk --no-cache add git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY [".", "."]
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o ./bin/app ./cmd/app/main.go

FROM scratch

COPY .env .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/bin/app /

CMD ["/app"]
