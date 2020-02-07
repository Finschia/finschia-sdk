package types

const (
	ModuleName = "safetybox"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	SafetyBoxKeyPrefix = []byte{0x04}
)

func SafetyBoxKey(safetyBoxID string) []byte {
	return append(SafetyBoxKeyPrefix, []byte(safetyBoxID)...)
}
