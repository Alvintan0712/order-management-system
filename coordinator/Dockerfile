FROM golang:1.23.2 AS builder

WORKDIR /app
COPY ./common /app/common
COPY ./coordinator /app/coordinator

WORKDIR /app/coordinator

RUN go mod tidy
RUN go build -v -o /usr/local/bin/coordinator .

EXPOSE 8800

CMD ["coordinator"]