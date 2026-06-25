# ========================
# Variables
# ========================
GO_TEST := go test
GO_TEST_ARGS := -v -cover

IMAGE_NAME ?= gin227/gobek03echo
GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMG_TAG := temporary

ifeq ($(BRANCH), main)
	IMG_TAG := dev
endif

ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

export IMG_TAG

COVERAGE_EXCLUDE=infrastructure|mocks|test|docs|main.go|config.go|client.go|api|helpers|model|mock.go|logger
COVERAGE_THRESHOLD = 80
COVERAGE_FOLDER=./coverage

# ========================
# Run app (dev only)
# ========================
.PHONY: run
run: swagger-gen
	go run ./cmd/api/main.go

# ========================
# Swagger (dev only)
# ========================
.PHONY: swagger-gen
swagger-gen:
	@command -v swag >/dev/null 2>&1 || { echo "❌ swag not installed"; exit 1; }
	swag init -g ./cmd/api/main.go --output ./docs

# ========================
# Unit tests (CI dùng)
# ========================
.PHONY: clean
clean:
	go clean -testcache
	rm -rf $(COVERAGE_FOLDER)

.PHONY: test
test: clean
	mkdir -p $(COVERAGE_FOLDER)

	go test ./... -coverprofile=$(COVERAGE_FOLDER)/coverage.tmp -covermode=atomic -coverpkg=./... -p 1

	grep -v -E "$(COVERAGE_EXCLUDE)" $(COVERAGE_FOLDER)/coverage.tmp > $(COVERAGE_FOLDER)/coverage.out || true
	# 👉 Generate HTML report
	go tool cover -html=$(COVERAGE_FOLDER)/coverage.out -o $(COVERAGE_FOLDER)/coverage.html

	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if awk "BEGIN {exit !($$total >= $(COVERAGE_THRESHOLD))}"; then \
		echo "[PASS] Coverage ($$total%) >= $(COVERAGE_THRESHOLD)%"; \
	else \
		echo "[FAIL] Coverage ($$total%) < $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	fi

# ========================
# Docker
# ========================

# ❗ KHÔNG phụ thuộc swagger-gen (CI sẽ fail nếu chưa cài swag)
.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

# CI-safe login
.PHONY: docker-login
docker-login:
	echo "$(DOCKER_HUB_ACCESS_TOKEN)" | docker login -u "$(DOCKER_HUB_USERNAME)" --password-stdin

.PHONY: docker-release
docker-release:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
.PHONY: docker-test
docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --target test -t bookmark_service:dev --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
    else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
   	fi	
# ========================
# Optional (dev only)
# ========================
.PHONY: test-one
test-one:
ifndef t
	$(error ❌ Missing test name: use t=<NameOfTest>)
endif
ifndef folder
	$(error ❌ Missing package path: use folder=./internal/...)
endif
	$(GO_TEST) $(GO_TEST_ARGS) -run $(t) $(folder)

.PHONY: test-all
test-all:
ifndef folder
	$(error ❌ Missing package path: use folder=./internal/...)
endif
	$(GO_TEST) $(GO_TEST_ARGS) $(folder)

.PHONY: mock-one
mock-one:
ifndef folder
	$(error ❌ Missing folder: use folder=./internal/...)
endif
	cd $(folder) && go generate ./...

.PHONY: mock-all
mock-all:
	go generate ./...
# ========================
# JSON Web Token Key Generator
# ========================
generate-rsa-key:
	openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in private_key.pem -out public_key.pem