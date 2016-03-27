FROM golang:1.6-wheezy

# until vendored
RUN go get github.com/emicklei/artreyu
RUN go get golang.org/x/net/context
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/storage/v1

RUN mkdir -p /go/src/github.com/emicklei/artreyu-google-cloud-storage
WORKDIR /go/src/github.com/emicklei/artreyu-google-cloud-storage
ADD . /go/src/github.com/emicklei/artreyu-google-cloud-storage

CMD make build