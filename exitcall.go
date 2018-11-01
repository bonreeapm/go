package bonree

import (
	"net/http"
)

// ExitCall is the remote call of the application.
type ExitCall interface {
	AddError(errorName string, summary string, details string, markBtAsError bool)

	AddException(exceptionName string, summary string, details string, markBtAsError bool)

	SetDetail(cmd string, details string) error

	RoundTripper() http.RoundTripper

	SetCrossResheader(header http.Header)

	End()
}