build:
	go build
dkbuild: build
	docker build -t reg.qiniu.com/wolfogre/jetbrains-licenser-proxy:${version} .
dkpush:
	docker push reg.qiniu.com/wolfogre/jetbrains-licenser-proxy:${version}
clean:
	rm -f jetbrains-licenser-proxy
