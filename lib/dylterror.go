package lib

import (
	"fmt"
	"runtime"
)

type DyltError struct {
	Err error
	Filename string
	Linenum int
}

func NewError (baseErr error) *DyltError {
	_, filename, linenum, ok := runtime.Caller(1)
	if !ok {
		filename = "N/A"
		linenum = 0
	}	

	err := DyltError{
		Err: baseErr,
		Filename: filename,
		Linenum: linenum,
	}

	return &err
}

func (err *DyltError) Error () string {
	return fmt.Sprintf("%s\n\t%s:%d", err.Err, err.Filename, err.Linenum)
}