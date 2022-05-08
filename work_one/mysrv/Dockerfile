FROM golang:1.17 AS build
WORKDIR /data/mysrv/
COPY . .
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN GOOS=linux go build -installsuffix cgo -o httpserver myserver.go

FROM busybox
COPY --from=build /data/mysrv/httpserver /data/mysrv/httpserver
EXPOSE 18086
ENV ENV local
WORKDIR /data/mysrv/
ENTRYPOINT ["./httpserver"]