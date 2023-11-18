SHELL = /bin/bash -o pipefail

.PHONY: test test-ci lint lint-ci fmt fmt-ci clean release lint-docs audit encrypt decrypt sops

test:
	rm -rf coverage
	mkdir -p coverage
	GOCOVERDIR=coverage go run -race -cover -covermode atomic _examples/single.go -c _examples/single.yml | egrep -q '^[0-9]{2}:[0-9]{2} INF Hello world! program=single$$'
	GOCOVERDIR=coverage go run -race -cover -covermode atomic _examples/multi.go --logging.console.type=nocolor plus 1 1 | egrep -q '^[0-9]{2}:[0-9]{2} INF 2$$'
	go tool covdata percent -i=coverage -pkg=gitlab.com/tozd/go/cli

test-ci: test
	go tool covdata textfmt -i=coverage -pkg=gitlab.com/tozd/go/cli -o coverage.txt
	gocover-cobertura < coverage.txt > coverage.xml
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run --timeout 4m --color always --allow-parallel-runners --fix

lint-ci:
	golangci-lint run --timeout 4m --out-format colored-line-number,code-climate:codeclimate.json

fmt:
	go mod tidy
	git ls-files --cached --modified --other --exclude-standard -z | grep -z -Z '.go$$' | xargs -0 gofumpt -w
	git ls-files --cached --modified --other --exclude-standard -z | grep -z -Z '.go$$' | xargs -0 goimports -w -local gitlab.com/tozd/go/cli

fmt-ci: fmt
	git diff --exit-code --color=always

clean:
	rm -rf coverage.* codeclimate.json tests.xml coverage

release:
	npx --yes --package 'release-it@15.4.2' --package '@release-it/keep-a-changelog@3.1.0' -- release-it

lint-docs:
	npx --yes --package 'markdownlint-cli@~0.34.0' -- markdownlint --ignore-path .gitignore --ignore testdata/ '**/*.md'

audit:
	go list -json -deps ./... | nancy sleuth --skip-update-check

encrypt:
	gitlab-config sops --encrypt --mac-only-encrypted --in-place --encrypted-comment-regex sops:enc .gitlab-conf.yml

decrypt:
	SOPS_AGE_KEY_FILE=keys.txt gitlab-config sops --decrypt --in-place .gitlab-conf.yml

sops:
	SOPS_AGE_KEY_FILE=keys.txt gitlab-config sops .gitlab-conf.yml
