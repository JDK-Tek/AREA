FROM golang:1.22

WORKDIR /usr/app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o backend backend.go

CMD ./backend
