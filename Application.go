package bonree

import (
	"net/http"
	"errors"
	"bonree/sdk"
)

// Application represents the application.
type Application interface {
	StartBusinessTransaction(name string, w http.ResponseWriter, r *http.Request) BusinessTransaction

	Release()
}

func NewApplication(appName string) (Application, error) {
	var handle sdk.AppHandle

	if appName == "" {
		handle = sdk.AppInit()
	}else {
		appConfig := &sdk.AppConfig{
			AppName: appName,
			AgentName: appName,
			TierName: "My Tier",
			ClusterName: "My Cluster", 
		}
	
		handle = sdk.AppInitWithCfg(appConfig)
	}	

	if handle == 0 {
		return nil, errors.New("create new app fail")
	}

	app := NewApp(handle)

	return app, nil
}
