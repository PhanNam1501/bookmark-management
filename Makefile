# ==========================================
# CẤU HÌNH CƠ BẢN
# ==========================================
APP_NAME         := bookmark_service
MAIN_PATH        := cmd/api/main.go
COVERAGE_EXCLUDE := "mocks|main.go|docs|test"
COVERAGE_DIR     := coverage

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

# ==========================================
# LẤY THÔNG TIN GIT & THIẾT LẬP IMAGE TAG
# ==========================================
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null)
BRANCH  := $(shell git rev-parse --abbrev-ref HEAD)

ifeq ($(BRANCH),main)
    IMAGE_TAG := dev
endif

ifneq ($(GIT_TAG),)
    IMAGE_TAG := $(GIT_TAG)
endif

export IMAGE_TAG

.PHONY: run swagger dev-run test docker-build docker-test clean docker-release docker-login

# ==========================================
# MÔI TRƯỜNG LOCAL
# ==========================================
run:
	go run $(MAIN_PATH)

swagger:
	swag init -g $(MAIN_PATH)

dev-run: swagger run

test:
	@echo "==> Running local tests..."
	@mkdir -p $(COVERAGE_DIR)
	go test ./... -coverprofile=$(COVERAGE_DIR)/coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -v -E $(COVERAGE_EXCLUDE) $(COVERAGE_DIR)/coverage.tmp > $(COVERAGE_DIR)/coverage.out
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@rm -f $(COVERAGE_DIR)/coverage.tmp
	@echo "==> Coverage report generated at $(COVERAGE_DIR)/coverage.html"

# ==========================================
# MÔI TRƯỜNG DOCKER
# ==========================================
docker-build:
	@echo "==> Building Docker image with tag: $(IMAGE_TAG)"
	docker build --target final -t $(APP_NAME):$(IMAGE_TAG) .

docker-test:
	@echo "==> Running tests inside Docker and exporting coverage..."
	docker build --target test --build-arg COVERAGE_EXCLUDE=$(COVERAGE_EXCLUDE) --output type=local,dest=./$(COVERAGE_DIR) .
	@echo "==> Coverage exported to ./$(COVERAGE_DIR)"

docker-login:
	@echo "==> Đang đăng nhập Docker Hub..."
	@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

docker-release:
	@echo "==> Đẩy image $(APP_NAME):$(IMAGE_TAG) lên registry..."
	docker push $(APP_NAME):$(IMAGE_TAG)

# ==========================================
# DỌN DẸP
# ==========================================
clean:
	rm -rf $(COVERAGE_DIR)