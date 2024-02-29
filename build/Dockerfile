FROM golang:alpine AS builder

RUN apk add --no-cache make gcc g++

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -v -o /app/build/gotube ./cmd/gotube


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app /app

CMD [ "/app/build/gotube" ]
