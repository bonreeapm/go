package bonree

import (
	"net/http"
	"github.com/bonreeapm/go/common"
)

// BusinessTransaction is the transaction of business.
type BusinessTransaction interface {
	http.ResponseWriter

	End()

	AddError(errorName string, summary string, details string, markBtAsError bool)

	AddException(exceptionName string, summary string, details string, markBtAsError bool)

	StartRPCExitCall(rpcType common.BR_RPC_TYPE, host string, port int) ExitCall

	StartSQLExitCall(sqlType common.BR_SQL_TYPE, host string, port int, dbschema string, vendor string, version string) ExitCall

	StartNoSQLExitCall(nosqlType common.BR_NOSQL_TYPE, serverPool string, port int, vendor string) ExitCall

	SnapshotFuncStart(className string, funcName string) SnapshotFunc

	SnapshotFuncEnd(snapshotFunc SnapshotFunc)
}

// GetCurrentTransaction get BusinessTransaction in http handler
func GetCurrentTransaction(w http.ResponseWriter) (BusinessTransaction) {
	if btn, ok := w.(BusinessTransaction); ok {
		return btn
	}

	return nil
}

func GetRoutineTransaction() BusinessTransaction {
	if _routineEngine != nil {
		routineValue := _routineEngine.Get()

		return (*btn)(routineValue)
	}

	return nil
}