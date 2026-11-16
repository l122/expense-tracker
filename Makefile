# Makefile for Go project

BINARY_NAME=main.exe
TMP_DIR=tmp
BINARY_PATH=$(TMP_DIR)/$(BINARY_NAME)

# Detect operating system and set platform-specific commands
ifeq ($(OS),Windows_NT)
    RM_CMD = if exist "$(TMP_DIR)" rmdir /s /q "$(TMP_DIR)" >nul 2>&1
    KILL_CMD = taskkill /f /im $(BINARY_NAME) >nul 2>&1 || (echo Process not running)
    RUN_BG = start /b "" "$(BINARY_PATH)"
    MKDIR_CMD = if not exist "$(TMP_DIR)" mkdir "$(TMP_DIR)"
else
    RM_CMD = rm -rf $(TMP_DIR)
    KILL_CMD = -pkill -SIGTERM -f "./$(BINARY_PATH)"
    RUN_BG = ./$(BINARY_PATH) &
    MKDIR_CMD = mkdir -p "$(TMP_DIR)"
endif

.PHONY: all build run clean test start stop watch

all: build test

build:
	@echo "Building..."
	@$(MKDIR_CMD)
	@go build -o $(BINARY_PATH) cmd/api/main.go
	@echo "Built!"

run: build
	@echo "Starting application..."
	@$(RUN_BG)
	@echo "Started!"

clean:
	@echo "Cleaning..."
	@go clean
	@$(RM_CMD)
	@echo "Cleaned!"

test:
	@echo "Testing..."
	@go test ./... -v
	@echo "Done!"

start: stop run

stop:
	@echo "Stopping application..."
	@$(KILL_CMD)
	@echo "Stopped!"

# Live Reload using Air
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"
