BEAT_NAME=udpproto
BEAT_PATH=github.com/sentient/udpproto
BEAT_URL=https://${BEAT_PATH}
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell govendor list -no-status +local)
PREFIX?=.
NOTICE_FILE=NOTICE

# Path to the libbeat Makefile
-include $(ES_BEATS)/metricbeat/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	$(MAKE) create-metricset
	$(MAKE) collect

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	ln -s ${PWD}/vendor/github.com/elastic/beats/metricbeat/scripts/generate_imports_helper.py ${PWD}/vendor/github.com/elastic/beats/script/generate_imports_helper.py
	rm -rf vendor/github.com/elastic/beats/.git vendor/github.com/elastic/beats/x-pack

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:	
	

.PHONY: generate-proto
generate-proto:
	# protoc --proto_path=protobuf --go_out ./module/udpproto/heap/. ./protobuf/golangheap.proto
	protoc  ./protobuf/golangheap.proto --go_out=testclient/generated --proto_path=protobuf
	protoc  ./protobuf/golangheap.proto --go_out=module/udpproto/heap --proto_path=protobuf

.PHONY: client
client: generate-proto
	go build -o ./build/testclient ./testclient