package bonree

import (
	"github.com/bonreeapm/go/sdk"
)

type snapshotFunc struct {
	snapshotFuncHandle sdk.SnapshotFuncHandle
}

func (snapshotFunc *snapshotFunc) End() {
	sdk.BtSnapshotFuncEnd(snapshotFunc.snapshotFuncHandle)
}

func (snapshotFunc *snapshotFunc) AddExitCall(exitCall ExitCall) {
	sdk.SnapshotExitcallAdd(snapshotFunc.snapshotFuncHandle, exitCall.Handle())
}