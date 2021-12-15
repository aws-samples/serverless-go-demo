STACK_NAME ?= serverless-go-demo
FUNCTIONS := get-products get-product put-product delete-product products-stream
REGION := eu-central-1

# To try different version of Go
GO := go

# Make sure to install aarch64 GCC compilers if you want to compile with GCC.
CC := aarch64-linux-gnu-gcc
GCCGO := aarch64-linux-gnu-gccgo-10

ci: build tests-unit

build:
		${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-${function})

build-%:
		cd functions/$* && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 ${GO} build -o bootstrap

build-gcc:
		${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-gcc-${function})

build-gcc-%:
		cd functions/$* && GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=${CC} ${GO} build -o bootstrap

build-gcc-optimized:
		${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-gcc-optimized-${function})

build-gcc-optimized-%:
		cd functions/$* && GOOS=linux GOARCH=arm64 GCCGO=${GCCGO} ${GO} build -compiler gccgo -gccgoflags '-static -Ofast -march=armv8.2-a+fp16+rcpc+dotprod+crypto -mtune=neoverse-n1 -moutline-atomics' -o bootstrap

invoke:
	@sam local invoke --env-vars env-vars.json GetProductsFunction

invoke-put:
	@sam local invoke --env-vars env-vars.json --event functions/put-product/event.json PutProductFunction

invoke-get:
	@sam local invoke --env-vars env-vars.json --event functions/get-product/event.json GetProductFunction

invoke-delete:
	@sam local invoke --env-vars env-vars.json --event functions/delete-product/event.json DeleteProductFunction

invoke-stream:
	@sam local invoke --env-vars env-vars.json --event functions/products-stream/event.json DDBStreamsFunction

clean:
	@rm $(foreach function,${FUNCTIONS}, functions/${function}/bootstrap)

deploy:
	if [ -f samconfig.toml ]; \
		then sam deploy --stack-name ${STACK_NAME}; \
		else sam deploy -g --stack-name ${STACK_NAME}; \
  fi

tests-unit:
	@go test -v -tags=unit -bench=. -benchmem -cover ./...

tests-integ:
	API_URL=$$(aws cloudformation describe-stacks --stack-name $(STACK_NAME) \
	  --region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text) go test -v -tags=integration ./...

tests-load:
	API_URL=$$(aws cloudformation describe-stacks --stack-name $(STACK_NAME) \
	  --region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text) artillery run load-testing/load-test.yml
