BUILDX_BUILD_PUSH=docker buildx build --platform linux/amd64,linux/arm64 --output=type=registry --push

build-push-image-%: ## build docker image
	if [ "$*" = "latest" ];then \
	$(BUILDX_BUILD_PUSH) -t openpitrix/release-app:latest . && \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	$(BUILDX_BUILD_PUSH) -t openpitrix/release-app:$* . && \
	fi

.PHONY: build
build:
	docker build -t openpitrix/release-app:latest .