package bonree

// SnapshotFunc is the function of business.
type SnapshotFunc interface {
	End()

	AddExitCall(exitCall ExitCall)
}