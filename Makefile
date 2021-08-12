SRCS != echo *.go
ROOT != pwd
GO = env GOPATH=$(ROOT)/.gopath go

all: $(SRCS) .gopath
	$(GO) fmt
	$(GO) test
	(cd ttyimg; $(GO) fmt)
	(cd ttyimg; $(GO) build)

.gopath:
	$(GO) mod tidy
	$(GO) get -d

clean:
	rm -rf .gopath ttyimg/ttyimg

