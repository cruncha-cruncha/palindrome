FROM golang:1.23-alpine
WORKDIR /usr/local/app

COPY . .

RUN go mod download
RUN go build -o main .

#ENV P_DELAY=10

CMD ["./main"]

EXPOSE 8090