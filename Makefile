clean:  ## Cleans built binaries and packages
	@rm -rf ./bin && mkdir ./bin
	@rm -rf ./pkg && mkdir ./pkg

cleandeps: ## Remove deps source and packages
	@rm -rf ./pkg && mkdir ./pkg
	@rm -rf ./vendor/src && mkdir ./vendor/src

cleanall: clean cleandeps ## Cleans everything

installgb:  ## Installs GB (https://getgb.io/)
	@which gb > /dev/null || go get github.com/constabulary/gb/...

deps: installgb  ## Fetch project dependencies
	@gb vendor restore

listdeps: installgb ## List used dependencies
	@gb vendor list

depadd: installgb ## Add a new dependency.Usage: make depadd DEP_URL=github.com/repo/url DEP_REV=some-revision-tag
	@gb vendor fetch -no-recurse -revision $(DEP_REV) $(DEP_URL)

build: installgb ## Builds sqrible
	@gb build all

help:
	@grep -hE '^[a-zA-Z_-\.]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
.PHONY: install-gb
