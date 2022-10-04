package main

import (
	"net/http"
)

type IntErr int

func ParseIntErr(s string) IntErr {
	t := EUnknown

	switch s {
	case "Unknown":
		t = EUnknown
	case "Validation":
		t = EValidation
	case "NotFound":
		t = ENotFound
	case "Internal":
		t = EInternal
	case "Duplicate":
		t = EDuplicate
	case "Empty":
		t = EEmpty
	case "InputBody":
		t = EInputBody
	}

	return t
}

const (
	ENo IntErr = iota
	EUnknown
	EValidation
	ENotFound
	EInternal
	EDuplicate
	EEmpty
	EInputBody
)

func (r IntErr) String() string {
	t := "Unknown"

	switch r {
	case EUnknown:
		t = "Unknown"
	case EValidation:
		t = "Validation"
	case ENotFound:
		t = "NotFound"
	case EInternal:
		t = "Internal"
	case EDuplicate:
		t = "Duplicate"
	case EEmpty:
		t = "Empty"
	case EInputBody:
		t = "InputBody"
	}

	return t
}

func (r IntErr) Status() int {
	status := http.StatusInternalServerError

	switch r {
	case EValidation:
		status = http.StatusUnprocessableEntity
	case ENotFound:
		status = http.StatusNotFound
	case EDuplicate:
		status = http.StatusConflict
	case EEmpty:
		status = http.StatusGone
	case EInputBody:
		status = http.StatusBadRequest
	}

	return status
}
