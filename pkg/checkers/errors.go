package checkers

type temporaryError string

func (e temporaryError) Error() string   { return string(e) }
func (e temporaryError) Temporary() bool { return true }
