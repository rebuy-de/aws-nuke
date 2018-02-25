package resources

func UnPtrBool(ptr *bool, def bool) bool {
	if ptr == nil {
		return def
	}
	return *ptr
}

func EqualStringPtr(v1, v2 *string) bool {
	if v1 == nil && v2 == nil {
		return true
	}

	if v1 == nil || v2 == nil {
		return false
	}

	return *v1 == *v2
}
