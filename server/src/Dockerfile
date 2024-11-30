FROM golang:1.22

COPY . /usr/app
WORKDIR /usr/app
RUN go get
RUN go build backend.go

CMD ./backend