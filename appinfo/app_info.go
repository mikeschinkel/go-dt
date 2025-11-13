package appinfo

import (
	"github.com/mikeschinkel/go-dt"
)

// TODO: Need to add comments to explain what each of these do

type AppInfo interface {
	AppName() string
	AppDescr() string
	AppVer() dt.Version
	AppSlug() dt.PathSegment
	ConfigSlug() dt.PathSegment
	ConfigFile() dt.RelFilepath
	ExeName() dt.Filename
	LogFile() dt.Filename
	LogPath() dt.PathSegments
	InfoURL() dt.URL
	ExtraInfo() map[string]any
}

var _ AppInfo = (*appInfo)(nil)

type appInfo struct {
	appName    string
	appDescr   string
	appVer     dt.Version
	appSlug    dt.PathSegment
	configSlug dt.PathSegment
	configFile dt.RelFilepath
	exeName    dt.Filename
	logFile    dt.Filename
	logPath    dt.PathSegments
	infoURL    dt.URL
	extraInfo  map[string]any
}

type Args struct {
	AppName    string
	AppDescr   string
	AppVer     dt.Version
	AppSlug    dt.PathSegment
	ConfigSlug dt.PathSegment
	ConfigFile dt.RelFilepath
	ExeName    dt.Filename
	LogFile    dt.Filename
	LogPath    dt.PathSegments
	InfoURL    dt.URL
	ExtraInfo  map[string]any
}

func New(args Args) AppInfo {
	if args.ExtraInfo == nil {
		args.ExtraInfo = make(map[string]any)
	}
	return &appInfo{
		appName:    args.AppName,
		appDescr:   args.AppDescr,
		appVer:     args.AppVer,
		appSlug:    args.AppSlug,
		configSlug: args.ConfigSlug,
		configFile: args.ConfigFile,
		exeName:    args.ExeName,
		logFile:    args.LogFile,
		logPath:    args.LogPath,
		infoURL:    args.InfoURL,
		extraInfo:  args.ExtraInfo,
	}
}

func (ai *appInfo) AppName() string {
	return ai.appName
}
func (ai *appInfo) AppDescr() string {
	return ai.appDescr
}
func (ai *appInfo) AppVer() dt.Version {
	return ai.appVer
}
func (ai *appInfo) InfoURL() dt.URL {
	return ai.infoURL
}
func (ai *appInfo) ExtraInfo() map[string]any {
	return ai.extraInfo
}
func (ai *appInfo) AppSlug() dt.PathSegment {
	return ai.appSlug
}
func (ai *appInfo) ConfigSlug() dt.PathSegment {
	return ai.configSlug
}
func (ai *appInfo) ConfigFile() dt.RelFilepath {
	return ai.configFile
}
func (ai *appInfo) ExeName() dt.Filename {
	return ai.exeName
}
func (ai *appInfo) LogFile() dt.Filename {
	return ai.logFile
}
func (ai *appInfo) LogPath() dt.PathSegments {
	return ai.logPath
}
