package common

func Contains(c interface{}, s interface{}) bool {
	cs, ok := s.([]interface{})
	if !ok {
		return false
	}
	for _, v := range cs {
		if v == c {
			return true
		}
	}
	return false
}
