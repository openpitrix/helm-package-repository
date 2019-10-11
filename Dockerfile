FROM alpine:3.7

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache util-linux jq curl

RUN mkdir -p /data/helm-pkg
COPY ./package/* /data/helm-pkg/
COPY init_repo.sh /usr/local/bin/
COPY config.json /usr/local/bin/
RUN chmod +x /usr/local/bin/init_repo.sh
COPY --from=openpitrix/openpitrix:latest /usr/local/bin/opctl /usr/local/bin/

#ENTRYPOINT ["init_repo.sh"]
