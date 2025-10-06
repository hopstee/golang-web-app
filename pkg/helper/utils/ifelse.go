package utils

func IfElse(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
