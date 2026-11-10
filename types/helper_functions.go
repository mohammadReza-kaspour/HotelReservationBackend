package types

func derefString(strPointer *string) string {
	if strPointer == nil {
		return ""
	}
	return *strPointer
}
func derefNullable[T any](destObj map[string]any, destFieldName string, fieldPointer *T) {
	if fieldPointer != nil {
		destObj[destFieldName] = *fieldPointer
	}
}
