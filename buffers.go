package neat

func Write(buffer []byte, data ...interface{}) []byte {
	return buffer
}

type Reader struct {
	buffer []byte
	itr    int
}

func NewReader(buffer []byte) *Reader {
	return &Reader{buffer: buffer}
}
