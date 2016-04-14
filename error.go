/*
  The MIT License (MIT)

  Copyright (c) 2016 Nirbhay Choubey

  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:

  The above copyright notice and this permission notice shall be included in all
  copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  SOFTWARE.
*/

package handlersocket

import (
	"fmt"
	"time"
)

type Error struct {
	code    uint16
	message string
	when    time.Time
}

// error codes
const (
	ErrUnknown = 10000 + iota
	ErrAuthentication
	ErrConnection
	ErrOperationFailed
	ErrInvalidPacket
)

var errFormat = map[uint16]string{
	ErrUnknown:         "Unknown error",
	ErrAuthentication:  "Authentication error (%s)",
	ErrConnection:      "Can't connect to the server (%s)",
	ErrOperationFailed: "Operation failed (%s)",
	ErrInvalidPacket:   "Invalid/unexpected packet received",
}

func myError(code uint16, a ...interface{}) *Error {
	return &Error{code: code,
		message: fmt.Sprintf(errFormat[code], a...),
		when:    time.Now()}
}

// Error returns the formatted error message (also required by Go's error
// interface)
func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.code, e.message)
}

// Code returns the error number.
func (e *Error) Code() uint16 {
	return e.code
}

// Message returns the error message.
func (e *Error) Message() string {
	return e.message
}

// When returns the time then error occured.
func (e *Error) When() time.Time {
	return e.when
}
