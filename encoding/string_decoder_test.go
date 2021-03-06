/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| encoding/string_decoder_test.go                          |
|                                                          |
| LastModified: Jun 2, 2020                                |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoding

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeString(t *testing.T) {
	sb := new(strings.Builder)
	enc := NewEncoder(sb)
	enc.Encode(-1)
	enc.Encode(0)
	enc.Encode(1)
	enc.Encode(123)
	enc.Encode(math.MinInt64)
	enc.Encode(-math.MaxInt64)
	enc.Encode(math.MaxInt64)
	enc.Encode(true)
	enc.Encode(false)
	enc.Encode(nil)
	enc.Encode(3.14)
	enc.Encode(math.NaN())
	enc.Encode(math.Inf(1))
	enc.Encode(math.Inf(-1))
	enc.Encode("")
	enc.Encode("1")
	enc.Encode("123")
	enc.Encode("N")
	enc.Encode("NaN")
	enc.Encode("Pokémon")
	enc.Encode("中文")
	enc.Encode("🐱🐶")
	enc.Encode("👩‍👩‍👧‍👧")
	dec := NewDecoder(([]byte)(sb.String()))
	var s string
	dec.Decode(&s)
	assert.Equal(t, "-1", s)
	dec.Decode(&s)
	assert.Equal(t, "0", s)
	dec.Decode(&s)
	assert.Equal(t, "1", s)
	dec.Decode(&s)
	assert.Equal(t, "123", s)
	dec.Decode(&s)
	assert.Equal(t, strconv.FormatInt(math.MinInt64, 10), s)
	dec.Decode(&s)
	assert.Equal(t, strconv.FormatInt(-math.MaxInt64, 10), s)
	dec.Decode(&s)
	assert.Equal(t, strconv.FormatInt(math.MaxInt64, 10), s)
	dec.Decode(&s)
	assert.Equal(t, "true", s)
	dec.Decode(&s)
	assert.Equal(t, "false", s)
	dec.Decode(&s)
	assert.Equal(t, "", s)
	dec.Decode(&s)
	assert.Equal(t, "3.14", s)
	dec.Decode(&s)
	assert.Equal(t, "NaN", s)
	dec.Decode(&s)
	assert.Equal(t, "+Inf", s)
	dec.Decode(&s)
	assert.Equal(t, "-Inf", s)
	dec.Decode(&s)
	assert.Equal(t, "", s)
	dec.Decode(&s)
	assert.Equal(t, "1", s)
	dec.Decode(&s)
	assert.Equal(t, "123", s)
	assert.NoError(t, dec.Error)
	dec.Decode(&s)
	assert.Equal(t, "N", s)
	dec.Decode(&s)
	assert.Equal(t, "NaN", s)
	dec.Decode(&s)
	assert.Equal(t, "Pokémon", s)
	dec.Decode(&s)
	assert.Equal(t, "中文", s)
	dec.Decode(&s)
	assert.Equal(t, "🐱🐶", s)
	dec.Decode(&s)
	assert.Equal(t, "👩‍👩‍👧‍👧", s)
}

func TestDecodeStringFromReader(t *testing.T) {
	sb := new(strings.Builder)
	enc := NewEncoder(sb)
	enc.Encode(-1)
	enc.Encode(0)
	enc.Encode(1)
	enc.Encode(123)
	enc.Encode(math.MinInt64)
	enc.Encode(-math.MaxInt64)
	enc.Encode(math.MaxInt64)
	enc.Encode(true)
	enc.Encode(false)
	enc.Encode(nil)
	enc.Encode(3.14)
	enc.Encode(math.NaN())
	enc.Encode(math.Inf(1))
	enc.Encode(math.Inf(-1))
	enc.Encode("")
	enc.Encode("1")
	enc.Encode("123")
	enc.Encode("N")
	enc.Encode("NaN")
	enc.Encode("Pokémon")
	enc.Encode("中文")
	enc.Encode("🐱🐶")
	enc.Encode("👩‍👩‍👧‍👧")
	dec := NewDecoderFromReader(bytes.NewReader(([]byte)(sb.String())), 32)
	var s *string
	dec.Decode(&s)
	assert.Equal(t, "-1", *s)
	dec.Decode(&s)
	assert.Equal(t, "0", *s)
	dec.Decode(&s)
	assert.Equal(t, "1", *s)
	dec.Decode(&s)
	assert.Equal(t, "123", *s)
	dec.Decode(&s)
	assert.Equal(t, strconv.FormatInt(math.MinInt64, 10), *s)
	dec.Decode(&s)
	assert.Equal(t, strconv.FormatInt(-math.MaxInt64, 10), *s)
	dec.Decode(&s)
	assert.Equal(t, strconv.FormatInt(math.MaxInt64, 10), *s)
	dec.Decode(&s)
	assert.Equal(t, "true", *s)
	dec.Decode(&s)
	assert.Equal(t, "false", *s)
	dec.Decode(&s)
	assert.Nil(t, s)
	dec.Decode(&s)
	assert.Equal(t, "3.14", *s)
	dec.Decode(&s)
	assert.Equal(t, "NaN", *s)
	dec.Decode(&s)
	assert.Equal(t, "+Inf", *s)
	dec.Decode(&s)
	assert.Equal(t, "-Inf", *s)
	dec.Decode(&s)
	assert.Equal(t, "", *s)
	dec.Decode(&s)
	assert.Equal(t, "1", *s)
	dec.Decode(&s)
	assert.Equal(t, "123", *s)
	assert.NoError(t, dec.Error)
	dec.Decode(&s)
	assert.Equal(t, "N", *s)
	dec.Decode(&s)
	assert.Equal(t, "NaN", *s)
	dec.Decode(&s)
	assert.Equal(t, "Pokémon", *s)
	dec.Decode(&s)
	assert.Equal(t, "中文", *s)
	dec.Decode(&s)
	assert.Equal(t, "🐱🐶", *s)
	dec.Decode(&s)
	assert.Equal(t, "👩‍👩‍👧‍👧", *s)
}
