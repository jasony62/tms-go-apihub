FROM golang:1.18.0-alpine3.15

#RUN apk --no-cache add gcc g++ make ca-certificates

RUN go env -w GOPROXY="https://goproxy.cn,direct"

ENV TZ=Asia/Shanghai  
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /home/tms-gah

COPY . .

RUN cd broker;go mod tidy; export CGO_ENABLED=0; go build -buildvcs=false -tags=jsoniter -o tms-gah-broker; ln -s ../example conf
