package bkperrors

import "errors"

var (
	ErrCannotAccessSrcDir = errors.New("can't access source directory")
	ErrCannotAccessDstDir = errors.New("can't access source directory")
)
