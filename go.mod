module github.com/openpitrix/helm-package-repository

go 1.12

require (
	github.com/gocraft/dbr v0.0.0-20190714181702-8114670a83bd // indirect
	github.com/golang/protobuf v1.3.2
	openpitrix.io/openpitrix v0.4.9-0.20200617102217-10d232395f06
)

replace github.com/gocraft/dbr => github.com/gocraft/dbr v0.0.0-20180507214907-a0fd650918f6

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190423201726-d2cfbce3f3b0
