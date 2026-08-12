package engine

import "testing"

func TestEngineStartup(t *testing.T) {
	engine, err := Startup()
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnginePut(t *testing.T) {

}

func TestEngineGet(t *testing.T) {

}

func TestEngineDelete(t *testing.T) {

}
