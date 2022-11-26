.DEFAULT: all

.PHONY: all
all: gen

.PHONY: gen
gen:
	hack/update-codegen.sh

.PHONY: clean
clean:
	rm -rf ./pkg/client
	rm -f ./pkg/apis/samplecontroller/v1/zz_generated.deepcopy.go