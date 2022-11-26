.DEFAULT: all

.PHONY: all
all: gen

.PHONY: gen
gen:
	generate-groups.sh \
	all \
	k8s-controller-custom-resource/pkg/client \
	k8s-controller-custom-resource/pkg/apis \
	samplecrd:v1 \
	--go-header-file=./hack/boilerplate.go.txt \
	--output-base ../
	deepcopy-gen \
    --input-dirs ./pkg/apis/samplecrd/v1 \
    -O zz_generated.deepcopy \
    --go-header-file=./hack/boilerplate.go.txt \

.PHONY: clean
clean:
	rm -rf ./pkg/client
	rm ./pkg/apis/samplecrd/v1/zz_generated.deepcopy.go