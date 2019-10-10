package narwhal_lib

import "testing"

func TestNarwhal_Save(t *testing.T) {
	n := Narwhal{false}
	n.Save("cyanprint","data","./")
}

func TestNarwhal_Load(t *testing.T) {
	n :=Narwhal{false}
	n.Load("ezvol", "./data.tar.gz")
}
