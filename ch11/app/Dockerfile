FROM golang:1.22

WORKDIR /demo

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app

RUN chmod +x ./app

EXPOSE 8080

CMD ["./app"]