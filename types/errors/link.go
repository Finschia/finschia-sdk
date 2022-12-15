package errors

const (
	linkCodespace = "link"
)

var (
	// ErrError             = Register(linkCodespace, 1, "error")
	ErrInvalidPermission = Register(linkCodespace, 2, "invalid permission")
	ErrInvalidDenom      = Register(linkCodespace, 3, "invalid denom")
)
