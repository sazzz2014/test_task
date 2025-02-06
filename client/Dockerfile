FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN cd /app && \
    go mod init client && \
    go mod tidy && \
    go build -o client ./cmd

CMD ["./client"] 