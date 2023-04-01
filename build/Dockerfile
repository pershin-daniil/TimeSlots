FROM golang:1.20-alpine AS builder
WORKDIR /app
ADD .. .
RUN go build -o /app/timeslots cmd/timeslots/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder ["/app/timeslots", "/app/timeslots"]

ENTRYPOINT ["/app/timeslots"]