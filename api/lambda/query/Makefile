build-ApiFunction:
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap main/main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.
	cp ./prompt.tmpl $(ARTIFACTS_DIR)/.
