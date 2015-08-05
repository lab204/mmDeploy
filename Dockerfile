FROM ckeyer/golang:1.4

WORKDIR /gopath/app
ENV GOPATH /gopath/app

EXPOSE 8080

COPY /main.go $GOPATH/src/main.go

WORKDIR $GOPATH/src/
CMD go run main.go

# CMD cd $GOPATH/src/web-server && ./web-server

