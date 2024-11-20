FROM golang:1.23.2 AS builder

WORKDIR /app
COPY ./common /app/common
COPY ./stock /app/stock

WORKDIR /app/stock

RUN go mod tidy
RUN go build -v -o /usr/local/bin/stock .

EXPOSE 8003

CMD ["stock"]