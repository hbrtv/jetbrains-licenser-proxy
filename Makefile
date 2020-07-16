IMAGE_NAME=registry.aliyuncs.com/wolfogre/jetbrains-licenser-proxy
VERSION=$(shell date "+%y.%m").$(shell git rev-list --count --since="$(shell date "+%Y-%m")-01T00:00:00+08:00" HEAD)

check:
	git diff HEAD --quiet || exit 1

build: check
	go build -v -o jetbrains-licenser-proxy

image: build
	docker build -t $(IMAGE_NAME):$(VERSION) .

push: image
	docker push $(IMAGE_NAME):$(VERSION)

run: image
	docker run --rm $(IMAGE_NAME):$(VERSION)

all: image push
