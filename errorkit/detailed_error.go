package errorkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

const (
	FlowErrHttpHeaderParamNotExists uint = iota
	FlowErrURLQueryNotExists
	DetailedErrLastIota
)

type ErrDescGenerator interface {
	GenerateDesc(uint, ...string) string
}

type ErrDescGeneratorFunc func(errorDescConst uint, args ...string) string

func (dgf ErrDescGeneratorFunc) GenerateDesc(errorDescConst uint, args ...string) string {
	return dgf(errorDescConst, args...)
}

type DetailedError struct {
	DateTime time.Time
	UUID     string
	// true if the error is caused by business flow or logic flow, false otherwise
	Flow         bool
	CallTrace    string
	WrappedErr   error
	ErrDescConst uint
	Desc         string
	logged       bool
}

func (de *DetailedError) Error() string {
	if de.DateTime.Format(time.RFC3339Nano) == "0001-01-01T00:00:00Z" {
		de.DateTime = time.Now().UTC()
	}
	switch {
	case !de.Flow && de.WrappedErr != nil:
		return fmt.Sprintf("date_time: %s uuid: %s flow: %t call_trace: %s desc: %s error: %s", de.DateTime, de.UUID, de.Flow, de.CallTrace, de.Desc, de.WrappedErr.Error())
	case !de.Flow && de.WrappedErr == nil:
		return fmt.Sprintf("date_time: %s uuid: %s flow: %t call_trace: %s desc: %s", de.DateTime, de.UUID, de.Flow, de.CallTrace, de.Desc)
	case de.Flow && de.WrappedErr != nil:
		return fmt.Sprintf("uuid: %s desc: %s error: %s", de.UUID, de.Desc, de.WrappedErr.Error())
	case de.Flow && de.WrappedErr == nil:
		return fmt.Sprintf("uuid: %s desc: %s", de.UUID, de.Desc)
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
	UUID         string `json:"uuid"`
	Desc         string `json:"desc"`
	ErrDescConst uint   `json:"err_desc_const,omitempty"`
}

type nonFlowStruct struct {
	DateTime     time.Time `json:"date_time"`
	UUID         string    `json:"uuid"`
	Flow         bool      `json:"flow"`
	CallTrace    string    `json:"call_trace"`
	WrappedErr   error     `json:"wrapped_err,omitempty"`
	ErrDescConst uint      `json:"err_desc_const,omitempty"`
	Desc         string    `json:"desc"`
	logged       bool      `json:"-"`
}

func (de DetailedError) MarshalJSON() ([]byte, error) {
	if de.Flow {
		return json.Marshal(flowStruct{de.UUID, de.Desc, de.ErrDescConst})
	}
	return json.Marshal(nonFlowStruct(de))
}

func (de *DetailedError) UnmarshalJSON(jsonData []byte) error {
	var result DetailedError
	return json.Unmarshal(jsonData, &result)
}

// NewDetailedError arg flow notating whether the cause is something from business flow or logic flow, or algorithmic one
func NewDetailedError(flow bool, callTrace string, wrappedErr error, errDescConst uint, descGenerator ErrDescGenerator, args ...string) *DetailedError {
	uuidRand, err := uuid.NewRandom()
	if err != nil {
		log.Println("error while generating random uuid")
		return nil
	}
	return &DetailedError{time.Now().UTC(), uuidRand.String(), flow, callTrace, wrappedErr, errDescConst, descGenerator.GenerateDesc(errDescConst, args...), false}
}

func IsNotNilThenLog(detailedErrs ...*DetailedError) bool {
	var notNilThenLog bool = false
	for i := range detailedErrs {
		if detailedErrs[i] != nil {
			if !detailedErrs[i].logged {
				log.Println(detailedErrs[i].Error())
				detailedErrs[i].logged = true
				notNilThenLog = true
			}
		}
	}
	return notNilThenLog
}
