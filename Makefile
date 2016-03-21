REROKU_PACKAGE := github.com/oliw/reroku

SRC_DIR := $(GOPATH)/src
REROKU_DIR := $(SRC_DIR)/$(REROKU_PACKAGE)
REROKU_MAIN := $(REROKU_DIR)/reroku
BUILD_DIR := $(REROKU_DIR)/build
REROKU_BIN := $(BUILD_DIR)/bin/reroku
VERSION := $(shell cat VERSION)
DEB_PACKAGE_DIR := $(BUILD_DIR)/deb

BUILD_OPTIONS = -ldflags "-X github.com/oliw/reroku/server.VERSION $(VERSION) -X github.com/oliw/reroku/client.VERSION $(VERSION)"

all: $(REROKU_BIN)

clean:
	@rm -rf $(BUILD_DIR)

$(REROKU_BIN):
	@mkdir -p  $(dir $@)
	@(cd $(REROKU_MAIN); go build $(BUILD_OPTIONS) -o $@)

dpkg: $(REROKU_BIN)
	@# place binary
	@mkdir -p $(DEB_PACKAGE_DIR)/usr/local/bin
	@cp $(REROKU_BIN) $(DEB_PACKAGE_DIR)/usr/local/bin
	@# place logrotation settings
	@mkdir -p $(DEB_PACKAGE_DIR)/etc/logrotate.d
	@cp $(REROKU_DIR)/packaging/debian/reroku.logrotate $(DEB_PACKAGE_DIR)/etc/logrotate.d/reroku
	@# place config settings
	@mkdir -p $(DEB_PACKAGE_DIR)/etc
	@cp $(REROKU_DIR)/packaging/debian/reroku.conf $(DEB_PACKAGE_DIR)/etc
	@(cd $(BUILD_DIR); fpm -s dir -t deb --deb-init $(REROKU_DIR)/packaging/debian/reroku.init -n reroku -v $(VERSION) -C $(DEB_PACKAGE_DIR) .)
