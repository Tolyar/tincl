package tincl

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func WriteToTelnet(t *Telnet) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		cmd := L.ToString(1) /* get argument */
		if n, err := t.WriteLine(cmd); err != nil || n < len(cmd) {
			return 0
		}

		return 1
	}
}

func ReadFromTelnet(t *Telnet) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		s, err := t.ReadLine()
		L.Push(lua.LString(s)) /* push result */
		if err != nil {
			return 0
		}

		return 1
	}
}

func RunScript(path string, t *Telnet) {
	L := lua.NewState()
	defer L.Close()

	// Register function for sending to telnet.
	L.SetGlobal("WriteToTelnet", L.NewFunction(WriteToTelnet(t)))
	L.SetGlobal("ReadFromTelnet", L.NewFunction(ReadFromTelnet(t)))

	if err := L.DoFile(path); err != nil {
		log.Fatalf("Can't load lua script '%s' : %v\n", path, err)
	}
}
