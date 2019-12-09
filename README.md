## Jobs for openpitrix
>1. sync-repo
> * what: 同步第三方仓库应用[应用信息以及package文件]到内置仓库[minio];

> * howto:
> ```
> 1.修改openpitrx-sync-app-job.yaml文件中.spec.template.spec.containers.command里边的-r,-t参数.
> eg.["sync-repo","-r","https://kubernetes-charts.storage.googleapis.com/","-t","https"].
> 2.apply -f openpitrx-sync-app-job.yaml
>```
> * result: 第三方仓库的App会显示在ks app store;