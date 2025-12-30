package src

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

type PickleWriter struct {
	buffer bytes.Buffer
}

func NewPickleWriter() *PickleWriter {
	return &PickleWriter{}
}

func (w *PickleWriter) WriteUInt32(value uint32) {
	binary.Write(&w.buffer, binary.LittleEndian, value)
}

func (w *PickleWriter) WriteString16(s string) {
	// Encode to UTF-16
	runes := []rune(s)
	encoded := utf16.Encode(runes)

	// Length is number of characters (uint16 units)
	charCount := uint32(len(encoded))
	w.WriteUInt32(charCount)

	// Write data (Little Endian for each uint16)
	for _, v := range encoded {
		binary.Write(&w.buffer, binary.LittleEndian, v)
	}

	// Padding
	// Total bytes written for string data = charCount * 2
	bytesWritten := int(charCount * 2)
	padding := (4 - (bytesWritten % 4)) % 4
	if padding > 0 {
		w.buffer.Write(make([]byte, padding))
	}
}

func (w *PickleWriter) GetPayload() []byte {
	return w.buffer.Bytes()
}
