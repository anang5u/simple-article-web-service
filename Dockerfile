# golang:1.20.6 sebagai base image
FROM golang:1.20.6

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]
