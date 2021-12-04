// Package util
// @Description
// @Author  liuguilin）  2021/10/12 00:46
// @Update   2021/10/12 00:46
package util

import "errors"

var ErrEof = errors.New("EOF")
var ErrClosedPipe = errors.New("io: read/write on closed pipe")
var ERR_NO_PROGRESS = errors.New("multiple Read calls return no data or error")
var ERR_SHORT_BUFFER = errors.New("short buffer")
var ERR_SHORT_WRITE = errors.New("short write")
var ERR_UNEXPECTED_EOF = errors.New("unexpected EOF")
