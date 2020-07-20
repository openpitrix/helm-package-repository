FROM openpitrix/openpitrix-builder as builder

RUN mkdir -p /workspace/helm-package-repository/
WORKDIR /workspace/helm-package-repository/
COPY . .

RUN mkdir -p /release_bin
RUN GOPROXY=https://goproxy.io CGO_ENABLED=0 GOBIN=/release_bin go install -ldflags '-w -s' -tags netgo ./release-app/...

RUN find /release_bin -type f -exec upx {} \;

FROM alpine:3.7
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /release_bin/* /usr/local/bin/
RUN mkdir -p /data/helm-pkg
COPY ./package/* /data/helm-pkg/


RUN apk add --update ca-certificates && \
    update-ca-certificates && \
    adduser -D -g openpitrix -u 1002 openpitrix && \
    chown -R openpitrix:openpitrix /usr/local/bin/ /data/
USER openpitrix

