FROM golang:1.12

WORKDIR /go/src/github.com/SummerCash/go-summercash
COPY . .

RUN go get -d -v ./...

CMD go run main.go