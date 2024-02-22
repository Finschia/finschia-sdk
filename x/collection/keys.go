package collection

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "collection"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

func CollectionAttrCanonicalKey(key string) string {
	convert := map[string]string{
		AttributeKeyBaseImgURI.String(): AttributeKeyURI.String(),
	}
	if converted, ok := convert[key]; ok {
		return converted
	}
	return key
}
