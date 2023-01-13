package errors

const (
	linkCodespace = "link"
)

var (
	ErrInvalidPermission = Register(linkCodespace, 2, "invalid permission")
	ErrInvalidDenom      = Register(linkCodespace, 3, "invalid denom")
)
