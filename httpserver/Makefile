export tag=v1.4
root:
	export ROOT=github.com/FuMengAtDJ/K8sLearning

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t mengfu521/httpserver:${tag} .

push: release
	echo "pushing dockerstudy/httpserver"
	docker push mengfu521/httpserver:${tag}
