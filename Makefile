REROKU_PACKAGE := github.com/oliw/reroku
VERSION := 0.0.0
BUILD := dev

SRC_DIR := $(GOPATH)/src
REROKU_DIR := $(SRC_DIR)/$(REROKU_PACKAGE)
REROKU_MAIN := $(REROKU_DIR)/reroku
BUILD_DIR := $(REROKU_DIR)/build
REROKU_BIN := $(BUILD_DIR)/bin/reroku

all: $(REROKU_BIN)

clean:
	@rm -rf $(BUILD_DIR)

$(REROKU_BIN):
	@mkdir -p  $(dir $@)
	@(cd $(REROKU_MAIN); go build -o $@)

dpkg: $(REROKU_BIN)
	mkdir -p $(BUILD_DIR)/deb/reroku/usr/local/bin
	cp $(REROKU_BIN) $(BUILD_DIR)/deb/reroku/usr/local/bin
	@(cd $(BUILD_DIR); fpm -s dir -t deb --deb-init $(REROKU_DIR)/packaging/debian/reroku.init -n reroku -v $(VERSION)-$(BUILD) -C deb/reroku .)
