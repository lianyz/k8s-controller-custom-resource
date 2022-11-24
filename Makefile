.DEFAULT: all

.PHONY: all
all: gen

.PHONY: gen
gen:
	./../../kubernetes/code-generator/generate-groups.sh \
	all \
	k8s-controller-custom-resource/pkg/client \
	k8s-controller-custom-resource/pkg/apis \
	samplecrd:v1 \
	--go-header-file=./hack/boilerplate.go.txt \
	--output-base ../

.PHONY: clean
clean:
	rm -rf ./pkg/client