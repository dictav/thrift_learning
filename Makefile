all:
	thrift -r --gen rb --out ruby/gen-rb tutorial.thrift
	thrift -r --gen js --out browser/gen-js tutorial.thrift
	thrift -r --gen js:node --out nodejs/gen-nodejs tutorial.thrift
	thrift -r --gen go:package_prefix=github.com/dictav/thrift_learning/go/gen-go/ --out go/gen-go tutorial.thrift
