NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
#DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

init:
	@echo "$(OK_COLOR)==> This project uses Godep, downloading...$(NO_COLOR)"
	go get github.com/tools/godep
	go get github.com/stretchr/testify
	go get github.com/gocraft/web

exampledocs:
	@echo "$(OK_COLOR)==> Generating Example Docs$(NO_COLOR)"
	go-swaggerLite -apiPackage="github.com/RobotsAndPencils/go-swaggerLite/example" -mainApiFile="github.com/RobotsAndPencils/go-swaggerLite/example/api.go" -basePath="http://127.0.0.1:3000"

test: exampledocs
	@echo "$(OK_COLOR)==> Testing$(NO_COLOR)"
	@echo "$(filter-out $@,$(MAKECMDGOALS))"
	go test -short $(filter-out $@,$(MAKECMDGOALS)) ./...

testv: exampledocs
	@echo "$(OK_COLOR)==> Testing$(NO_COLOR)"
	go test -v -short $(filter-out $@,$(MAKECMDGOALS)) ./...