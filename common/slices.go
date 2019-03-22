package common

func ContainsString(c string, s []string) bool {
	for _, v := range s {
		if v == c {
			return true
		}
	}
	return false
}
