package errors

func RegisteredErrors() []*Error {
	es := make([]*Error, 0, len(usedCodes))
	for _, e := range usedCodes {
		es = append(es, e)
	}
	return es
}
