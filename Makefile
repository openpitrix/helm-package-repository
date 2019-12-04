build-image-%: ## build docker image
	if [ "$*" = "latest" ];then \
	docker build -t openpitrix/release-app:latest . ;\
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker build -t openpitrix/release-app:$* . ;\
	fi

push-image-%: ## push docker image
	@if [ "$*" = "latest" ];then \
	docker push openpitrix/release-app:latest; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker push openpitrix/release-app:$*; \
	fi


.PHONY: build
build:
	docker build -t openpitrix/release-app:latest .

.PHONY: debug
debug:
	docker build -t openpitrix/release-app:debug .
	docker push openpitrix/release-app:debug