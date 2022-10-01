@_help:
  just --list

# run linter
@lint:
  golangci-lint run \
    --exclude-use-default=false \
    --enable revive \
    --enable errcheck \
    --enable govet \
    --enable ineffassign \
    --enable typecheck \
    --disable gosimple \
    --disable staticcheck \
    --disable unused

# check todos
@todo:
  rg 'TODO' --glob '**/*.go'

# run tests
test:
  go test ./... -count=10 -race -cover

# check opens test cover in browser
@cover:
  go test -coverprofile=coverage.out ./...
  go tool cover -html coverage.out
  rm coverage.out

# run ci checks
@ci: lint test

# install git precommit hook
@setup:
  echo "#!/bin/env just ci" > .git/hooks/precommit
  chmod +x .git/hooks/precommit
