package neat

import "math"

type Buffer []byte

func NewBuffer(v ...[]byte) Buffer {
	if len(v) == 1 {
		return Buffer(v[0])
	}
	return Buffer(make([]byte, 0))
}

func (b Buffer) WriteBool(v bool) Buffer {
	var value uint8
	if v {
		value++
	}
	return append(b, byte(value))
}

func (b Buffer) WriteU8(v uint8) Buffer {
	return append(b, byte(v))
}

func (b Buffer) Write8(v int8) Buffer {
	return b.WriteU8(uint8(v))
}

func (b Buffer) WriteU16(v uint16) Buffer {
	b = append(b, byte(v), byte(v>>8))
	return b
}

func (b Buffer) Write16(v int16) Buffer {
	return b.WriteU16(uint16(v))
}

func (b Buffer) WriteU32(v uint32) Buffer {
	b = append(b, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	return b
}

func (b Buffer) Write32(v int32) Buffer {
	return b.WriteU32(uint32(v))
}

func (b Buffer) WriteF32(v float32) Buffer {
	return b.WriteU32(math.Float32bits(v))
}

func (b Buffer) WriteU64(v uint64) Buffer {
	b = append(b, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
	return b
}

func (b Buffer) Write64(v int64) Buffer {
	return b.WriteU64(uint64(v))
}

func (b Buffer) WriteF64(v float64) Buffer {
	return b.WriteU64(math.Float64bits(v))
}

func (b Buffer) WriteStr(v string) Buffer {
	bstring := []byte(v)
	len := uint16(len(bstring))
	b = append(b, byte(len), byte(len>>8))
	b = append(b, bstring...)
	return b
}

func (b Buffer) WriteBytes(v []byte) Buffer {
	len := uint16(len(v))
	b = append(b, byte(len), byte(len>>8))
	b = append(b, v...)
	return b
}

func (b Buffer) Write(v ...interface{}) Buffer {
	for i := range v {
		switch value := v[i].(type) {
		case bool:
			b = b.WriteBool(value)
		case uint8:
			b = b.WriteU8(value)
		case int8:
			b = b.Write8(value)
		case uint16:
			b = b.WriteU16(value)
		case int16:
			b = b.Write16(value)
		case uint32:
			b = b.WriteU32(value)
		case int32:
			b = b.Write32(value)
		case float32:
			b = b.WriteF32(value)
		case uint64:
			b = b.WriteU64(value)
		case int64:
			b = b.Write64(value)
		case float64:
			b = b.WriteF64(value)
		case string:
			b = b.WriteStr(value)
		case []byte:
			b = b.WriteBytes(value)
		}
	}
	return b
}

func (b Buffer) Readable() *BufferReadable {
	return &BufferReadable{Buffer: b}
}

type BufferReadable struct {
	Buffer
	Index int
}

func (b *BufferReadable) ReadBool() bool {
	v := uint8(b.Buffer[b.Index])
	b.Index++
	return v == 1
}

func (b *BufferReadable) ReadU8() uint8 {
	v := uint8(b.Buffer[b.Index])
	b.Index++
	return v
}

func (b *BufferReadable) ReadU16() uint16 {
	v := uint16(b.Buffer[b.Index]) | uint16(b.Buffer[b.Index+1])<<8
	b.Index += 2
	return v
}

func (b *BufferReadable) Read16() int16 {
	return int16(b.ReadU16())
}

func (b *BufferReadable) ReadU32() uint32 {
	v := uint32(b.Buffer[b.Index]) | uint32(b.Buffer[b.Index+1])<<8 | uint32(b.Buffer[b.Index+2])<<16 | uint32(b.Buffer[b.Index+3])<<24
	b.Index += 4
	return v
}

func (b *BufferReadable) Read32() int32 {
	return int32(b.ReadU32())
}

func (b *BufferReadable) ReadF32() float32 {
	return math.Float32frombits(b.ReadU32())
}

func (b *BufferReadable) ReadU64() uint64 {
	v := uint64(b.Buffer[b.Index]) | uint64(b.Buffer[b.Index+1])<<8 | uint64(b.Buffer[b.Index+2])<<16 | uint64(b.Buffer[b.Index+3])<<24 | uint64(b.Buffer[b.Index+4])<<32 | uint64(b.Buffer[b.Index+5])<<40 | uint64(b.Buffer[b.Index+6])<<48 | uint64(b.Buffer[b.Index+7])<<56
	b.Index += 8
	return v
}

func (b *BufferReadable) Read64() int64 {
	return int64(b.ReadU64())
}

func (b *BufferReadable) ReadF64() float64 {
	return math.Float64frombits(b.ReadU64())
}

func (b *BufferReadable) ReadStr() string {
	len := int(b.ReadU16())
	v := string(b.Buffer[b.Index : b.Index+len])
	b.Index += len
	return v
}

func (b *BufferReadable) ReadBytes() []byte {
	len := int(b.ReadU16())
	v := b.Buffer[b.Index : b.Index+len]
	b.Index += len
	return v
}
