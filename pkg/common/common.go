// Released under an MIT license. See LICENSE.

package common

import "errors"

type ReadStringer interface {
	ReadString(delim byte) (line string, err error)
}

const (
	ErrNotExecutable = "oh: 126: error/runtime: "
	ErrNotFound      = "oh: 127: error/runtime: "
	ErrSyntax        = "oh: 1: error/syntax: "
)

var CtrlCPressed = errors.New("ctrl-c pressed")
