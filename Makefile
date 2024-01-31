export GOPROXY=https://goproxy.cn,direct

default: build

build: export GO111MODULE=on

build:
	rm -rf bin
	go get
	go build  -o metricsvr main.go

	# mkdir bin
	# cp -rf conf bin/
	# cp -rf metricsvr bin/
	# cp stop.sh bin/
	# cp start.sh bin/

install:
	cp -rf bin/* ../bin/metricsvr/

clean:
	rm -rf bin
	rm metricsvr
