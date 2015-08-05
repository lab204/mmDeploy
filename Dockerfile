FROM daocloud.io/golang:1.3-onbuild

WORKDIR /gopath/app
ENV GOPATH /gopath/app

EXPOSE 8080

COPY /test.go $GOPATH/src/main.go

WORKDIR $GOPATH/src/
CMD go run main.go

# CMD cd $GOPATH/src/web-server && ./web-server

