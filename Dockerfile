FROM golang:1.26.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /task

FROM alpine:3

WORKDIR /app

COPY --from=builder /task .

CMD ["./task"]