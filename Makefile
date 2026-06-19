# ========================
# Variables
# ========================
GO_TEST := go test
GO_TEST_ARGS := -v -cover

# ========================
# Run app
# ========================
.PHONY: run
run:swagger-gen
	go run ./cmd/api/main.go
# ========================
# Run app
# ========================
.PHONY: swagger-gen
swagger-gen:
	swag init -g ./cmd/api/main.go --output ./docs
# ========================
# Unit tests
# ========================

# Run single test function
.PHONY: test-one
test-one:
ifndef t
	$(error ❌ Missing test name: use t=<NameOfTest>)
endif
ifndef folder
	$(error ❌ Missing package path: use folder=./internal/...)
endif
	$(GO_TEST) $(GO_TEST_ARGS) -run $(t) $(folder)

# Run all tests in a package
.PHONY: test-all
test-all:
ifndef folder
	$(error ❌ Missing package path: use folder=./internal/...)
endif
	$(GO_TEST) $(GO_TEST_ARGS) $(folder)
# ========================
# Mock generation
# ========================

# Generate mock trong 1 package cụ thể
# Usage: make mock-one folder=./internal/service
.PHONY: mock-one
mock-one:
ifndef folder
	$(error ❌ Missing folder: use folder=./internal/...)
endif
	@echo "→ Generating mocks in $(folder)..."
	cd $(folder) && go generate ./...

# Generate toàn bộ mock trong project
# Usage: make mock-all
.PHONY: mock-all
mock-all:
	@echo "→ Generating all mocks..."
	go generate ./...

# ========================
# Variables
# ========================
IMAGE_NAME ?= your-dockerhub-username/your-image
IMAGE_TAG ?= latest
# ========================
# Docker login
# ========================
.PHONY: docker-login
docker-login:
	docker login
# ========================
# Docker hub image build
# ========================
.PHONY: docker-build
docker-build:swagger-gen
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
# ========================
# Docker hub image push
# ========================
.PHONY: docker-release
docker-release:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

# ========================
# Coverage - unit test (using Ci)
# ========================
COVERAGE_EXCLUDE=infrastructure|mocks|test|docs|main.go|config.go|client.go|api|utils
COVERAGE_THRESHOLD = 80
COVERAGE_FOLDER=./coverage
.PHONY: clean
clean:
	go clean -testcache
	rm -rf $(COVERAGE_FOLDER)
# 	docker rm -f redis || true
.PHONY: test
test: clean
	mkdir -p $(COVERAGE_FOLDER)

	go test ./... -coverprofile=$(COVERAGE_FOLDER)/coverage.tmp -covermode=atomic -coverpkg=./... -p 1

	grep -v -E "$(COVERAGE_EXCLUDE)" $(COVERAGE_FOLDER)/coverage.tmp > $(COVERAGE_FOLDER)/coverage.out || true

	go tool cover -html=$(COVERAGE_FOLDER)/coverage.out -o $(COVERAGE_FOLDER)/coverage.html || true

	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if awk "BEGIN {exit !($$total >= $(COVERAGE_THRESHOLD))}"; then \
		echo "[PASS] Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	else \
		echo "[FAIL] Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	fi