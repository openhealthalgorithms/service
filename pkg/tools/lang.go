package tools

// TernaryString returns success value if condition is true, otherwise - it returns failure.
func TernaryString(condition bool, success, failure string) string {
	if !condition {
		return failure
	}

	return success
}

// TernaryInt returns success value if condition is true, otherwise - it returns failure.
func TernaryInt(condition bool, success, failure int) int {
	if !condition {
		return failure
	}

	return success
}

// TernaryInt32 returns success value if condition is true, otherwise - it returns failure.
func TernaryInt32(condition bool, success, failure int32) int32 {
	if !condition {
		return failure
	}

	return success
}

// TernaryInt64 returns success value if condition is true, otherwise - it returns failure.
func TernaryInt64(condition bool, success, failure int64) int64 {
	if !condition {
		return failure
	}

	return success
}

// TernaryBool returns success value if condition is true, otherwise - it returns failure.
func TernaryBool(condition bool, success, failure bool) bool {
	if !condition {
		return failure
	}

	return success
}
