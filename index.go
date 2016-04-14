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

const (
	_IDX_FLAG_OP = 1 << iota
	_IDX_FLAG_VALUES
	_IDX_FLAG_LIMIT
	_IDX_FLAG_OFFSET
	_IDX_FLAG_IN
	_IDX_FLAG_FILTER
	_IDX_FLAG_MODIFY_OP
	_IDX_FLAG_MODIFY_VALUES
)

type Index struct {
	hs      *HandlerSocket
	id      int
	name    string
	schema  string
	table   string
	columns []string

	flags    uint64
	op       string
	limit    int
	offset   int
	inColumn int
	fType    rune
	fOp      string
	fColumn  int
	fValue   NullString
	mOp      string

	values   []NullString
	inValues []NullString
	mValues  []NullString
}

// <!-- command builder -->

func (idx *Index) Operator(op string) *Index {
	idx.op = op
	idx.flags |= _IDX_FLAG_OP
	return idx
}

func (idx *Index) Values(values []NullString) *Index {
	if values != nil {
		idx.values = values
		idx.flags |= _IDX_FLAG_VALUES
	}
	return idx
}

func (idx *Index) Limit(limit int) *Index {
	idx.limit = limit
	idx.flags |= _IDX_FLAG_LIMIT
	return idx
}

func (idx *Index) Offset(offset int) *Index {
	idx.offset = offset
	idx.flags |= _IDX_FLAG_OFFSET
	return idx
}

func (idx *Index) In(column int, values []NullString) *Index {
	idx.inColumn = column

	if values != nil {
		idx.inValues = values
		idx.flags |= _IDX_FLAG_IN
	}

	return idx
}

func (idx *Index) Filter(typ rune, op string, column int, value NullString) *Index {
	idx.fType = typ
	idx.fOp = op
	idx.fColumn = column
	idx.fValue = value

	idx.flags |= _IDX_FLAG_FILTER
	return idx
}

func (idx *Index) Reset() *Index {
	idx.flags = 0
	return idx
}

func (idx *Index) Select() ([]Row, error) {
	return idx.find()
}

func (idx *Index) Insert(values []NullString) error {
	idx.flags |= _IDX_FLAG_OP
	idx.op = "+"

	if values != nil {
		idx.flags |= _IDX_FLAG_VALUES
		idx.values = values
	}
	return idx.insert()
}

func (idx *Index) Update(values []NullString) (Result, error) {
	idx.flags |= _IDX_FLAG_MODIFY_OP
	idx.mOp = "U"

	if values != nil {
		idx.flags |= _IDX_FLAG_MODIFY_VALUES
		idx.mValues = values
	}
	return idx.findModify()
}

func (idx *Index) Delete() (Result, error) {
	idx.flags |= _IDX_FLAG_MODIFY_OP
	idx.mOp = "D"
	return idx.findModify()
}
