package main

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unicode/utf16"
)

func TestPickleWriter_WriteUInt32(t *testing.T) {
	w := NewPickleWriter()
	w.WriteUInt32(42)

	payload := w.GetPayload()
	if len(payload) != 4 {
		t.Errorf("Expected length 4, got %d", len(payload))
	}

	var val uint32
	binary.Read(bytes.NewReader(payload), binary.LittleEndian, &val)
	if val != 42 {
		t.Errorf("Expected value 42, got %d", val)
	}
}

func TestPickleWriter_WriteString16(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLen  int // Expected total bytes written (length + data + padding)
		wantChar uint32
	}{
		{
			name:     "Empty String",
			input:    "",
			wantLen:  4, // 4 bytes for length (0)
			wantChar: 0,
		},
		{
			name:     "ASCII String",
			input:    "test",
			wantLen:  4 + 8, // 4 (len) + 8 (4 chars * 2 bytes) + 0 padding
			wantChar: 4,
		},
		{
			name:     "Odd Length String",
			input:    "abc",
			wantLen:  4 + 6 + 2, // 4 (len) + 6 (3 chars * 2 bytes) + 2 padding
			wantChar: 3,
		},
		{
			name:     "Japanese String",
			input:    "あいう",
			wantLen:  4 + 6 + 2, // 4 (len) + 6 (3 chars * 2 bytes) + 2 padding
			wantChar: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewPickleWriter()
			w.WriteString16(tt.input)

			payload := w.GetPayload()
			if len(payload) != tt.wantLen {
				t.Errorf("Expected total length %d, got %d", tt.wantLen, len(payload))
			}

			// Verify Length (first 4 bytes)
			var charCount uint32
			binary.Read(bytes.NewReader(payload[0:4]), binary.LittleEndian, &charCount)
			if charCount != tt.wantChar {
				t.Errorf("Expected char count %d, got %d", tt.wantChar, charCount)
			}

			// Verify Data
			expectedRunes := utf16.Encode([]rune(tt.input))
			dataBytes := payload[4 : 4+len(expectedRunes)*2]
			
			readRunes := make([]uint16, len(expectedRunes))
			binary.Read(bytes.NewReader(dataBytes), binary.LittleEndian, &readRunes)

			for i, r := range expectedRunes {
				if readRunes[i] != r {
					t.Errorf("Char mismatch at index %d: expected %x, got %x", i, r, readRunes[i])
				}
			}
		})
	}
}

func TestPickleWriter_ComplexStructure(t *testing.T) {
	// Simulate the structure used in main.go
	// Entry Count (2) + Key1 + Value1
	w := NewPickleWriter()
	
	w.WriteUInt32(2)
	w.WriteString16("key")
	w.WriteString16("val")

	payload := w.GetPayload()
	
	// Expected size:
	// Count: 4 bytes
	// Key "key": 4 (len) + 6 (data) + 2 (padding) = 12 bytes
	// Val "val": 4 (len) + 6 (data) + 2 (padding) = 12 bytes
	// Total: 28 bytes
	expectedSize := 4 + 12 + 12
	
	if len(payload) != expectedSize {
		t.Errorf("Expected complex payload size %d, got %d", expectedSize, len(payload))
	}
}
