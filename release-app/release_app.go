package main

import (
	"context"
	"flag"
	"github.com/go-openapi/strfmt"
	"io/ioutil"
	"openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"os"
	"strings"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "./package/", "need package path.eg./your/path/to/pkg/")
	flag.Parse()

	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(nil, "read dir path error: %s", err.Error())
	}

	for _, f := range fileInfoList {
		filePath := path + f.Name()
		var appName string
		segName := strings.Split(f.Name(), "-")
		if len(segName) > 0 {
			appName = strings.Join(segName[:len(segName)-1], "-")
		}
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			logger.Error(nil, "")
		}
		pkg := strfmt.Base64(content)

		client, err := app.NewAppManagerClient()
		if err != nil {
			panic(err)
		}

		createReq := &pb.CreateAppRequest{
			VersionPackage: pbutil.ToProtoBytes(pkg),
			Name:           pbutil.ToProtoString(appName),
			VersionType:    pbutil.ToProtoString("helm"),
		}

		ctxFunc := func() (ctx context.Context) {
			ctx = context.Background()
			ctx = ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
			return
		}
		apps, err := client.DescribeActiveApps(ctxFunc(), &pb.DescribeAppsRequest{Name: []string{appName}})
		if err != nil {
			logger.Error(nil, "describe app error: %s", err.Error())
			os.Exit(-1)
		}
		if apps.TotalCount > 0 {
			continue
		}

		res, err := client.CreateApp(ctxFunc(), createReq)
		if err != nil {
			logger.Error(nil, "create app error: %s", err.Error())
			os.Exit(-1)
		}
		submitReq := &pb.SubmitAppVersionRequest{
			VersionId: res.VersionId,
		}

		_, err = client.SubmitAppVersion(ctxFunc(), submitReq)
		if err != nil {
			logger.Error(nil, "submit app error: %s", err.Error())
			os.Exit(-1)
		}

		passReq := &pb.PassAppVersionRequest{
			VersionId: res.VersionId,
		}
		_, err = client.AdminPassAppVersion(ctxFunc(), passReq)
		if err != nil {
			logger.Error(nil, "pass app error: %s", err.Error())
			os.Exit(-1)
		}

		releaseReq := &pb.ReleaseAppVersionRequest{
			VersionId: res.VersionId,
		}
		_, err = client.ReleaseAppVersion(ctxFunc(), releaseReq)
		if err != nil {
			logger.Error(nil, "release app error: %s", err.Error())
			os.Exit(-1)
		}
	}

}
