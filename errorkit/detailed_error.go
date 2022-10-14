package errorkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type ErrDescConst uint

type DescGenerator interface {
	GenerateDesc(ErrDescConst, ...string) string
}

type DescGeneratorFunc func(edc ErrDescConst, args ...string) string

func (dgf DescGeneratorFunc) GenerateDesc(edc ErrDescConst, args ...string) string {
	return dgf(edc, args...)
}

type DetailedError struct {
	DateTime time.Time
	// Notating whether the cause is something from business flow or logic flow, or algorithmic one
	Flow         bool
	CallTrace    string
	WrappedErr   error
	ErrDescConst ErrDescConst
	Desc         string
	logged       bool
}

func (de *DetailedError) Error() string {
	if de.DateTime.Format(time.RFC3339Nano) == "0001-01-01T00:00:00Z" {
		de.DateTime = time.Now().UTC()
	}
	switch {
	case de.WrappedErr != nil && !de.Flow:
		return fmt.Sprintf("date_time: %s flow: %t call_trace: %s desc: %s error: %s", de.DateTime, de.Flow, de.CallTrace, de.Desc, de.WrappedErr.Error())
	case de.WrappedErr == nil && !de.Flow:
		return fmt.Sprintf("date_time: %s flow: %t call_trace: %s desc: %s", de.DateTime, de.Flow, de.CallTrace, de.Desc)
	case de.Flow:
		return de.Desc
	default:
		return ""
	}
}

func (de *DetailedError) Unwrap() error { return de.WrappedErr }

func (de *DetailedError) Is(target error) bool {
	var targetAsDetailedError *DetailedError
	if !errors.As(target, &targetAsDetailedError) {
		return false
	}
	return de.Flow == targetAsDetailedError.Flow && de.CallTrace == targetAsDetailedError.CallTrace && de.Desc == targetAsDetailedError.Desc
}

func (de *DetailedError) IsWrappedErrNotNilThenLog() bool {
	if de.WrappedErr != nil {
		log.Printf("\n%s", de.WrappedErr.Error())
		return true
	}
	return false
}

type flowStruct struct {
	Desc string `json:"desc"`
}

type nonFlowStruct struct {
	DateTime     time.Time    `json:"date_time"`
	Flow         bool         `json:"flow"`
	CallTrace    string       `json:"call_trace"`
	WrappedErr   error        `json:"wrapped_err,omitempty"`
	ErrDescConst ErrDescConst `json:"err_desc_const,omitempty"`
	Desc         string       `json:"desc"`
	logged       bool         `json:"-"`
}

func (de DetailedError) MarshalJSON() ([]byte, error) {
	if de.Flow {
		return json.Marshal(flowStruct{de.Desc})
	}
	return json.Marshal(nonFlowStruct(de))
}

func (de *DetailedError) UnmarshalJSON(jsonData []byte) error {
	var result DetailedError
	return json.Unmarshal(jsonData, &result)
}

// NewDetailedError arg flow notating whether the cause is something from business flow or logic flow, or algorithmic one
func NewDetailedError(flow bool, callTrace string, wrappedErr error, errDescConst ErrDescConst, descGenerator DescGenerator, args ...string) *DetailedError {
	return &DetailedError{time.Now().UTC(), flow, callTrace, wrappedErr, errDescConst, descGenerator.GenerateDesc(errDescConst, args...), false}
}

func IsNotNilThenLog(err error) (*DetailedError, bool) {
	if err != nil {
		if detailedErr, ok := err.(*DetailedError); ok {
			if !detailedErr.logged {
				log.Printf("\n%s", detailedErr.Error())
				detailedErr.logged = true
			}
			return detailedErr, true
		}
		log.Printf("\ndate_time: %s error: %s", time.Now().UTC(), err.Error())
		return nil, true
	}
	return nil, false
}
