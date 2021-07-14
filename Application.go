package bonree

import (
	"net/http"
	"errors"
	"unsafe"
	"github.com/bonreeapm/go/sdk"
)

// Application represents the application.
type Application interface {
	StartBusinessTransaction(name string, w http.ResponseWriter, r *http.Request) BusinessTransaction

	Release()
}

type Engine interface {
	Get() unsafe.Pointer
	Set(p unsafe.Pointer)
}

var _routineEngine Engine = nil

// NewApplication is create a new Application
func NewApplication(appName string) (Application, error) {
	var handle sdk.AppHandle

	if appName == "" {
		handle = sdk.AppInit()
	}else {
		appConfig := &sdk.AppConfig{
			AppName: appName,
			AgentName: appName,
			TierName: appName,
			ClusterName: appName, 
		}
	
		handle = sdk.AppInitWithCfg(appConfig)
	}	

	if handle == 0 {
		return nil, errors.New("create new app fail")
	}

	app := NewApp(handle)

	return app, nil
}

// RoutineEngineInit is initialize RoutineEngine for sdk
func RoutineEngineInit(engine Engine) {
	_routineEngine = engine
}
