FROM golang:alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN go build -o -v /app/build/api-server ./cmd/api-server


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app /app

CMD [ "/app/build/api-server" ]