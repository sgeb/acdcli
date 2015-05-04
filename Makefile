# This file was inspired by hashicorp's terraform at:
# https://github.com/hashicorp/terraform
#
TEST?=./...
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

default: test

# bin generates the releaseable binaries
bin:
	@sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev:
	@ACD_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

quickdev:
	@ACD_QUICKDEV=1 ACD_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# test runs the unit tests and vets the code
test:
	ACD_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4
	@$(MAKE) vet

# testacc runs acceptance tests
testacc:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	ACD_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 45m

# testrace runs the race checker
testrace:
	ACD_ACC= go test -race $(TEST) $(TESTARGS)

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@go list -f '{{.Dir}}' ./... \
		| grep -v '.*github.com/sgeb/acdcli$$' \
		| xargs go tool vet ; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Vet found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

clean:
	rm -rf bin
	rm -rf pkg

.PHONY: bin default test vet clean
