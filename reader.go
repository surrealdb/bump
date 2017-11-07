// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bump

import (
	"io"
)

const readerSize = 1024

// Reader represents a buffer for reading
// from an io.Reader, or a byte slice.
type Reader struct {
	pos int
	sze int
	buf []byte
	out []byte
	rdr io.Reader
	arr [readerSize]byte
}

// NewReader creates a new Reader which
// reads from an underlying io.Reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{rdr: r}
}

// NewReaderBytes creates a new Reader
// which reads from a byte slice.
func NewReaderBytes(b []byte) *Reader {
	return &Reader{out: b}
}

// Reset resets the Reader, and instructs it
// to read from the specified io.Reader.
func (r *Reader) Reset(i io.Reader) error {
	r.pos = 0
	r.rdr = i
	r.out = nil
	return nil
}

// ResetBytes resets the Reader, and instructs
// it to read from the specified byte slice.
func (r *Reader) ResetBytes(b []byte) error {
	r.pos = 0
	r.out = b
	r.rdr = nil
	return nil
}

// PeekByte returns the next byte in the
// stream without advancing the position
// of the reader.
func (r *Reader) PeekByte() (byte, error) {
	if r.out != nil {
		return r.peekByteFromBytes()
	}
	return r.peekByteFromReader()
}

// ReadByte reads a single byte from the
// underlying io.Reader, or byte slice,
// and advances the position.
func (r *Reader) ReadByte() (byte, error) {
	if r.out != nil {
		return r.readByteFromBytes()
	}
	return r.readByteFromReader()
}

// ReadBytes reads the specified number
// of bytes from the underlying io.Reader,
// or byte slice, and advances the position.
func (r *Reader) ReadBytes(l int) ([]byte, error) {
	if r.out != nil {
		return r.readBytesFromBytes(l)
	}
	return r.readBytesFromReader(l)
}

func (r *Reader) peekByteFromBytes() (byte, error) {

	// Return an error if there is no more data.

	if r.pos+1 > len(r.out) {
		return byte(0), io.EOF
	}

	// Get the next byte without advancing.

	b := r.out[r.pos]

	// Everything went ok.

	return b, nil

}

func (r *Reader) readByteFromBytes() (byte, error) {

	// Return an error if there is no more data.

	if r.pos+1 > len(r.out) {
		return byte(0), io.EOF
	}

	// Get the next byte from the byte slice.

	b := r.out[r.pos]

	// Advance the buffer position.

	r.pos += 1

	// Everything went ok.

	return b, nil

}

func (r *Reader) readBytesFromBytes(l int) ([]byte, error) {

	// Return an error if there is no more data.

	if r.pos+l > len(r.out) {
		return nil, io.EOF
	}

	// Get the data from the byte slice.

	b := r.out[r.pos : r.pos+l]

	// Advance the buffer position.

	r.pos += l

	// Everything went ok.

	return b, nil

}

func (r *Reader) peekByteFromReader() (byte, error) {

	// Initialise the underlying buffer if needed.

	if r.buf == nil {
		r.buf = r.arr[0:]
	}

	// Fill the buffer with data if there is not enough.

	if r.sze == 0 || r.pos+1 > r.sze {
		err := r.fill()
		if err != nil {
			return byte(0), err
		}
	}

	// Get the next byte without advancing.

	b := r.buf[r.pos]

	// Everything went ok.

	return b, nil

}

func (r *Reader) readByteFromReader() (byte, error) {

	// Initialise the underlying buffer if needed.

	if r.buf == nil {
		r.buf = r.arr[0:]
	}

	// Fill the buffer with data if there is not enough.

	if r.sze == 0 || r.pos+1 > r.sze {
		err := r.fill()
		if err != nil {
			return byte(0), err
		}
	}

	// Get the data from the underlying buffer.

	b := r.buf[r.pos]

	// Advance the buffer position.

	r.pos += 1

	// Everything went ok.

	return b, nil

}

func (r *Reader) readBytesFromReader(l int) ([]byte, error) {

	// Initialise the underlying buffer if needed.

	if r.buf == nil {
		r.buf = r.arr[0:]
	}

	// Initialise the byte slice for returning.

	b := make([]byte, l)

	// Loop through until we have filled the byte slice.

	for p := 0; l > 0; {

		// Fill the buffer with data if there is not enough.

		if r.sze == 0 || r.pos+l > r.sze {
			err := r.fill()
			if err != nil {
				return nil, err
			}
		}

		// Get the data from the underlying buffer.

		n := copy(b[p:], r.buf[r.pos:])

		// Advance the buffer position.

		r.pos += n
		p += n
		l -= n

	}

	// Everything went ok.

	return b, nil

}

func (r *Reader) fill() error {
	if r.pos > 0 && r.pos < r.sze {
		copy(r.buf, r.buf[r.pos:r.sze])
		r.sze -= r.pos
	} else {
		r.sze = 0
	}
	n, err := r.rdr.Read(r.buf[r.sze:])
	if err != nil {
		return err
	}
	r.sze += n
	r.pos = 0
	return nil
}
