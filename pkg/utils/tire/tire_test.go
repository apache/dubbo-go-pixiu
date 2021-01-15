package tire

import "testing"

func TestName(t *testing.T) {
	keys := []string{"aa", "bb", "cc"}
	keys2 := keys[1:]
	keys3 := keys2[1:]
	keys4 := keys3[1:]
	keys5 := keys4[1:]
	println(keys5)
}
