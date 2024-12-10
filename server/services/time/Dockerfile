FROM golang:1.22

WORKDIR /usr/app
COPY go.mod ./
RUN go mod tidy
COPY . .
RUN go build

CMD ./time
