FROM golang:1.18.3-alpine3.16

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD server ./server
COPY cmd/40/main.go ./

RUN go build -o api .

EXPOSE 5000

CMD ["./api"]