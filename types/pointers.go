package types

// StringPtr returns a pointer to a string
func StringPtr(x string) *string {
	return &x
}

// Int64Ptr returns a pointer to an int64
func Int64Ptr(x int64) *int64 {
	return &x
}
