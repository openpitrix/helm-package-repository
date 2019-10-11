default:
	docker build -t openpitrix/init-repo .
	@echo "ok"

pull:
	docker pull openpitrix/init-repo
	@echo "ok"

test:
	docker run --rm -it -v /data/helm-pkg:/data/helm-pkg sh /usr/local/bin/init-repo.sh

clean:
	@echo "ok"
