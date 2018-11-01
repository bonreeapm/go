package bonree

import (
	"errors"
	"bonree/common"
	"bonree/sdk"
	"net/http"
)

type exitcall struct {
	exitcallHandle sdk.ExitcallHandle
}

func (exitcall *exitcall) AddError(errorName string, summary string, details string, markBtAsError bool) {
	_markBtAsError := 0

	if markBtAsError {
		_markBtAsError = 1
	}

	sdk.ExitcallAddError(exitcall.exitcallHandle, common.BR_ERROR_TYPE_HTTP, errorName, summary, details, _markBtAsError)
}

func (exitcall *exitcall) AddException(exceptionName string, summary string, details string, markBtAsError bool) {
	_markBtAsError := 0

	if markBtAsError {
		_markBtAsError = 1
	}

	sdk.ExitcallAddError(exitcall.exitcallHandle, common.BR_ERROR_TYPE_EXCEPTION, exceptionName, summary, details, _markBtAsError)
}

func (exitcall *exitcall) SetDetail(cmd string, details string) error {
	ret := sdk.ExitcallSetDetail(exitcall.exitcallHandle, cmd, details)
	if ret != 0 {
		return errors.New("ExitCall setDetail is fail")
	}

	return nil
}

func (exitcall *exitcall) RoundTripper() http.RoundTripper {
	crossReqheader := sdk.ExitcallGenerateCrossReqheader(exitcall.exitcallHandle)

	return bonreeRoundTripper(crossReqheader)
}

func (exitcall *exitcall) SetCrossResheader(header http.Header) {
	if header == nil {
		return
	}

	crossResponseHeader := header.Get(common.CrossResponseHeader)

	if crossResponseHeader == "" {
		return
	}

	sdk.ExitcallSetCrossResheader(exitcall.exitcallHandle, crossResponseHeader)
}

func (exitcall *exitcall) End() {
	sdk.ExitcallEnd(exitcall.exitcallHandle)
}