FROM golang:1.18.0-alpine3.15

#RUN apk --no-cache add gcc g++ make ca-certificates

RUn go env -w GOPROXY="https://goproxy.cn,direct"

WORKDIR /home/tms-gah

COPY . .

RUN cd broker;go mod tidy; export CGO_ENABLED=0; go build -buildvcs=false -o tms-gah-broker
