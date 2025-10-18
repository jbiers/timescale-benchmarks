FROM golang:1.24.9-bookworm AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o benchmarker ./cmd/benchmarker

FROM debian:bookworm-slim

ENV DB_URL=""

WORKDIR /app
COPY --from=builder /build/benchmarker .
ENTRYPOINT ["./benchmarker"]
CMD [""]