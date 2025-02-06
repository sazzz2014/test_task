FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN cd /app && \
    go mod init server && \
    go mod tidy && \
    go build -o server ./cmd

EXPOSE 8080

CMD ["./server"] 