APP_NAME=unipidns
IMAGE_NAME=unipidns
REGISTRY=containers.cjparish.uk/
VERSION=0.1.0

run:
	go run main.go

compile:
	echo "Compiling..."
	GOOS=darwin GOARCH=arm64 go build -o bin/$(APP_NAME)_macos_arm64
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME)_macos_amd64
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)_linux_amd64
	GOOS=linux GOARCH=arm64 go build -o bin/$(APP_NAME)_linux_arm64
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME)_windows_amd64.exe
	echo "Done."

docker:
	docker build -t $(REGISTRY)$(IMAGE_NAME):$(VERSION) -t $(REGISTRY)$(IMAGE_NAME):latest .

docker-push:
	docker push $(REGISTRY)$(IMAGE_NAME):$(VERSION)
	docker push $(REGISTRY)$(IMAGE_NAME):latest

docker-multiarch:
	docker buildx build --provenance=false --platform linux/amd64,linux/arm64 -t $(REGISTRY)$(IMAGE_NAME):$(VERSION) -t $(REGISTRY)$(IMAGE_NAME):latest -o type=registry .

clean:
	rm -rf bin
	docker rmi $(REGISTRY)$(IMAGE_NAME):$(VERSION) $(REGISTRY)$(IMAGE_NAME):latest

