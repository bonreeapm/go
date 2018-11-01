package sdk

import (
	"os"
    "os/signal"
    "sync"
	"syscall"
	"log"
	"bonree/common"
)

// AppHandle is the Handle of app object.
type AppHandle uint64

// BtHandle is the Handle of app object.
type BtHandle uint64

// BackendHandle is the Handle of app object.
type BackendHandle uint64

// SnapshotThreadHandle is the Handle of app object.
type SnapshotThreadHandle uint64

// SnapshotFuncHandle is the Handle of app object.
type SnapshotFuncHandle uint64

// ExitcallHandle is the Handle of app object.
type ExitcallHandle uint64

// AppConfig is including the info of app.
type AppConfig struct {
	AppName, TierName, ClusterName, AgentName string
}

var appHandleSlice []AppHandle

func init() {
	err := sdkInit()

	if err != nil {
		return
	}

	var stopLock sync.Mutex
    stop := false
    stopChan := make(chan struct{}, 1)
    signalChan := make(chan os.Signal, 1)
    go func() {
        <-signalChan
        stopLock.Lock()
        stop = true
        stopLock.Unlock()
		stopChan <- struct{}{}

		for i := 0; i < len(appHandleSlice); i++ {
			AppRelease(appHandleSlice[i])
		}
				
		sdkRelease()
        os.Exit(0)
    }()
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

// Init Initializes the Bonree SDK.
// Returns an error on failure.
func sdkInit() error {	
	log.Println("SDK Init().")
	return nil
}

// Release Release the Bonree SDK.
func sdkRelease() {
	log.Println("SDK Release().")
}

func appendAppHandle(appHandle AppHandle) {
	appHandleSlice = append(appHandleSlice, appHandle)
}

func removeAppHandle(appHandle AppHandle) {
	var _appHandleSlice []AppHandle
	for i := 0; i < len(appHandleSlice); i++ {
		if appHandleSlice[i] == appHandle {
			_appHandleSlice = append(appHandleSlice[:i], appHandleSlice[i+1:]...)			
		}
	}
	appHandleSlice = _appHandleSlice
}

// AppInit is Init the app
func AppInit() AppHandle {
	log.Println("SDK AppInit().")
	var appHandle AppHandle = 1
	appendAppHandle(appHandle)
	return appHandle
}

// AppInitWithCfg is Init the app with appConfig
func AppInitWithCfg(appConfig *AppConfig) AppHandle {
	var appHandle AppHandle = 1
	appendAppHandle(appHandle)
	log.Println("SDK AppInitWithCfg().")
	return appHandle
}

// AppRelease Release the app
func AppRelease(appHandle AppHandle) {
	if appHandle != 0 {
		removeAppHandle(appHandle)
		log.Println("SDK AppRelease().")
		appHandle = 0
	}
}

// BtBegin Begin the BT
func BtBegin(appHandle AppHandle, name string) BtHandle {
	log.Println("SDK BtBegin().")
	return 0
}

// BtBeginEx Begin the BT with crossRequestHeader
func BtBeginEx(appHandle AppHandle, name string, crossRequestHeader string) BtHandle {
	log.Println("SDK BtBeginEx().")
	return 0
}

// BtGenerateCrossResheader generate crossResponseHeader
func BtGenerateCrossResheader(btHandle BtHandle) string {
	log.Println("SDK BtGenerateCrossResheader().")
	return "nbg=1a2f69d6-4964-4eb2-43e3-6f0a744044d8;sst=2;sag=1eb82b05-c6e9-4ffd-bcb4-98757ee3b078;sbn=/ylexamples/fopen6.php;sbt=1514537122492;srh=;sbg=293e7c2d-4334-48c1-4c66-65626def4464"
}

// BtEnd End the BT
func BtEnd(btHandle BtHandle) {
	log.Println("SDK BtEnd().")
}

// BtSetURL Set the URL to BT
func BtSetURL(btHandle BtHandle, url string) {
	log.Println("SDK BtSetURL().")
}

// BtAddError Add the Error to BT
func BtAddError(btHandle BtHandle, errorType common.BR_ERROR_TYPE, errorName string, summary string, details string, markBtAsError int) {
	log.Println("SDK BtAddError().")
}

// BackendDeclareSQL declare the sql backend
func BackendDeclareSQL(sqlType common.BR_SQL_TYPE, host string, port int, dbschema string, vendor string, version string) BackendHandle {
	log.Println("SDK BackendDeclareSql().")
	return 0
}

// BackendDeclareNosql declare the nosql backend
func BackendDeclareNosql(nosqlType common.BR_NOSQL_TYPE, serverPool string, port int, vendor string) BackendHandle {
	log.Println("SDK BackendDeclareNosql().")
	return 0
}

// BackendDeclareRPC the rpc backend
func BackendDeclareRPC(rpcType common.BR_RPC_TYPE, host string, port int) BackendHandle {
	log.Println("SDK BackendDeclareRpc().")
	return 0
}

// ExitcallBegin begin the exitcall
func ExitcallBegin(btHandle BtHandle, backend BackendHandle) ExitcallHandle {
	log.Println("SDK ExitcallBegin().")
	return 0
}

// ExitcallSetDetail set the detail of exitcall
func ExitcallSetDetail(exitcall ExitcallHandle, cmd string, details string) int {
	log.Println("SDK ExitcallSetDetail().")
	return 0
}

// ExitcallAddError add the error to exitcall
func ExitcallAddError(exitcall ExitcallHandle, errorType common.BR_ERROR_TYPE, errorName string, summary string, details string, markAsError int) {
	log.Println("SDK ExitcallAddError().")
}

// ExitcallEnd end the exitcall
func ExitcallEnd(exitcall ExitcallHandle) {
	log.Println("SDK ExitcallEnd().")
}

// ExitcallGenerateCrossReqheader generate the crossrequestheader string
func ExitcallGenerateCrossReqheader(exitcall ExitcallHandle) string {
	log.Println("SDK ExitcallGenerateCrossReqheader().")
	return "nbg=1a2f69d6-4964-4eb2-43e3-6f0a744044d8;sst=2;sag=1eb82b05-c6e9-4ffd-bcb4-98757ee3b078;sbn=/ylexamples/fopen6.php;sbt=1514537122492;srh=;sbg=293e7c2d-4334-48c1-4c66-65626def4464"
}

// ExitcallSetCrossResheader set the crossresponseheader to exitcall
func ExitcallSetCrossResheader(exitcall ExitcallHandle, crossResponseHeader string) {
	log.Println("SDK ExitcallSetCrossResheader().")
}

// func BtIsSnapshotting(btHandle BtHandle) byte {
// 	_btHandle := C.uint_to_bthandle(C.uintptr_t(btHandle))
// 	return byte(C.br_bt_is_snapshotting(_btHandle))
// }

// func BtSnapshotData(btHandle BtHandle, key string, value string) {
// 	_btHandle := C.uint_to_bthandle(C.uintptr_t(btHandle))
// 	_key := C.CString(key)
// 	_value := C.CString(value)
// 	C.br_bt_snapshot_data(_btHandle, _key, _value)
// }

// func BtSnapshotThreadStart(btHandle BtHandle) SnapshotThreadHandle {
// 	_btHandle := C.uint_to_bthandle(C.uintptr_t(btHandle))
// 	return SnapshotThreadHandle(C.threadhandle_to_uint(C.br_bt_snapshot_thread_start(_btHandle)))
// }

// func BtSnapshotThreadEnd(snapshotThreadHandle SnapshotThreadHandle) {
// 	_snapshotThreadHandle := C.uint_to_threadhandle(C.uintptr_t(snapshotThreadHandle))
// 	C.br_bt_snapshot_thread_end(_snapshotThreadHandle)
// }

// func BtSnapshotFuncStart(snapshotThreadHandle SnapshotThreadHandle, className string, funcName string, fileName string, lineno int) SnapshotFuncHandle {
// 	_snapshotThreadHandle := C.uint_to_threadhandle(C.uintptr_t(snapshotThreadHandle))
// 	_className := C.CString(className)
// 	_funcName := C.CString(funcName)
// 	_fileName := C.CString(fileName)
// 	_lineno := C.int(lineno)
// 	return SnapshotFuncHandle(C.funchandle_to_uint(C.br_bt_snapshot_func_start(_snapshotThreadHandle, _className, _funcName, _fileName, _lineno)))
// }

// func BtSnapshotFuncEnd(snapshotFuncHandle SnapshotFuncHandle) {
// 	_snapshotFuncHandle := C.uint_to_funchandle(C.uintptr_t(snapshotFuncHandle))
// 	C.br_bt_snapshot_func_end(_snapshotFuncHandle)
// }

// func SnapshotExitcallAdd(snapshotFuncHandle SnapshotFuncHandle, exitcall ExitcallHandle) {
// 	_snapshotFuncHandle := C.uint_to_funchandle(C.uintptr_t(snapshotFuncHandle))
// 	_exitcallHandle := C.uint_to_exithandle(C.uintptr_t(exitcall))
// 	C.br_snapshot_exitcall_add(_snapshotFuncHandle, _exitcallHandle)
// }

// func SnapshotErrorAdd(snapshotThreadHandle SnapshotThreadHandle, errorName string, summary string, details string) {
// 	_snapshotThreadHandle := C.uint_to_threadhandle(C.uintptr_t(snapshotThreadHandle))
// 	_errorName := C.CString(errorName)
// 	_summary := C.CString(summary)
// 	_details := C.CString(details)
// 	C.br_snapshot_error_add(_snapshotThreadHandle, _errorName, _summary, _details)
// }
