package pkg

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Type byte

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Value
type Value struct {
	typ     Type
	str     string
	integer int
	array   []Value
	null    bool
}

func (v Value) Type() Type {
	return v.typ
}

func (v Value) IsArray() bool {
	return v.typ == ARRAY
}

func (v Value) IsBulk() bool {
	return v.typ == BULK
}

func (v Value) IsError() bool {
	return v.typ == ERROR
}

func (v Value) IsInteger() bool {
	return v.typ == INTEGER
}

func (v Value) IsString() bool {
	return v.typ == STRING
}

func (v Value) IsNull() bool {
	return v.null
}

func (v Value) Array() []Value {
	return v.array
}

func (v Value) Bulk() string {
	return v.str
}

func (v Value) Error() string {
	return v.str
}

func (v Value) Integer() int {
	return v.integer
}

func (v Value) String() string {
	return v.str
}

func BulkString(str string) Value {
	return Value{
		typ: BULK,
		str: str,
	}
}

// Resp
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd *bufio.Reader) *Resp {
	return &Resp{reader: rd}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n++
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	case '\n':
	case '\r':
	}
	return Value{}, fmt.Errorf("unknowntype: %v\n", string(_type))
}

func (r *Resp) readArray() (Value, error) {
	// read the length of the array
	len, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	if len > 1024*1024 {
		return Value{}, fmt.Errorf("array too long: %d", len)
	}
	var array []Value
	for i := 0; i < len; i++ {
		v, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		array = append(array, v)
	}
	return Value{typ: ARRAY, array: array}, nil
}

func (r *Resp) readBulk() (Value, error) {
	// read the length of the bulk string
	len, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	if len > 512*1024*1024 {
		return Value{}, fmt.Errorf("bulk string too long: %d", len)
	}

	buf := make([]byte, len)
	_, err = r.reader.Read(buf)
	if err != nil {
		return Value{}, err
	}
	// Read the trailing CRLF
	r.readLine()
	return Value{typ: BULK, str: string(buf)}, nil
}

// Writer
type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

func (w *Writer) WriteSimpleString(str string) error {
	_, err := w.w.Write([]byte("+" + str + "\r\n"))
	return err
}

func (w *Writer) WriteError(str string) error {
	_, err := w.w.Write([]byte("-" + str + "\r\n"))
	return err
}

func (w *Writer) WriteInteger(i int) error {
	_, err := w.w.Write([]byte(":" + strconv.Itoa(i) + "\r\n"))
	return err
}

func (w *Writer) WriteBulkString(str string) error {
	_, err := w.w.Write([]byte("$" + strconv.Itoa(len(str)) + "\r\n" + str + "\r\n"))
	return err
}

func (w *Writer) WriteNull() error {
	_, err := w.w.Write([]byte("$-1\r\n"))
	return err
}

func (w *Writer) WriteArray(value Value) error {
	_, err := w.w.Write([]byte("*" + strconv.Itoa(len(value.array)) + "\r\n"))

	if err != nil {
		return err
	}

	for _, v := range value.array {
		switch v.typ {
		case STRING:
			err = w.WriteSimpleString(v.str)
		case ERROR:
			err = w.WriteError(v.str)
		case INTEGER:
			err = w.WriteInteger(v.integer)
		case BULK:
			err = w.WriteBulkString(v.str)
		}
	}
	return err
}

func (v Value) MarshalResp() ([]byte, error) {
	return marshalResp(v)
}

func marshalResp(v Value) ([]byte, error) {
	switch v.typ {
	case STRING:
		return []byte("+" + v.str + "\r\n"), nil
	case ERROR:
		return []byte("-" + v.str + "\r\n"), nil
	case INTEGER:
		return []byte(":" + strconv.Itoa(v.integer) + "\r\n"), nil
	case BULK:
		return []byte("$" + strconv.Itoa(len(v.str)) + "\r\n" + v.str + "\r\n"), nil
	case ARRAY:
		buf := []byte("*" + strconv.Itoa(len(v.array)) + "\r\n")
		for _, value := range v.array {
			m, err := marshalResp(value)
			if err != nil {
				return nil, err
			}
			buf = append(buf, m...)
		}
		return buf, nil
	}
	return nil, fmt.Errorf("unknown type: %v", v.typ)
}

// Reader
type Reader struct {
	r    *bufio.Reader
	resp *Resp
}

func NewReader(r io.Reader) *Reader {
	rd := bufio.NewReader(r)
	return &Reader{
		r:    rd,
		resp: NewResp(rd),
	}
}
