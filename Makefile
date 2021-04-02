BUILD = `date +%FT%T%z`
# APP_VERSION := `git describe --tags --abbrev=4`

_CURDIR := `git rev-parse --show-toplevel 2>/dev/null | sed -e 's/(//'`

# PKG := "gitlab.com/ovsinc/memory-rate-limits"
PKG_LIST := $(shell go list ${_CURDIR}/... 2>/dev/null)


linter := golangci-lint

mockery := mockery

test := go test -short -run=^Example


.PHONY: all
all: test lint

.PHONY: build_all
build_all: go_security lint build

.PHONY: full_tests
full_tests: go_lint_max unit_tests race msan ## Full tests

.PHONY: lint
lint: go_lint go_security ## Full liner checks

.PHONY: go_security
go_security: ## Check bugs
	@${linter} run --disable-all \
	-E gosec -E govet -E exportloopref -E staticcheck -E typecheck

.PHONY: go_lint
go_lint: ## Lint the files
	@${linter} run

.PHONY: go_lint_max
go_lint_max: ## Max lint checks the files
	@${linter} run \
	-p bugs -p complexity -p unused -p format \
	-E gosec -E govet -E exportloopref -E staticcheck -E typecheck

.PHONY: go_style
go_style: ## check style of code
	@${linter} run -p style

.PHONY: test
test: unit_test race msan ## Run unittests

.PHONY: unit_test
unit_test: ## Run unittests
	@${test} ${PKG_LIST}

.PHONY: race
race: ## Run data race detector
	@${test} -race ${PKG_LIST}

.PHONY: msan
msan: ## Run memory sanitizer
	@CXX=clang++ CC=clang \
	${test} -msan ${PKG_LIST}

.PHONY: bench
bench: ## Run benchmark tests
	@${test} -benchmem -run=^# -bench=. ${_CURDIR}

.PHONY: bench_vendors
bench_vendors: ## Run benchmarks tests for comparison with other vendors
	@${test} -tags=vendors -benchmem -run=^# -bench=Vendor ${_CURDIR}

.PHONY: coverage
coverage: ## Generate global code coverage report
	@go test -cover ${_CURDIR}

.PHONY: coverhtml
coverhtml: ## Generate global code coverage report in HTML
	[ -x /opt/tools/bin/coverage.sh ] && /opt/tools/bin/coverage.sh html || \
	${_CURDIR}/scripts/tools/coverage.sh html ${BUILD}/coverage.html

.PHONY: dep
dep: ## Get the dependencies
	@go mod vendor


.PHONY: moc_gen
mock_gen: ## Genarate mock-files
	@${mockery} --all --keeptree \
	--dir ${_CURDIR}/log \
	--output ${_CURDIR}/log/mocks


.PHONY: build_docker
build_docker: build ## Pack into docker image
	@docker build -t "test-pospiskud/pospiskud:latest" -f ${BUILD}/docker/pospiskud.docker ${BUILD}




.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
