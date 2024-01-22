FROM golang:alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN make build/app path=/app/build


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app /app

CMD [ "/app/build/api-server" ]