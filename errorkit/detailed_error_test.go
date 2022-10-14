package errorkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

const callTraceFileErrorTest = "/errorkit/error_test.go"

const (
	Err1stLayerNotFound ErrDescConst = iota
	Err1stLayerInvalidType
	LastErr1stLayer
)

const (
	Err2ndLayerDriverBusy ErrDescConst = iota + LastErr1stLayer
	LastErr2ndLayer
)

const (
	Err3rdLayerClientTimeout ErrDescConst = iota + LastErr2ndLayer
)

type DescGeneration string

func (dg DescGeneration) GenerateDesc(edc ErrDescConst, args ...string) string {
	switch edc {
	case Err1stLayerNotFound:
		return fmt.Sprintf("%s not found", dg)
	case Err1stLayerInvalidType:
		return fmt.Sprintf("%s invalid type", dg)
	case Err2ndLayerDriverBusy:
		return fmt.Sprintf("%s driver busy", dg)
	case Err3rdLayerClientTimeout:
		return fmt.Sprintf("%s client timeout", dg)
	default:
		return "error description constant undefined"
	}
}

type Response struct {
	Data   string  `json:"data,omitempty"`
	Errors []error `json:"errors"`
}

func TestDetailedError(t *testing.T) {
	var callTraceFunc = fmt.Sprintf("%s#TestDetailedError", callTraceFileErrorTest)
	internal1 := NewDetailedError(true, callTraceFunc, nil, Err1stLayerNotFound, DescGeneration("internal1"))
	response1 := Response{Data: "this is the data1", Errors: []error{internal1}}
	response1Json, err := json.Marshal(response1)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("%s", response1Json)

	wrappedErr := errors.New("wrapped err inernal3")
	internal2 := NewDetailedError(true, callTraceFunc, nil, Err1stLayerNotFound, DescGeneration("internal2"))
	internal3 := NewDetailedError(true, callTraceFunc, wrappedErr, Err1stLayerNotFound, DescGeneration("internal3"))
	t.Logf("errors.Is(internal2, internal1): %t", errors.Is(internal2, internal1))
	t.Logf("errors.Is(internal3, internal1): %t", errors.Is(internal3, internal1))
	var internal4 *DetailedError
	t.Logf("errors.As(internal3, &internal4): %t", errors.As(internal3, &internal4))
	response2 := Response{Data: "this is the data2", Errors: []error{internal4}}
	response2Json, err := json.Marshal(response2)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("%s", response2Json)

	internal5 := NewDetailedError(true, callTraceFunc, wrappedErr, Err1stLayerNotFound, DescGeneration("internal3"))
	t.Logf("errors.Is(internal4, internal5): %t", errors.Is(internal4, internal5))
	response3 := Response{Data: "this is the data3", Errors: []error{internal5}}
	response3Json, err := json.Marshal(response3)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("%s", response3Json)

	external1 := internal5
	external1.Flow = false
	response4 := Response{Data: "this is the data4", Errors: []error{external1}}
	response4Json, err := json.Marshal(response4)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("%s", response4Json)

	detailedErr1 := NewDetailedError(false, callTraceFunc, nil, Err1stLayerInvalidType, DescGeneration("external1"))
	IsNotNilThenLog(detailedErr1)

	detailedErr2 := NewDetailedError(true, callTraceFunc, nil, Err1stLayerInvalidType, DescGeneration("detailed internal error 2"))
	IsNotNilThenLog(detailedErr2)
}
