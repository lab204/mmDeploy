FROM ckeyer/nginx:1.7
FROM ckeyer/golang:1.4

ENV GOPATH /gopath/app
ENV DATA /data
WORKDIR $GOPATH/src/

EXPOSE 80

ADD ./nginx.conf /etc/nginx/nginx.conf
ADD ./server.go $GOPATH/src/server.go
ADD ./test.html /data/test.html

CMD go run server.go
