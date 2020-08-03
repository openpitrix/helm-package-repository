package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/ptypes/wrappers"

	"openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const Helm = "helm"

func main() {
	logger.SetLevelByString("debug")

	var path string
	flag.StringVar(&path, "path", "./package/", "need package path.eg./your/path/to/pkg/")
	flag.Parse()

	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(nil, "read dir path error: %s", err.Error())
		panic(err)
	}

	var appIds []string

	client, err := app.NewAppManagerClient()
	if err != nil {
		panic(err)
	}
	ctxFunc := func() (ctx context.Context) {
		ctx = context.Background()
		ctx = ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
		return
	}

	for _, f := range fileInfoList {
		filePath := path + f.Name()

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		vi, err := repoiface.LoadPackage(context.Background(), Helm, content)
		appName := vi.GetName()
		versionName := vi.GetVersionName()

		apps, err := client.DescribeActiveApps(ctxFunc(), &pb.DescribeAppsRequest{
			Name: []string{appName},
			Isv:  []string{"\u0000"},
		})
		if err != nil {
			logger.Error(nil, "describe app error: %s", err.Error())
			os.Exit(-1)
		}
		var versionId *wrappers.StringValue
		var appId string
		if apps.TotalCount > 0 {
			appId = apps.AppSet[0].AppId.GetValue()
			describeVersionReq := &pb.DescribeAppVersionsRequest{
				AppId: []string{appId},
				Name:  []string{versionName},
				Limit: 200,
			}
			res, err := client.DescribeActiveAppVersions(ctxFunc(), describeVersionReq)
			if err != nil {
				logger.Error(nil, "describe app version error: %s", err.Error())
				os.Exit(-1)
			}
			if res.TotalCount == 0 {
				createReq := &pb.CreateAppVersionRequest{
					AppId:   pbutil.ToProtoString(appId),
					Package: pbutil.ToProtoBytes(content),
					Type:    pbutil.ToProtoString(Helm),
				}
				res, err := client.CreateAppVersion(ctxFunc(), createReq)
				if err != nil {
					logger.Error(nil, "create app version error: %s", err.Error())
					os.Exit(-1)
				}

				versionId = res.VersionId
			} else {
				logger.Info(nil, "app %v version %s is exists, skip...", appName, versionName)
				appIds = append(appIds, appId)
				continue
			}
		} else {
			createReq := &pb.CreateAppRequest{
				VersionPackage: pbutil.ToProtoBytes(content),
				Name:           pbutil.ToProtoString(appName),
				VersionType:    pbutil.ToProtoString(Helm),
			}
			res, err := client.CreateApp(ctxFunc(), createReq)
			if err != nil {
				logger.Error(nil, "create app error: %s", err.Error())
				os.Exit(-1)
			}

			versionId = res.VersionId
			appId = res.AppId.GetValue()
		}

		submitReq := &pb.SubmitAppVersionRequest{
			VersionId: versionId,
		}

		_, err = client.SubmitAppVersion(ctxFunc(), submitReq)
		if err != nil {
			logger.Error(nil, "submit app version error: %s", err.Error())
			os.Exit(-1)
		}

		passReq := &pb.PassAppVersionRequest{
			VersionId: versionId,
		}
		_, err = client.AdminPassAppVersion(ctxFunc(), passReq)
		if err != nil {
			logger.Error(nil, "pass app version error: %s", err.Error())
			os.Exit(-1)
		}

		releaseReq := &pb.ReleaseAppVersionRequest{
			VersionId: versionId,
		}
		_, err = client.ReleaseAppVersion(ctxFunc(), releaseReq)
		if err != nil {
			logger.Error(nil, "release app verison error: %s", err.Error())
			os.Exit(-1)
		}

		logger.Info(nil, "app %v version %s is done", appName, versionName)
		appIds = append(appIds, appId)
		continue
	}

	_, err = client.ResortApps(ctxFunc(), &pb.ResortAppsRequest{
		AppId: appIds,
	})
	if err != nil {
		logger.Error(nil, "resort apps failed error: %s", err.Error())
		panic(err)
	}
}
