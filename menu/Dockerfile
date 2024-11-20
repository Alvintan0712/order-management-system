FROM golang:1.23.2 AS builder

WORKDIR /app
COPY ./common /app/common
COPY ./menu /app/menu

WORKDIR /app/menu

RUN go mod tidy
RUN go build -v -o /usr/local/bin/menu .

EXPOSE 8002

CMD ["menu"]