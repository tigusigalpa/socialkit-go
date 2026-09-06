package socialkit

// StringPtr returns a pointer to the given string. This is a convenience
// helper for creating optional string fields in request structs.
func StringPtr(v string) *string {
	return &v
}

// IntPtr returns a pointer to the given int.
func IntPtr(v int) *int {
	return &v
}

// BoolPtr returns a pointer to the given bool.
func BoolPtr(v bool) *bool {
	return &v
}

// Float64Ptr returns a pointer to the given float64.
func Float64Ptr(v float64) *float64 {
	return &v
}
