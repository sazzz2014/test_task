FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

RUN cd /app && \
    go mod tidy && \
    go build -o client ./cmd

CMD ["./client"] 