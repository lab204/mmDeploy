FROM ckeyer/nginx:1.7
FROM ckeyer/golang:1.4

ENV GOPATH /gopath/app
ENV DATA_DIR /data
WORKDIR $GOPATH/src/

EXPOSE 80

ADD ./server.go $GOPATH/src/server.go

CMD go run server.go
