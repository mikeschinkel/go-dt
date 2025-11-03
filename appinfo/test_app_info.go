package appinfo

import (
	"github.com/mikeschinkel/go-dt"
)

type Test struct {
	*appInfo
}

func NewTest(ai AppInfo) Test {
	return Test{appInfo: ai.(*appInfo)}
}

func (ai *Test) SetExtraInfo(info map[string]any) {
	ai.extraInfo = info
}

//	func (ai *Test) SetConfigPath(cd dt.PathSegments) {
//		ai.configPath = cd
//	}
func (ai *Test) SetConfigFile(cf dt.RelFilepath) {
	ai.configFile = cf
}
func (ai *Test) SetLogFile(lf dt.Filename) {
	ai.logFile = lf
}
