.DEFAULT: all

.PHONY: all
all: gen

.PHONY: gen
gen:
	hack/update-codegen.sh
	rm -rf ./github.com
	deepcopy-gen \
    --input-dirs ./pkg/apis/samplecontroller/v1 \
    -O zz_generated.deepcopy \
    --go-header-file=./hack/boilerplate.go.txt

.PHONY: clean
clean:
	rm -rf ./pkg/client
	rm -f ./pkg/apis/samplecontroller/v1/zz_generated.deepcopy.go