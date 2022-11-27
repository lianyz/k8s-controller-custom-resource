.DEFAULT: all

.PHONY: all
all: gen build

.PHONY: gen
gen:
	hack/update-codegen.sh samplecrd v1

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o bin/samplecrd-controller .

.PHONY: run
run:
	./bin/samplecrd-controller -kubeconfig=$(HOME)/.kube/config -alsologtostderr=true
.PHONY: clean
clean:
	rm -rf ./pkg/client
	rm -f ./pkg/apis/*/*/zz_generated.deepcopy.go