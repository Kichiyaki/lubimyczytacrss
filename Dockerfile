FROM golang:1.18.3-alpine3.16 as builder

WORKDIR /app/lubimyczytacrss

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o lubimyczytacrss .

######## Start a new stage from scratch #######
FROM alpine:3.16

WORKDIR /app/lubimyczytacrss

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/lubimyczytacrss/lubimyczytacrss .

EXPOSE 9234/tcp

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9234/health || exit 1

CMD ["./lubimyczytacrss"]
