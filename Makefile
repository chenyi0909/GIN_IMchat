BIN=bin

all: x86

x86 :
	export GOARCH=amd64 GOOS=linux GO111MODULE=on CGO_ENABLED=0  && go build -a -tags netgo -ldflags '-w -extldflags "-static"'  -o $(BIN)/chatProc
	upx $(BIN)/chatProc -1