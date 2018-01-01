package common

func EitherString(s, defaultS string) string {
	if s == "" {
		return defaultS
	}
	return s
}
