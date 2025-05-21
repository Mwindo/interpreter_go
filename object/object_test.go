package object

import "testing"

func TestStringHashKey(t *testing.T) {
	sameString1 := &String{Value: "Hello World"}
	sameString2 := &String{Value: "Hello World"}
	diffString1 := &String{Value: "hello World"}

	if sameString1.HashKey() != sameString2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if sameString1.HashKey() == diffString1.HashKey() {
		t.Errorf("strings with different content have different hash keys")
	}
}
