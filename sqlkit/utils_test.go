package sqlkit_test

import (
	"testing"

	"github.com/ilhammhdd/go-toolkit/sqlkit"
)

func TestGeneratePlaceHolder(t *testing.T) {
	testCases := []struct {
		n        uint16
		expected string
	}{
		{0, ""},
		{1, "(?)"},
		{2, "(?,?)"},
		{3, "(?,?,?)"},
		{10, "(?,?,?,?,?,?,?,?,?,?)"},
	}

	var result string
	for i := 0; i < len(testCases); i++ {
		result = string(sqlkit.GeneratePlaceHolder(testCases[i].n))
		if result != testCases[i].expected {
			t.Fatalf("n: %d expected: \"%s\" got: \"%s\"", testCases[i].n, testCases[i].expected, result)
		}
	}
}

func TestGenerateNPlaceHolder(t *testing.T) {
	testCases := []struct {
		nPlaceHolders, nParams uint16
		expected               string
	}{
		{0, 0, ""},
		{1, 1, "(?)"},
		{1, 2, "(?,?)"},
		{2, 2, "(?,?),(?,?)"},
		{2, 3, "(?,?,?),(?,?,?)"},
		{10, 2, "(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?)"}, //59
	}

	var result string
	for i := 0; i < len(testCases); i++ {
		result = string(sqlkit.GenerateNPlaceHolder(testCases[i].nPlaceHolders, testCases[i].nParams))
		if result != testCases[i].expected {
			t.Fatalf("nPlaceholders: %d nParams: %d expected: \"%s\" got: \"%s\"", testCases[i].nPlaceHolders, testCases[i].nParams, testCases[i].expected, result)
		}
	}
}
