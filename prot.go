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
	"io"
	"strconv"
)

// Request headers
const (
	_PROT_HEADER_AUTH       = 'A'
	_PROT_HEADER_OPEN_INDEX = 'P'
	_PROT_SEPARATOR         = '\t'
	_PROT_TERMINATOR        = '\n'
	_PROT_AUTH_TYPE         = '1'
	_PROT_NULL              = 0x00
)

func (hs *HandlerSocket) auth(secret string) error {
	var (
		err error
	)

	hs.createAuthReq(secret)

	if err = hs.writeRequest(); err != nil {
		return err
	}

	// read auth response
	if err = hs.readResponse(); err != nil {
		return err
	}

	return parseAuthResponse(hs.in)
}

func (hs *HandlerSocket) createAuthReq(secret string) error {
	hs.out.Reset()
	hs.out.WriteRune(_PROT_HEADER_AUTH)
	hs.out.WriteRune(_PROT_SEPARATOR)
	hs.out.WriteRune(_PROT_AUTH_TYPE)
	hs.out.WriteRune(_PROT_SEPARATOR)
	hs.out.WriteString(secret)
	hs.out.WriteRune(_PROT_TERMINATOR)
	return nil
}

func parseError(code uint16, b *bytes.Buffer) error {
	var (
		err error
		msg []byte
	)

	// move past [ <error code> \t '1' \t ]
	_, err = b.ReadString(_PROT_SEPARATOR)
	_, err = b.ReadString(_PROT_SEPARATOR)

	// check if we reached the end already..
	switch err {
	case nil:
		// there is a message, proceed
	case io.EOF:
		// there was no message
		return myError(code, "")
	default:
		return myError(ErrInvalidPacket)
	}

	msg, err = b.ReadBytes(_PROT_TERMINATOR)

	if err != nil {
		return myError(ErrInvalidPacket)
	} else {
		// discard the terminator
		return myError(code, string(msg[:len(msg)-1]))
	}
}

func parseAuthResponse(b *bytes.Buffer) error {
	if b.Len() < 4 {
		return myError(ErrInvalidPacket)
	}

	switch b.Bytes()[0] {
	case '0':
		return nil // success
	default:
		return parseError(ErrAuthentication, b)
	}
}

func (hs *HandlerSocket) openIndex(idx *Index) error {
	var (
		err error
	)

	hs.createOpenIndexReq(idx)

	if err = hs.writeRequest(); err != nil {
		return err
	}

	// read open_index response
	if err = hs.readResponse(); err != nil {
		return err
	}

	return parseOpenIndexResponse(hs.in)
}

func (hs *HandlerSocket) createOpenIndexReq(idx *Index) error {
	var (
		err    error
		length int
	)
	length = len(idx.columns)

	hs.out.Reset()
	hs.out.WriteRune(_PROT_HEADER_OPEN_INDEX)
	hs.out.WriteRune(_PROT_SEPARATOR)
	hs.out.WriteString(strconv.Itoa(idx.id))
	hs.out.WriteRune(_PROT_SEPARATOR)
	hs.out.WriteString(idx.schema)
	hs.out.WriteRune(_PROT_SEPARATOR)
	hs.out.WriteString(idx.table)
	hs.out.WriteRune(_PROT_SEPARATOR)
	hs.out.WriteString(idx.name)
	hs.out.WriteRune(_PROT_SEPARATOR)
	for i := 0; i < length; i++ {
		hs.out.WriteString(idx.columns[i])
		if i < length-1 {
			hs.out.WriteRune(',')
		}
	}
	hs.out.WriteRune(_PROT_TERMINATOR)

	return err
}

func parseOpenIndexResponse(b *bytes.Buffer) error {
	if b.Len() < 4 {
		return myError(ErrInvalidPacket)
	}

	switch b.Bytes()[0] {
	case '0':
		return nil // success
	default:
		return parseError(ErrOperationFailed, b)
	}
}

func (idx *Index) find() ([]Row, error) {
	var (
		err error
		hs  *HandlerSocket
	)

	hs = idx.hs

	idx.createFindReq()

	if err = hs.writeRequest(); err != nil {
		return nil, err
	}

	// read find response
	if err = hs.readResponse(); err != nil {
		return nil, err
	}

	switch hs.in.Bytes()[0] {
	case '0':
		return hs.parseResultSet()
	default:
		return nil, parseError(ErrOperationFailed, hs.in)
	}
	return nil, nil
}

func (idx *Index) createFindReq() error {
	var (
		err  error
		vLen int
		hs   *HandlerSocket
	)

	hs = idx.hs

	vLen = len(idx.values)

	hs.out.Reset()
	hs.out.WriteString(strconv.Itoa(idx.id))
	hs.out.WriteRune(_PROT_SEPARATOR)
	hs.out.WriteString(idx.op)
	hs.out.WriteRune(_PROT_SEPARATOR)
	hs.out.WriteString(strconv.Itoa(vLen))

	for i := 0; i < vLen; i++ {
		hs.out.WriteRune(_PROT_SEPARATOR)
		writeNullString(hs.out, idx.values[i])
	}

	hs.out.WriteRune(_PROT_TERMINATOR)
	return err
}

func (idx *Index) findModify() (Result, error) {
	var (
		err error
		res Result
		hs  *HandlerSocket
	)

	hs = idx.hs

	idx.createFindModifyReq()

	if err = hs.writeRequest(); err != nil {
		return res, err
	}

	// read find_modify response
	if err = hs.readResponse(); err != nil {
		return res, err
	}

	switch hs.in.Bytes()[0] {
	case '0':
		return hs.parseFindModifyResult()
	default:
		return res, parseError(ErrOperationFailed, hs.in)
	}
	return res, nil

}

func (idx *Index) createFindModifyReq() error {
	var (
		vLen  int
		inLen int
		mLen  int
		hs    *HandlerSocket
	)

	hs = idx.hs

	vLen = len(idx.values)
	inLen = len(idx.inValues)
	mLen = len(idx.mValues)

	hs.out.Reset()

	hs.out.WriteString(strconv.Itoa(idx.id))
	hs.out.WriteRune(_PROT_SEPARATOR)

	hs.out.WriteString(idx.op)
	hs.out.WriteRune(_PROT_SEPARATOR)

	hs.out.WriteString(strconv.Itoa(vLen))
	for i := 0; i < vLen; i++ {
		hs.out.WriteRune(_PROT_SEPARATOR)
		writeNullString(hs.out, idx.values[i])
	}

	// use default limit "1" if not specified
	hs.out.WriteRune(_PROT_SEPARATOR)
	if (idx.flags & _IDX_FLAG_LIMIT) > 0 {
		hs.out.WriteString(strconv.Itoa(idx.limit))
	} else {
		hs.out.WriteString("1")
	}

	// use default offset "0" if not specified
	hs.out.WriteRune(_PROT_SEPARATOR)
	if (idx.flags & _IDX_FLAG_OFFSET) > 0 {
		hs.out.WriteString(strconv.Itoa(idx.offset))
	} else {
		hs.out.WriteString("0")
	}

	if (idx.flags & _IDX_FLAG_IN) > 0 {
		hs.out.WriteRune(_PROT_SEPARATOR)
		hs.out.WriteRune('@')

		hs.out.WriteRune(_PROT_SEPARATOR)
		hs.out.WriteString(strconv.Itoa(idx.inColumn))

		hs.out.WriteRune(_PROT_SEPARATOR)
		hs.out.WriteString(strconv.Itoa(inLen))

		for i := 0; i < inLen; i++ {
			hs.out.WriteRune(_PROT_SEPARATOR)
			writeNullString(hs.out, idx.inValues[i])
		}
	}

	if (idx.flags & _IDX_FLAG_FILTER) > 0 {
		hs.out.WriteRune(_PROT_SEPARATOR)
		hs.out.WriteRune(idx.fType)

		hs.out.WriteRune(_PROT_SEPARATOR)
		hs.out.WriteString(idx.fOp)

		hs.out.WriteRune(_PROT_SEPARATOR)
		hs.out.WriteString(strconv.Itoa(idx.fColumn))

		hs.out.WriteRune(_PROT_SEPARATOR)
		writeNullString(hs.out, idx.fValue)
	}

	if (idx.flags & _IDX_FLAG_MODIFY_OP) > 1 {
		hs.out.WriteRune(_PROT_SEPARATOR)
		hs.out.WriteString(idx.mOp)
	}

	if (idx.flags & _IDX_FLAG_MODIFY_VALUES) > 1 {
		for i := 0; i < mLen; i++ {
			hs.out.WriteRune(_PROT_SEPARATOR)
			writeNullString(hs.out, idx.mValues[i])
		}
	}

	hs.out.WriteRune(_PROT_TERMINATOR)

	return nil
}

func (idx *Index) insert() error {
	var (
		err error
		hs  *HandlerSocket
	)

	hs = idx.hs

	idx.createInsertReq()

	if err = hs.writeRequest(); err != nil {
		return err
	}

	// read find_modify response
	if err = hs.readResponse(); err != nil {
		return err
	}

	switch hs.in.Bytes()[0] {
	case '0':
		return hs.parseInsertResult()
	default:
		return parseError(ErrOperationFailed, hs.in)
	}

	return nil

}

func (idx *Index) createInsertReq() error {
	var (
		vLen int
		hs   *HandlerSocket
	)

	hs = idx.hs

	vLen = len(idx.values)

	hs.out.Reset()

	hs.out.WriteString(strconv.Itoa(idx.id))
	hs.out.WriteRune(_PROT_SEPARATOR)

	hs.out.WriteString(idx.op)
	hs.out.WriteRune(_PROT_SEPARATOR)

	hs.out.WriteString(strconv.Itoa(vLen))
	for i := 0; i < vLen; i++ {
		hs.out.WriteRune(_PROT_SEPARATOR)
		writeNullString(hs.out, idx.values[i])
	}

	hs.out.WriteRune(_PROT_TERMINATOR)

	return nil
}

func (hs *HandlerSocket) parseFindModifyResult() (Result, error) {
	var (
		err          error
		b            []byte
		res          Result
		rowsAffected int
	)

	// move past [ '0' \t '1' \t ]
	hs.in.Next(4)

	if b, err = hs.in.ReadBytes(_PROT_TERMINATOR); err != nil {
		return res, err
	}

	// _PROT_SEPARATOR has also been read, must be discarded
	if rowsAffected, err = strconv.Atoi(string(b[:len(b)-1])); err != nil {
		return res, err
	}

	res.rowsAffected = int64(rowsAffected)
	return res, nil
}

func (hs *HandlerSocket) parseInsertResult() error {
	// move past [ '0' \t '1' 0x0A ]
	hs.in.Next(4)

	return nil
}

func (hs *HandlerSocket) parseResultSet() ([]Row, error) {
	var (
		err     error
		b       []byte
		rows    []Row
		numCols int
		done    bool
	)

	// TODO: verify that numCols must be equal to Index's len(columns)

	// read till _PROT_SEPARATOR
	hs.in.ReadBytes(_PROT_SEPARATOR)

	// read num of columns
	if b, err = hs.in.ReadBytes(_PROT_SEPARATOR); err != nil {
		// check if we reached past _PROT_TERMINATOR
		if err == io.EOF && len(b) >= 2 && b[len(b)-1] == byte(_PROT_TERMINATOR) {
			// empty result set
			return nil, nil
		}
		return nil, myError(ErrInvalidPacket, err)
	}

	// _PROT_SEPARATOR has also been read, must be discarded
	if numCols, err = strconv.Atoi(string(b[:len(b)-1])); err != nil {
		return nil, myError(ErrInvalidPacket, err)
	}

	rows = make([]Row, 0)

	for !done {
		row := new(Row)
		row.row = make([]NullString, 0)

		for i := 0; i < numCols; i++ {
			// read a value
			if b, err = hs.in.ReadBytes(_PROT_SEPARATOR); err != nil {
				// check whether this is a real error
				if err == io.EOF {
					done = true
				} else {
					return nil, err
				}
			}
			if b[0] == 0 {
				row.row = append(row.row, NullString{"", false})
			} else {
				row.row = append(row.row, NullString{string(b[:len(b)-1]), true})
			}
		}
		rows = append(rows, *row)
	}
	return rows, nil
}

// writeRequest writes the content in 'out' buffer to the network.
func (hs *HandlerSocket) writeRequest() error {
	var (
		err error
	)

	if _, err = hs.conn.Write(hs.out.Bytes()); err != nil {
		return err
	}
	return nil
}

// readResponse reads server's response from network into 'in' buffer.
func (hs *HandlerSocket) readResponse() error {
	var (
		buf [1]byte
		err error
		n   int
	)

	hs.in.Reset()

	for {
		if n, err = hs.conn.Read(buf[0:]); err != nil {
			if err != io.EOF {
				// real error
				return err
			}
		}

		hs.in.Write(buf[0:n])

		// 0x0A at the end signifies end of response
		if buf[n-1] == _PROT_TERMINATOR {
			break
		}
	}
	return nil
}

func writeNullString(b *bytes.Buffer, v NullString) (n int, err error) {
	if !v.Valid {
		return 1, b.WriteByte(0)
	}
	return b.WriteString(v.String)
}
