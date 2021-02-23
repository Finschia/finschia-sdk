package errors

const (
	CodespaceLink = "link"
)

var (
	ErrError             = Register(CodespaceLink, 1, "error")
	ErrInvalidPermission = Register(CodespaceLink, 2, "invalid permission")
	ErrInvalidDenom      = Register(CodespaceLink, 3, "invalid denom")
)
