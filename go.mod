module github.com/openpitrix/helm-package-repository

go 1.12

require (
	github.com/gocraft/dbr v0.0.0-20190714181702-8114670a83bd // indirect
	github.com/golang/protobuf v1.3.2
	helm.sh/helm/v3 v3.0.1 // indirect
	openpitrix.io/openpitrix v0.4.9-0.20200803141156-2fd66cfea5c5
)

replace github.com/gocraft/dbr => github.com/gocraft/dbr v0.0.0-20180507214907-a0fd650918f6

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190423201726-d2cfbce3f3b0

replace helm.sh/helm/v3 => github.com/openpitrix/helm/v3 v3.0.0-20200725015400-ebf6d7e5b2b0
