

deepcopy-gen \
--input-dirs ./pkg/apis/samplecrd/v1 \
-O zz_generated.deepcopy \
--go-header-file=./hack/boilerplate.go.txt \

