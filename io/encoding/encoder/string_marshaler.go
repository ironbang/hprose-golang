/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| io/encoding/encoder/string_marshaler.go                  |
|                                                          |
| LastModified: Mar 1, 2020                                |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoder

import (
	"github.com/hprose/hprose-golang/v3/io"
	"github.com/modern-go/reflect2"
)

// StringMarshaler is the implementation of Marshaler for string.
type StringMarshaler struct{}

var stringMarshaler StringMarshaler

func (m StringMarshaler) encode(enc *Encoder, s string) (err error) {
	length := utf16Length(s)
	switch length {
	case 0:
		err = enc.Writer.WriteByte(io.TagEmpty)
	case 1:
		if err = enc.Writer.WriteByte(io.TagUTF8Char); err == nil {
			_, err = enc.Writer.Write(reflect2.UnsafeCastString(s))
		}
	default:
		var ok bool
		if ok, err = enc.WriteStringReference(s); !ok && err == nil {
			enc.SetStringReference(s)
			err = writeString(enc.Writer, s, length)
		}
	}
	return
}

func (m StringMarshaler) write(enc *Encoder, s string) (err error) {
	enc.SetStringReference(s)
	return writeString(enc.Writer, s, utf16Length(s))
}

// Encode writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as reference
func (m StringMarshaler) Encode(enc *Encoder, v interface{}) (err error) {
	return m.encode(enc, v.(string))
}

// Write writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as value
func (m StringMarshaler) Write(enc *Encoder, v interface{}) (err error) {
	return m.write(enc, v.(string))
}
