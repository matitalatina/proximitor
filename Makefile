.PHONY: build build-arm

BUILD_CMD=go build -o dist/proximitor cmd/proximitor/main.go

build-arm:
	GOOS=linux GOARCH=arm GOARM=7 $(BUILD_CMD)

build:
	$(BUILD_CMD)

cp-serina:
	scp dist/proximitor serina-rpi:proximitor
	