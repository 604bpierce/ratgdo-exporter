FROM golang:1.24 AS builder

ENV CGO_ENABLED=0
WORKDIR /app
COPY . .

RUN go build -ldflags="-s -w -extldflags=-static" -o ratgdo-exporter ./cmd/ratgdo-exporter

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/ratgdo-exporter /ratgdo-exporter

EXPOSE 9100
ENTRYPOINT ["/ratgdo-exporter"]
