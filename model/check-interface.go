package model

type CheckInterface interface {
	BuildContext(cx CheckContext)
	Validate() (bool, error)
	Check() (bool, bool, error)
}
