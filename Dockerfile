FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o app cmd/server/main.go

FROM alpine
WORKDIR /build
COPY --from=builder /build/app /build/app
COPY --from=builder /build/configs /build/configs
COPY --from=builder /build/migration.sql /build/migration.sql
CMD ["/build/app"]