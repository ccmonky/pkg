package utils

// Abs for int64
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

// Divmod like python divmod
func Divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
