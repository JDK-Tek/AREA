FROM golang:1.23.4

WORKDIR /usr/app
COPY go.mod ./
RUN go mod tidy
COPY . .
RUN go build

CMD ./google
