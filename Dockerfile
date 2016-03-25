FROM golang:1.6-wheezy

RUN go get github.com/emicklei/artreyu

RUN mkdir -p /go/src/github.com/emicklei/artreyu-gcs
WORKDIR /go/src/github.com/emicklei/artreyu-gcs
ADD . /go/src/github.com/emicklei/artreyu-gcs

CMD make build