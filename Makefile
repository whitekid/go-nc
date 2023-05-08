TARGET=bin/gn
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*_test.go")

LDFLAGS ?= "-s -w"
BUILD_FLAGS ?= -v -ldflags ${LDFLAGS}

.PHONY: test

all: build
build: $(TARGET)

$(TARGET): $(SRC)
	go build -o $(TARGET) $(BUILD_FLAGS) -v

linux: ${SRC}
	@${MAKE} ${TARGET} CGO_ENABLED=0 GOOS=linux GOARCH=amd64 BUILD_FLAGS='-ldflags="-s -w"'

clean:
	rm -f ${TARGET}

test:
	go test ${TEST_FLAGS} ./...

dep:
	@rm -f go.mod go.sum
	@go mod init go-nc
	@${MAKE} tidy

tidy:
	go mod tidy
