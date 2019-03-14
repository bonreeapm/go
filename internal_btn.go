package bonree

import (
	"net/http"
	"github.com/bonreeapm/go/common"
	"github.com/bonreeapm/go/sdk"
	"runtime"
)

type btn struct {
	W http.ResponseWriter
	btHandle sdk.BtHandle
	snapshotThreadHandle sdk.SnapshotThreadHandle
}

func (btn *btn) End() {
	sdk.BtSnapshotThreadEnd(btn.snapshotThreadHandle)
	sdk.BtEnd(btn.btHandle)
}

func (btn *btn) SetURL(url string) {
	sdk.BtSetURL(btn.btHandle, url)
}

func (btn *btn) AddError(errorName string, summary string, details string, markBtAsError bool) {
	_markBtAsError := 0

	if markBtAsError {
		_markBtAsError = 1
	}

	sdk.BtAddError(btn.btHandle, common.BR_ERROR_TYPE_HTTP, errorName, summary, details, _markBtAsError)
	sdk.SnapshotErrorAdd(btn.snapshotThreadHandle, errorName, summary, details)
}

func (btn *btn) AddException(exceptionName string, summary string, details string, markBtAsError bool) {
	_markBtAsError := 0

	if markBtAsError {
		_markBtAsError = 1
	}

	sdk.BtAddError(btn.btHandle, common.BR_ERROR_TYPE_EXCEPTION, exceptionName, summary, details, _markBtAsError)
}

func (btn *btn) StartRPCExitCall(rpcType common.BR_RPC_TYPE, host string, port int) ExitCall {
	backendHandle := sdk.BackendDeclareRPC(rpcType, host, port)
	exitcallHandle := sdk.ExitcallBegin(btn.btHandle, backendHandle)

	return &exitcall{
		exitcallHandle: exitcallHandle,
	}
}

func (btn *btn) StartSQLExitCall(sqlType common.BR_SQL_TYPE, host string, port int, dbschema string, vendor string, version string) ExitCall {
	backendHandle := sdk.BackendDeclareSQL(sqlType, host, port, dbschema, vendor, version)
	exitcallHandle := sdk.ExitcallBegin(btn.btHandle, backendHandle)

	return &exitcall{
		exitcallHandle: exitcallHandle,
	}
}

func (btn *btn) StartNoSQLExitCall(nosqlType common.BR_NOSQL_TYPE, serverPool string, port int, vendor string) ExitCall {
	backendHandle := sdk.BackendDeclareNosql(nosqlType, serverPool, port, vendor)
	exitcallHandle := sdk.ExitcallBegin(btn.btHandle, backendHandle)

	return &exitcall{
		exitcallHandle: exitcallHandle,
	}
}

func (btn *btn) Header() http.Header { 
	return btn.W.Header() 
}

func (btn *btn) Write(b []byte) (int, error) {
	return btn.W.Write(b)
}

func (btn *btn) WriteHeader(code int) {
	btn.W.WriteHeader(code)
}

func newBtn(app *app, name string, w http.ResponseWriter, r *http.Request) BusinessTransaction {
	var handle sdk.BtHandle

	crossReqHeader :=  r.Header.Get(common.CrossRequestHeader)
	if crossReqHeader == "" {
		handle = sdk.BtBegin(app.appHandle, name)
	} else {		
		handle = sdk.BtBeginEx(app.appHandle, name, crossReqHeader)

		crossResHeader := sdk.BtGenerateCrossResheader(handle)

		w.Header().Add(common.CrossResponseHeader, crossResHeader)
	}

	sdk.BtSetURL(handle, r.URL.RequestURI())

	snapshotThreadHandle := sdk.BtSnapshotThreadStart(handle)

	return &btn{
		W: w,
		btHandle: handle,
		snapshotThreadHandle: snapshotThreadHandle,
	}
}

func (btn *btn) SnapshotFuncStart(className string, funcName string) SnapshotFunc {
	_,file,line,ok := runtime.Caller(1)
	if ok {
		_snapshotFuncHandle := sdk.BtSnapshotFuncStart(btn.snapshotThreadHandle, className, funcName, file, line)

		return &snapshotFunc{
			snapshotFuncHandle: _snapshotFuncHandle,
		}
	}

	return nil
}

func (btn *btn) SnapshotFuncEnd(snapshotFunc SnapshotFunc) {
	snapshotFunc.End()
}