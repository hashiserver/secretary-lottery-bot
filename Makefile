.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./build/main

build:
	GOOS=linux GOARCH=amd64 go build -o build/main ./src

package:
	sam package --template-file template.yaml --output-template-file output-template.yaml --s3-bucket <S3BUCKETNAME> --profile <YOURENAME>

deploy:
	sam deploy --template-file output-template.yaml --stack-name <APPNAME> --capabilities CAPABILITY_IAM --profile <YOURENAME>
