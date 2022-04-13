FROM golang:1.18.0-alpine3.15

RUN apk --no-cache add gcc g++ make ca-certificates

RUn go env -w GOPROXY="https://goproxy.cn,direct"

WORKDIR /home/tms-gah/broker

COPY . .

RUN go mod tidy

RUN go build -o tms-gah-broker
