FROM golang:1.13 AS build-env
ADD . /go/src/demo-server
WORKDIR /go/src/demo-server
RUN ./bash/build/build.sh example
FROM alpine
RUN apk --no-cache add tzdata  && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
COPY --from=build-env /go/src/demo-server/bin /bin
COPY --from=build-env /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
ENV ETCD_ENDPOINTS=10.118.8.250:2379,10.118.9.245:2379,10.118.10.252:2379
ENV ETCD_PREFIX=/root/demo/demo/example/online
ENV ETCD_USER="demo"
ENV ETCD_PWD="demo"
EXPOSE 8086
CMD ["/bin/example"]