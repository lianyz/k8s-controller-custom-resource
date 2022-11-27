.DEFAULT: all

.PHONY: all
all: gen build

.PHONY: gen
gen:
	hack/update-codegen.sh samplecrd v1

.PHONY: build
build:
	go build -o bin/samplecrd-controller .

.PHONY: clean
clean:
	rm -rf ./pkg/client
	rm -f ./pkg/apis/*/*/zz_generated.deepcopy.go