package exec

import (
	"github.com/yuin/gopher-lua"
	"testing"
)

func testExp(t *testing.T, line string, expect lua.LValue) {
	L := NewExecutor()
	defer L.Close()
	// Execute
	if err := L.DoString("out = " + line); err != nil {
		t.Fatal(err)
	}
	// Check
	if L.GetGlobal("out") != expect {
		t.Fatalf("Expect %v, receive %v", expect, L.GetGlobal("out"))
	}
}

func TestLoad(t *testing.T) {
	testExp(t, "len(load(\"ml-100k\"))", lua.LNumber(100000))
}
