FROM golang:1.6-wheezy

RUN mkdir -p /go/src/github.com/emicklei/artreyu-google-cloud-storage
WORKDIR /go/src/github.com/emicklei/artreyu-google-cloud-storage
ADD . /go/src/github.com/emicklei/artreyu-google-cloud-storage

CMD make build