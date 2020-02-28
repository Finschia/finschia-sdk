package types

type Change struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func NewChange(field string, value string) Change {
	return Change{
		Field: field,
		Value: value,
	}
}

type Changes []Change

func NewChanges(changes ...Change) Changes {
	return changes
}

func NewChangesWithMap(changesMap map[string]string) Changes {
	changes := make([]Change, len(changesMap))
	idx := 0
	for k, v := range changesMap {
		changes[idx] = Change{Field: k, Value: v}
		idx++
	}
	return NewChanges(changes...)
}
