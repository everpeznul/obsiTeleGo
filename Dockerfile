FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG REPO
RUN CGO_ENABLED=0 GOOS=linux go build -tags ${REPO} -o /bot ./cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /bot .
CMD ["./bot"]