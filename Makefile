VERSION := $(shell cat VERSION)

.DEFAULT_GOAL := test

bin/$(NAME): $(SRCS)
	go build -o bin/$(NAME)

.PHONY: test
test:
	go test -count=1 $${TEST_ARGS}

.PHONY: testrace
testrace:
	go test -count=1 $${TEST_ARGS} -race

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: fmtci
fmtci:
	! gofmt -d . | grep '^'

.PHONY: lint
lint:
	golint -set_exit_status

.PHONY: vet
vet:
	go vet

.PHONY: tag
tag:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --tags

.PHONY: release
release: tag
	git push origin master
