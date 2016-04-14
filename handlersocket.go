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
	"bytes"
	"net"
)

type HandlerSocket struct {
	conn net.Conn
	in   *bytes.Buffer
	out  *bytes.Buffer
}

func Connect(address, secret string) (*HandlerSocket, error) {
	var (
		err error
		hs  *HandlerSocket
	)

	if hs, err = connect(address); err != nil {
		return nil, err
	}

	// authentication
	if err = hs.auth(secret); err != nil {
		return nil, err
	}

	return hs, nil
}

func connect(address string) (*HandlerSocket, error) {
	var (
		err error
		hs  *HandlerSocket
	)

	hs = new(HandlerSocket)

	if hs.conn, err = net.Dial("tcp", address); err != nil {
		return nil, err
	}

	hs.in = bytes.NewBuffer(nil)
	hs.out = bytes.NewBuffer(nil)

	return hs, nil
}

func (hs *HandlerSocket) Close() error {
	hs.conn.Close()
	return nil
}

func (hs *HandlerSocket) OpenIndex(id int, name, schema, table string, columns []string) (*Index, error) {
	var (
		idx *Index
		err error
	)

	idx = new(Index)
	idx.id = id
	idx.name = name
	idx.schema = schema
	idx.table = table
	idx.columns = columns

	if err = hs.openIndex(idx); err != nil {
		return nil, err
	}

	idx.hs = hs
	return idx, nil
}
