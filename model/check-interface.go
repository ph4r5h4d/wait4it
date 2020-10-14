package model

type CheckInterface interface {
	BuildContext(cx CheckContext)
	Validate() error
	Check() (bool, bool, error)
}
