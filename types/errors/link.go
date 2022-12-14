package errors

const (
	linkCodespace = "link"
)

var (
	ErrError             = Register(linkCodespace, 2, "error")
	ErrInvalidPermission = Register(linkCodespace, 3, "invalid permission")
	ErrInvalidDenom      = Register(linkCodespace, 4, "invalid denom")
)
