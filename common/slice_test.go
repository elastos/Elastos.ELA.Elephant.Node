package common

import "testing"

func Test_Contains(t *testing.T) {
	s := []string{"11", "22"}
	println(ContainsString("11", s))
}
