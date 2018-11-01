package sdk
//#cgo CFLAGS: -I/opt/bonree-agent-sdk/sdk_lib
//#cgo LDFLAGS: -L/opt/bonree-agent-sdk/sdk_lib/lib -lbonree_sdk -ldl -Wl,-rpath,/opt/bonree-agent-sdk/sdk_lib/lib
//#cgo linux LDFLAGS: -lrt
//#include "bonree.h"
//#include <stdlib.h>
//#include <stdint.h>
/*
uintptr_t apphandle_to_uint(br_app_t apphandle) {
    return (uintptr_t) apphandle;
}
br_app_t uint_to_apphandle(uintptr_t apphandle) {
    return (br_app_t) apphandle;
}

uintptr_t bthandle_to_uint(br_bt_t bthandle) {
    return (uintptr_t) bthandle;
}
br_bt_t uint_to_bthandle(uintptr_t bthandle) {
    return (br_bt_t) bthandle;
}

uintptr_t behandle_to_uint(br_backend_t behandle) {
    return (uintptr_t) behandle;
}
br_backend_t uint_to_behandle(uintptr_t behandle) {
    return (br_backend_t) behandle;
}

uintptr_t threadhandle_to_uint(br_snapshot_thread_t threadhandle) {
    return (uintptr_t) threadhandle;
}
br_snapshot_thread_t uint_to_threadhandle(uintptr_t threadhandle) {
    return (br_snapshot_thread_t) threadhandle;
}

uintptr_t funchandle_to_uint(br_snapshot_func_t funchandle) {
    return (uintptr_t) funchandle;
}
br_snapshot_func_t uint_to_funchandle(uintptr_t funchandle) {
    return (br_snapshot_func_t) funchandle;
}

uintptr_t exithandle_to_uint(br_exitcall_t exithandle) {
    return (uintptr_t) exithandle;
}
br_exitcall_t uint_to_exithandle(uintptr_t exithandle) {
    return (br_exitcall_t) exithandle;
}
*/
import "C"

import (
	"os"
    "os/signal"
    "sync"
	"syscall"
	"errors"
	"unsafe"
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
	result := int(C.br_sdk_init())

	if result == 0 {
		text := "Could not initialize the Golang SDK."
        return errors.New(text)
	}
	
	return nil
}

// Release Release the Bonree SDK.
func sdkRelease() {
	C.br_sdk_release()
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

func marshalAppConfig(from *AppConfig) C.br_app_config_t {
	to := C.br_app_config_t{}
	
	to.app_name = C.CString(from.AppName)
	to.tier_name = C.CString(from.TierName)
	to.cluster_name = C.CString(from.ClusterName)
	to.agent_name = C.CString(from.AgentName)

	return to
}

func freeAppConfigMembers(cAppConfig C.br_app_config_t) {
    if cAppConfig.app_name != nil {
        C.free(unsafe.Pointer(cAppConfig.app_name))
    }
    if cAppConfig.tier_name != nil {
        C.free(unsafe.Pointer(cAppConfig.tier_name))
    }
    if cAppConfig.cluster_name != nil {
        C.free(unsafe.Pointer(cAppConfig.cluster_name))
    }
    if cAppConfig.agent_name != nil {
        C.free(unsafe.Pointer(cAppConfig.agent_name))
    }
}

// AppInit is Init the app
func AppInit() AppHandle {
	appHandle := AppHandle(C.apphandle_to_uint(C.br_app_init()))
	appendAppHandle(appHandle)
	return appHandle
}

// AppInitWithCfg is Init the app with appConfig
func AppInitWithCfg(appConfig *AppConfig) AppHandle {
	cAppConfig := marshalAppConfig(appConfig)

	defer freeAppConfigMembers(cAppConfig)

	appHandle := AppHandle(C.apphandle_to_uint(C.br_app_init_with_cfg(&cAppConfig)))
	appendAppHandle(appHandle)
	return appHandle
}

// AppRelease Release the app
func AppRelease(appHandle AppHandle) {
	_appHandle := C.uint_to_apphandle(C.uintptr_t(appHandle))
	C.br_app_release(_appHandle)
	removeAppHandle(appHandle)
	appHandle = 0
}

// BtBegin Begin the BT
func BtBegin(appHandle AppHandle, name string) BtHandle {
	_appHandle := C.uint_to_apphandle(C.uintptr_t(appHandle))
	_name := C.CString(name)

	return BtHandle(C.bthandle_to_uint(C.br_bt_begin(_appHandle, _name)))
}

// BtBeginEx Begin the BT with crossRequestHeader
func BtBeginEx(appHandle AppHandle, name string, crossRequestHeader string) BtHandle {
	_appHandle := C.uint_to_apphandle(C.uintptr_t(appHandle))
	_name := C.CString(name)
	_crossRequestHeader := C.CString(crossRequestHeader)
	return BtHandle(C.bthandle_to_uint(C.br_bt_begin_ex(_appHandle, _name, _crossRequestHeader)))
}

// BtGenerateCrossResheader generate crossResponseHeader
func BtGenerateCrossResheader(btHandle BtHandle) string {
	_btHandle := C.uint_to_bthandle(C.uintptr_t(btHandle))
	return C.GoString(C.br_bt_generate_cross_resheader(_btHandle))
}

// BtEnd End the BT
func BtEnd(btHandle BtHandle) {
	_btHandle := C.uint_to_bthandle(C.uintptr_t(btHandle))

	C.br_bt_end(_btHandle)
}

// BtSetURL Set the URL to BT
func BtSetURL(btHandle BtHandle, url string) {
	_btHandle := C.uint_to_bthandle(C.uintptr_t(btHandle))
	_url := C.CString(url)

	C.br_bt_set_url(_btHandle, _url)
}

// BtAddError Add the Error to BT
func BtAddError(btHandle BtHandle, errorType common.BR_ERROR_TYPE, errorName string, summary string, details string, markBtAsError int) {
	_btHandle := C.uint_to_bthandle(C.uintptr_t(btHandle))
	_errorType := C.br_error_type(errorType)
	_errorName := C.CString(errorName)
	_summary := C.CString(summary)
	_details := C.CString(details)
	_markBtAsError := C.int(markBtAsError)

	C.br_bt_add_error(_btHandle, _errorType, _errorName, _summary, _details, _markBtAsError)
}

func BackendDeclareSQL(sqlType common.BR_SQL_TYPE, host string, port int, dbschema string, vendor string, version string) BackendHandle {
	C_ret := C.br_backend_declare_sql(C.br_sql_type(sqlType), C.CString(host), C.int(port), C.CString(dbschema), C.CString(vendor), C.CString(version))
       ret := BackendHandle(C.behandle_to_uint(C_ret))
       return ret
}

func BackendDeclareNosql(nosqlType common.BR_NOSQL_TYPE, serverPool string, port int, vendor string) BackendHandle {
	C_ret := C.br_backend_declare_nosql(C.br_nosql_type(nosqlType), C.CString(serverPool), C.int(port), C.CString(vendor))
       ret := BackendHandle(C.behandle_to_uint(C_ret))
       return ret
}

func BackendDeclareRPC(rpcType common.BR_RPC_TYPE, host string, port int) BackendHandle {
	C_ret := C.br_backend_declare_rpc(C.br_rpc_type(rpcType), C.CString(host), C.int(port))
       ret := BackendHandle(C.behandle_to_uint(C_ret))
       return ret
}

func ExitcallBegin(btHandle BtHandle, backend BackendHandle) ExitcallHandle {
	C_ret := C.br_exitcall_begin(C.uint_to_bthandle(C.uintptr_t(btHandle)), C.uint_to_behandle(C.uintptr_t(backend)))
       ret := ExitcallHandle(C.exithandle_to_uint(C_ret))
       return ret
}

func ExitcallSetDetail(exitcall ExitcallHandle, cmd string, details string) int {
	C_ret := C.br_exitcall_set_detail(C.uint_to_exithandle(C.uintptr_t(exitcall)), C.CString(cmd), C.CString(details))
       ret := int(C_ret)
       return ret
}

func ExitcallAddError(exitcall ExitcallHandle, errorType common.BR_ERROR_TYPE, errorName string, summary string, details string, markAsError int) {
	C.br_exitcall_add_error(C.uint_to_exithandle(C.uintptr_t(exitcall)), C.br_error_type(errorType), C.CString(errorName), C.CString(summary), C.CString(details), C.int(markAsError))
}

func ExitcallEnd(exitcall ExitcallHandle) {
	handle := C.uint_to_exithandle(C.uintptr_t(exitcall))
	C.br_exitcall_end(handle)
}

func ExitcallGenerateCrossReqheader(exitcall ExitcallHandle) string {
	_exitcallHandle := C.uint_to_exithandle(C.uintptr_t(exitcall))
	return C.GoString(C.br_exitcall_generate_cross_reqheader(_exitcallHandle))
}

func ExitcallSetCrossResheader(exitcall ExitcallHandle, crossResponseHeader string) {
	_exitcallHandle := C.uint_to_exithandle(C.uintptr_t(exitcall))
	_crossResponseHeader := C.CString(crossResponseHeader)
	C.br_exitcall_set_cross_resheader(_exitcallHandle, _crossResponseHeader)
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
