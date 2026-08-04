FROM golang:1.26.2-bookworm AS builder
WORKDIR /app
COPY go.sum go.mod ./ 
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/cohesive/exe /app/cmd/cohesive

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/cmd/cohesive/exe /app
COPY --from=builder /app/public /app/public
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]
