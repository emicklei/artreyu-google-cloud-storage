local:
	go build -ldflags "-X main.VERSION=`git rev-parse HEAD` -X main.BUILDDATE=`date -u +%Y:%m:%d.%H:%M:%S`" -o $(GOPATH)/bin/artreyu-gcs
	
build:
	mkdir -p /target
	go build -ldflags "-X main.VERSION='$(VERSION)' -X main.BUILDDATE=`date -u +%Y:%m:%d.%H:%M:%S`" -o /target/artreyu-gcs *.go	
	
# this task exists for Jenkins	
dockerbuild:
	docker build --no-cache=true --tag=artreyu-gcs-builder .
	docker run --rm -e VERSION=$(GIT_COMMIT) -v $(TARGET):/target -t artreyu-gcs-builder
	
# this task exists for local docker	
docker:
	docker build --no-cache=true --tag=artreyu-gcs-builder .
	docker run --rm -v target:/target -t artreyu-gcs-builder	