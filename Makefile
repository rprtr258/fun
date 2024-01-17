@_help:
	just --list

lint: # run linter
	@golangci-lint run \
		--exclude-use-default=false \
		--enable revive \
		--enable errcheck \
		--enable govet \
		--enable ineffassign \
		--enable typecheck \
		--disable gosimple \
		--disable staticcheck \
		--disable unused

todo: # check todos
	@rg 'TODO' --glob '**/*.go' || echo 'All done!'

# run tests
test:
	@#go run gotest.tools/gotestsum@latest
	GOSUMDB=sum.golang.org GOEXPERIMENT=rangefunc gotestsum --format dots-v2
	@# go test ./... -count=10 -race

cover: # check opens test cover in browser
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html coverage.out
	@rm coverage.out

ci: lint test # run ci checks

setup: # install git precommit hook
	@echo "#!/bin/env sh\nmake ci" > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
