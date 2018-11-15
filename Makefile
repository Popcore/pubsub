EXECUTABLE=jaakpubsub
BUILD_DIR = build

.PHONY: build
build:
	@echo "=> build executable"
	@export CGO_ENABLED=0
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(EXECUTABLE)