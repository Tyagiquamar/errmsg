package errmsg_test

import (
	"encoding/json"
	"errors"
	"github.com/sardine-ai/errmsg"
	errmsg_testing "github.com/sardine-ai/errmsg/testing"
	"testing"
)

type birthday struct {
	Month int8 `json:"month"`
	Day   int8 `json:"day"`
}

type count struct {
	Value uint8 `json:"value"`
}

type user struct {
	ID string `json:"id"`
	BirthDay birthday `json:"birthday"`
}

func TestJsonTypeError(t *testing.T) {
	var u user
	input := `{"id": 1}`
	err := json.Unmarshal([]byte(input), &u)
	errmsg_testing.AssertEquals(t, errmsg.Message(err), `'id' should be string but received number`)
	errmsg_testing.AssertType(t, errors.Unwrap(errmsg.Wrap(err)), &json.UnmarshalTypeError{})
}

func TestJsonTypeErrorForUintField(t *testing.T) {
	var c count
	input := `{"value": "x"}`
	err := json.Unmarshal([]byte(input), &c)
	errmsg_testing.AssertEquals(t, errmsg.Message(err), `'value' should be integer but received string`)
}

func TestJsonTypeErrorForNestedField(t *testing.T) {
	var u user
	input := `{"birthday": {"month": ""}}`
	err := json.Unmarshal([]byte(input), &u)
	errmsg_testing.AssertEquals(t, errmsg.Message(err), `'birthday.month' should be integer but received string`)
}

func TestJsonTypeErrorForNestedFieldOverflow(t *testing.T) {
	var u user
	input := `{"birthday": {"month": 4294967296}}`
	err := json.Unmarshal([]byte(input), &u)
	errmsg_testing.AssertEquals(t, errmsg.Message(err), `'birthday.month' should be int8 but received number 4294967296`)
}

func TestInvalidJsonError(t *testing.T) {
	var u user
	input := `{"id": }`
	err := json.Unmarshal([]byte(input), &u)
	errmsg_testing.AssertHasPrefix(t, errmsg.Message(err), "JSON syntax error")
	errmsg_testing.AssertType(t, errors.Unwrap(errmsg.Wrap(err)), &json.SyntaxError{})
}