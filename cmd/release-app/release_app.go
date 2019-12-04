package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"io/ioutil"
	"openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type App struct {
	Appname  string    `json:"appname"`
	Category string    `json:"category"`
	Versions []Version `json:"versions"`
}

type Version struct {
	Pkgname string `json:"pkgname"`
	Url     string `json:"url"`
}

type Apps []App

func ReadConfig(path string) *Apps {
	var apps *Apps
	apps = new(Apps)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error(nil, "read config.json error: %s", err.Error())
	}
	err = json.Unmarshal(content, apps)
	if err != nil {
		logger.Error(nil, "unmarshal error: %s", err.Error())
	}
	return apps
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

func main() {
	var path string
	flag.StringVar(&path, "path", "../package/", "need package path.eg./your/path/to/pkg/")
	flag.Parse()
	items := ReadConfig(path + "/config.json")
	basePath := path
	client, err := app.NewAppManagerClient()
	if err != nil {
		panic(err)
	}
	ctxFunc := func() (ctx context.Context) {
		ctx = context.Background()
		ctx = ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
		return
	}

	for _, item := range *items {
		filePath := basePath + item.Appname + "/" + item.Versions[0].Pkgname
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			logger.Error(nil, "read package [%s] error: %s", filePath, err.Error())
		}
		pkg := strfmt.Base64(content)
		createReq := &pb.CreateAppRequest{
			VersionPackage: pbutil.ToProtoBytes(pkg),
			Name:           pbutil.ToProtoString(item.Appname),
			VersionType:    pbutil.ToProtoString("helm"),
		}
		resp, err := client.CreateApp(ctxFunc(), createReq)
		if err != nil {
			logger.Error(nil, "create app error: %s", err.Error())
		}
		AuditApp(ctxFunc, client, resp.VersionId)

		//Create AppVersion
		for _, version := range item.Versions[1:len(item.Versions)] {
			filePath := basePath + item.Appname + "/" + version.Pkgname
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				logger.Error(nil, "read package [%s] error: %s", filePath, err.Error())
			}
			pkg := strfmt.Base64(content)
			createVersion := &pb.CreateAppVersionRequest{
				AppId:       resp.AppId,
				Name:        pbutil.ToProtoString(""),
				Description: pbutil.ToProtoString(""),
				Type:        pbutil.ToProtoString("helm"),
				Package:     pbutil.ToProtoBytes(pkg),
			}
			res, err := client.CreateAppVersion(ctxFunc(), createVersion)
			if err != nil {
				logger.Error(nil, "create AppVersion error: %s", err.Error())
			}
			AuditApp(ctxFunc, client, res.VersionId)
		}
	}

}
