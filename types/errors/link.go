package errors

const (
	linkCodespace = "link"
)

// NO additional errors allowed into this codespace
var (
	ErrInvalidPermission = Register(linkCodespace, 2, "invalid permission")
	ErrInvalidDenom      = Register(linkCodespace, 3, "invalid denom")
)
