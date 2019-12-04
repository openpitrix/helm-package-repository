package main

import (
	"context"
	"flag"
	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"sort"
)

func main() {
	Sync()
}

func Sync() {
	var url string
	var scheme string
	flag.StringVar(&url, "r", "", "need chart repository")
	flag.StringVar(&scheme, "t", "", "repository scheme: http/https/s3")
	flag.Parse()
	logger.Info(nil, "repository URL: %s, SCHEME: %s", url, scheme)

	if url == "" || scheme == "" {
		logger.Error(nil, "need parameters[-r,-t]")
		return
	}

	ctxFunc := func() (ctx context.Context) {
		ctx = context.Background()
		ctx = ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
		return
	}

	r := &pb.Repo{
		Type: pbutil.ToProtoString(scheme),
		Url:  pbutil.ToProtoString(url),
	}
	reader, err := repoiface.NewReader(ctxFunc(), r)
	if err != nil {
		logger.Error(nil, "read repository error: %s", err.Error())
	}

	indexFile, err := reader.GetIndex(ctxFunc())
	if err != nil {
		logger.Error(nil, "get index file error: %s", err.Error())
	}

	appClient, _ := app.NewAppManagerClient()
	for appName, appVersions := range indexFile.GetEntries() {
		logger.Info(nil, "create app: %s", appName)
		sort.Sort(appVersions)
		packageName := appVersions[0].GetPackageName()
		versionName := appVersions[0].GetVersionName()
		content, err := reader.ReadFile(ctxFunc(), packageName)
		pkg := strfmt.Base64(content)
		if err != nil {
			logger.Error(nil, "read file error: %s", err.Error())
		}
		createAppReq := &pb.CreateAppRequest{
			Name:           pbutil.ToProtoString(appName),
			VersionType:    pbutil.ToProtoString("helm"),
			VersionPackage: pbutil.ToProtoBytes(pkg),
			VersionName:    pbutil.ToProtoString(versionName),
		}
		createAppResp, err := appClient.CreateApp(ctxFunc(), createAppReq)
		if err != nil {
			logger.Error(nil, "create app error: %s", err.Error())
		}
		AuditApp(ctxFunc, appClient, createAppResp.GetVersionId())
		for _, version := range appVersions[1:] {
			logger.Info(nil, "create app version: %s", version.GetPackageName())
			content, err := reader.ReadFile(ctxFunc(), version.GetPackageName())
			if err != nil {
				logger.Error(nil, "read file error: %s", err.Error())
			}
			pkg := strfmt.Base64(content)
			createAppVersionReq := &pb.CreateAppVersionRequest{
				AppId:       createAppResp.GetAppId(),
				Name:        pbutil.ToProtoString(version.GetVersionName()),
				Description: pbutil.ToProtoString(version.GetDescription()),
				Type:        pbutil.ToProtoString("helm"),
				Package:     pbutil.ToProtoBytes(pkg),
			}
			createAppVersionResp, err := appClient.CreateAppVersion(ctxFunc(), createAppVersionReq)
			if err != nil {
				logger.Error(nil, "create app version error: %s", err.Error())
			}
			AuditApp(ctxFunc, appClient, createAppVersionResp.GetVersionId())
		}
	}
}

func AuditApp(ctxFunc func() context.Context, client *app.Client, versionId *wrappers.StringValue) {
	submitReq := &pb.SubmitAppVersionRequest{
		VersionId: versionId,
	}

	_, err := client.SubmitAppVersion(ctxFunc(), submitReq)
	if err != nil {
		logger.Error(nil, "submit app error: %s", err.Error())
	}

	passReq := &pb.PassAppVersionRequest{
		VersionId: versionId,
	}
	_, err = client.AdminPassAppVersion(ctxFunc(), passReq)
	if err != nil {
		logger.Error(nil, "pass app error: %s", err.Error())
	}

	releaseReq := &pb.ReleaseAppVersionRequest{
		VersionId: versionId,
	}
	_, err = client.ReleaseAppVersion(ctxFunc(), releaseReq)
	if err != nil {
		logger.Error(nil, "release app error: %s", err.Error())
	}
}

func YetSync() {
	ctxFunc := func() (ctx context.Context) {
		ctx = context.Background()
		ctx = ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
		return
	}

	repoClient, _ := repo.NewRepoManagerClient()
	createRepoReq := &pb.CreateRepoRequest{
		Name:             pbutil.ToProtoString("emqqq"),
		Type:             pbutil.ToProtoString("HELM"),
		Url:              pbutil.ToProtoString("http://ec2-13-57-33-89.us-west-1.compute.amazonaws.com/charts/"),
		Credential:       pbutil.ToProtoString(""),
		AppDefaultStatus: pbutil.ToProtoString("active"),
	}
	createRepoResp, _ := repoClient.CreateRepo(ctxFunc(), createRepoReq)

	appClient, _ := app.NewAppManagerClient()
	syncRepoReq := &pb.SyncRepoRequest{
		RepoId: createRepoResp.GetRepoId().String(),
	}
	syncRepoResp, _ := appClient.SyncRepo(ctxFunc(), syncRepoReq)
	_ = syncRepoResp

	describeAppReq := &pb.DescribeAppsRequest{
		RepoId: []string{createRepoResp.GetRepoId().String()},
	}
	describeAppResp, _ := appClient.DescribeActiveApps(ctxFunc(), describeAppReq)
	for _, app := range describeAppResp.AppSet {
		appVersionReq := &pb.DescribeAppVersionsRequest{
			AppId: []string{app.GetAppId().String()},
		}
		appVersions, _ := appClient.DescribeAppVersions(ctxFunc(), appVersionReq)
		for _, version := range appVersions.AppVersionSet {
			versionPackageReq := &pb.GetAppVersionPackageRequest{
				VersionId: version.GetVersionId(),
			}
			versionPackageResp, _ := appClient.GetAppVersionPackage(ctxFunc(), versionPackageReq)
			deleteVersionReq := &pb.DeleteAppVersionRequest{
				VersionId: version.VersionId,
			}
			deleteVersionResp, _ := appClient.DeleteAppVersion(ctxFunc(), deleteVersionReq)
			_ = deleteVersionResp
			createAppVersionReq := &pb.CreateAppVersionRequest{
				AppId:       version.GetAppId(),
				Name:        version.GetName(),
				Description: version.GetDescription(),
				Type:        version.GetType(),
				Package:     pbutil.ToProtoBytes(versionPackageResp.GetPackage()),
			}
			createAppVersionResp, _ := appClient.CreateAppVersion(ctxFunc(), createAppVersionReq)
			_ = createAppVersionResp
		}
	}

}
