package dao

import (
	"os"
	"testing"
)

func TestCreateFunc(t *testing.T) {
	if len(os.Getenv("CLEARDB_DATABASE_URL")) == 0 {
		t.Error("Could not test")
		return
	}
	inputMappings := map[string]string{
		"family": "Input!B6",
		"child":  "Input!B7",
	}
	outputMappings := map[string]string{
		"subtotal": "Input!B10",
		"total":    "Input!B11",
	}
	if id, err := CreateFunc("91931784", "testFun", "testFile.xlsx",
		inputMappings, outputMappings); err != nil {
		t.Errorf("Could not create function. Reason: %v", err)
	} else if id < 0 {
		t.Errorf("Expected ID > 0. Actual: %d", id)
	}
}
