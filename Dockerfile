FROM golang:1.23-alpine
WORKDIR /usr/local/app

COPY . .

RUN go mod download
RUN go build -o main .

CMD ["./main"]

EXPOSE 8090