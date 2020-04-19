/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| encoding/decoder.go                                      |
|                                                          |
| LastModified: Apr 19, 2020                               |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoding

import (
	"bytes"
	"io"
	"reflect"
)

const defaultBufferSize = 8192

// Decoder is a io.Reader like object, with hprose specific read functions.
// Error is not returned as return value, but stored as Error member on this decoder instance.
type Decoder struct {
	reader io.Reader
	buf    []byte
	head   int
	tail   int
	refer  *decoderRefer
	ref    []reflect.Type
	Error  error
}

// NewDecoder creates an Decoder instance from byte array
func NewDecoder(input []byte) *Decoder {
	return &Decoder{
		reader: nil,
		buf:    input,
		head:   0,
		tail:   len(input),
	}
}

// NewDecoderFromReader creates an Decoder instance from io.Reader
func NewDecoderFromReader(reader io.Reader, bufSize int) *Decoder {
	if bufSize < 32 {
		bufSize = defaultBufferSize
	}
	return &Decoder{
		reader: reader,
		buf:    make([]byte, bufSize),
		head:   0,
		tail:   0,
	}
}

// Decode a data from the Decoder
func (dec *Decoder) Decode(p interface{}) {
	tag := dec.NextByte()
	switch p.(type) {
	case *int:
		intdec.Decode(dec, p, tag)
	case *int8:
		int8dec.Decode(dec, p, tag)
	case *int16:
		int16dec.Decode(dec, p, tag)
	case *int32:
		int32dec.Decode(dec, p, tag)
	case *int64:
		int64dec.Decode(dec, p, tag)
	case *uint:
		uintdec.Decode(dec, p, tag)
	case *uint8:
		uint8dec.Decode(dec, p, tag)
	case *uint16:
		uint16dec.Decode(dec, p, tag)
	case *uint32:
		uint32dec.Decode(dec, p, tag)
	case *uint64:
		uint64dec.Decode(dec, p, tag)
	case *string:
		strdec.Decode(dec, p, tag)
	}
}

// Reset the value reference and struct type reference
func (dec *Decoder) Reset() {
	if dec.refer != nil {
		dec.refer.Reset()
	}
	dec.ref = dec.ref[:0]
}

// Simple resets the encoder to simple mode or not
func (dec *Decoder) Simple(simple bool) {
	if simple {
		dec.refer = nil
	} else {
		dec.refer = &decoderRefer{}
	}
	dec.ref = dec.ref[:0]
}

// AddReference adds o to the reference
func (dec *Decoder) AddReference(o interface{}) {
	if dec.refer != nil {
		dec.refer.Add(o)
	}
}

// SetReference sets o to the reference at index i
func (dec *Decoder) SetReference(i int, o interface{}) {
	if dec.refer != nil {
		dec.refer.Set(i, o)
	}
}

// LastReferenceIndex returns the last index of the reference
func (dec *Decoder) LastReferenceIndex() int {
	if dec.refer != nil {
		dec.refer.Last()
	}
	return -1
}

// ResetReader reuse decoder instance by specifying another reader
func (dec *Decoder) ResetReader(reader io.Reader) *Decoder {
	dec.reader = reader
	dec.head = 0
	dec.tail = 0
	return dec
}

// ResetBytes reuse decoder instance by specifying another byte array as input
func (dec *Decoder) ResetBytes(input []byte) *Decoder {
	dec.reader = nil
	dec.buf = input
	dec.head = 0
	dec.tail = len(input)
	return dec
}

// NextByte reads and returns the next byte from the dec. If no byte is available, it returns 0.
func (dec *Decoder) NextByte() (b byte) {
	if (dec.head == dec.tail) && !dec.loadMore() {
		return 0
	}
	b = dec.buf[dec.head]
	dec.head++
	return b
}

// Skip the next byte from the dec.
func (dec *Decoder) Skip() {
	if (dec.head == dec.tail) && !dec.loadMore() {
		return
	}
	dec.head++
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read method.
func (dec *Decoder) Next(n int) (data []byte) {
	if (dec.head == dec.tail) && !dec.loadMore() {
		return
	}
	remain := dec.tail - dec.head
	if remain >= n {
		data = dec.buf[dec.head : dec.head+n]
		dec.head += n
		return
	}
	data = make([]byte, 0, n)
	data = append(data, dec.buf[dec.head:dec.tail]...)
	n -= remain
	for {
		if !dec.loadMore() {
			return
		}
		if dec.tail >= n {
			dec.head = n
			return append(data, dec.buf[0:n]...)
		}
		data = append(data, dec.buf[:dec.tail]...)
		n -= dec.tail
	}
}

// Remains reads and returns all bytes data in this iter that has not been read.
func (dec *Decoder) Remains() (data []byte) {
	if (dec.head == dec.tail) && !dec.loadMore() {
		return
	}
	for {
		data = append(data, dec.buf[dec.head:dec.tail]...)
		if !dec.loadMore() {
			return data
		}
	}
}

// Until reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and not including the delimiter.
func (dec *Decoder) Until(delim byte) (data []byte) {
	if (dec.head == dec.tail) && !dec.loadMore() {
		return
	}
	if i := bytes.IndexByte(dec.buf[dec.head:dec.tail], delim); i >= 0 {
		data = dec.buf[dec.head : dec.head+i]
		dec.head += i + 1
		return
	}
	for {
		data = append(data, dec.buf[dec.head:dec.tail]...)
		if !dec.loadMore() {
			return
		}
		if i := bytes.IndexByte(dec.buf[dec.head:dec.tail], delim); i >= 0 {
			data = append(data, dec.buf[dec.head:dec.head+i]...)
			dec.head += i + 1
			return
		}
	}
}

func (dec *Decoder) loadMore() bool {
	if dec.reader == nil {
		dec.head = dec.tail
		if dec.Error == nil {
			dec.Error = io.EOF
		}
		return false
	} else if dec.buf == nil {
		dec.buf = make([]byte, defaultBufferSize)
	}
	for {
		n, err := dec.reader.Read(dec.buf)
		dec.head = 0
		dec.tail = n
		if n > 0 {
			return true
		}
		if err != nil {
			if dec.Error == nil {
				dec.Error = err
			}
			return false
		}
	}
}

func (dec *Decoder) castError(p interface{}) {
	var iface interface{}
	dec.Decode(&iface)
	dec.Error = &CastError{
		Source:      reflect.TypeOf(iface),
		Destination: reflect.TypeOf(p).Elem(),
	}
}
