package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file   *os.File
	rd     *bufio.Reader
	mu     sync.Mutex
	closed bool
	atEnd  bool
}

func NewAof(file string) (*Aof, error) {

	if file == "" {
		file = "pkg.aof"
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}
	// start a goroutine to flush the buffer to disk every second
	go func() {
		for {
			aof.mu.Lock()
			if aof.closed {
				aof.mu.Unlock()
				return
			}
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()
	return aof, nil
}

func (a *Aof) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.closed = true
	return a.file.Close()
}

func (a *Aof) ReadValues(iterator func(Value) bool) error {
	a.atEnd = false
	if _, err := a.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	rd := NewReader(a.file)
	for {
		v, err := rd.resp.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read value failed: ", err)
			return err
		}

		if iterator != nil && !iterator(v) {
			return nil
		}
	}

	_, err := a.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	a.atEnd = true
	return nil
}

func (a *Aof) Write(b []byte) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.closed {
		return 0, fmt.Errorf("aof is closed")
	}
	return a.file.Write(b)
}

func (a *Aof) Append(v Value) error {
	return a.AppendMany([]Value{v})
}

func (a *Aof) AppendMany(vs []Value) error {
	var buf []byte
	for _, v := range vs {
		b, err := v.MarshalResp()
		if err != nil {
			return err
		}

		buf = append(buf, b...)

	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return fmt.Errorf("aof is closed")
	}

	if !a.atEnd {
		a.ReadValues(nil)
		if !a.atEnd {
			return fmt.Errorf("aof is not at end")
		}
	}

	_, err := a.file.Write(buf)

	return err
}
