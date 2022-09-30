lint:
  golangci-lint run \
    --exclude-use-default=false \
    --enable revive \
    --enable deadcode \
    --enable errcheck \
    --enable govet \
    --enable ineffassign \
    --enable structcheck \
    --enable typecheck \
    --enable varcheck \
    --disable gosimple \
    --disable staticcheck \
    --disable unused

todo:
  rg 'TODO' --glob '**/*.go'

test:
  go test ./... -count=10 -race -cover

cover:
  go test -coverprofile=coverage.out ./...
  go tool cover -html coverage.out
  rm coverage.out
