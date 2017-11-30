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

const writerSize = 1024

// Writer represents a buffer for writing
// to an io.Writer, or a byte slice.
type Writer struct {
	pos int
	buf []byte
	out *[]byte
	wtr io.Writer
	arr [writerSize]byte
}

// stringer represents an io.Writer which
// can write a string without conversion.

type stringer interface {
	WriteString(string) (int, error)
}

// NewWriter creates a new Writer which
// writes to an underlying io.Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{wtr: w}
}

// NewWriterBytes creates a new Writer
// which writes to a byte slice.
func NewWriterBytes(b *[]byte) *Writer {
	return &Writer{out: b}
}

// Reset resets the Writer, and instructs it
// to write to the specified io.Writer.
func (w *Writer) Reset(i io.Writer) error {
	w.pos = 0
	w.wtr = i
	w.out = nil
	return nil
}

// ResetBytes resets the Writer, and instructs
// it to write to the specified byte slice.
func (w *Writer) ResetBytes(b *[]byte) error {
	w.pos = 0
	w.out = b
	w.wtr = nil
	return nil
}

// Flush flushes any remaining buffered data
// to the underlying io.Writer. When writing
// to a byte slice, this function does not
// do anything, as data is written immediately.
func (w *Writer) Flush() error {

	// Don't flush if there is no data.

	if w.pos == 0 {
		return nil
	}

	// Don't flush if we are writing to a slice.

	if w.out != nil {
		return nil
	}

	// Write the data to the underlying writer.

	n, err := w.wtr.Write(w.buf[:w.pos])
	if err != nil {
		return err
	}

	// If not all data was sent, then error.

	if n < w.pos {
		return io.ErrShortWrite
	}

	// Reset the buffer position.

	w.pos = 0

	// Everything went ok.

	return nil

}

// WriteByte writes a single byte to the
// underlying io.Writer, or byte slice.
func (w *Writer) WriteByte(v byte) error {
	if w.out != nil {
		return w.writeByteToBytes(v)
	}
	return w.writeByteToWriter(v)
}

// WriteBytes writes a slice of bytes to the
// underlying io.Writer, or byte slice.
func (w *Writer) WriteBytes(v []byte) error {
	if w.out != nil {
		return w.writeBytesToBytes(v)
	}
	return w.writeBytesToWriter(v)
}

// WriteString writes a string to the
// underlying io.Writer, or byte slice.
func (w *Writer) WriteString(v string) error {
	if w.out != nil {
		return w.writeStringToBytes(v)
	}
	return w.writeStringToWriter(v)
}

func (w *Writer) writeByteToBytes(v byte) error {

	// Grow the underlying buffer if needed.

	if w.pos+1 >= len(*w.out) {
		if w.pos+1 < cap(*w.out) {
			*w.out = (*w.out)[:w.pos+1]
		} else {
			bs := make([]byte, len(*w.out)+1, len(*w.out)+1+writerSize)
			copy(bs, (*w.out)[:w.pos])
			*w.out = bs
		}
	}

	// Insert the specified byte into the buffer.

	(*w.out)[w.pos] = v

	// Increment the current buffer position.

	w.pos += 1

	// Trim the slice to the correct length.

	*w.out = (*w.out)[:w.pos]

	// Everything went ok.

	return nil

}

func (w *Writer) writeBytesToBytes(v []byte) error {

	// Grow the underlying buffer if needed.

	if w.pos+len(v) >= len(*w.out) {
		if w.pos+len(v) < cap(*w.out) {
			*w.out = (*w.out)[:w.pos+len(v)]
		} else {
			bs := make([]byte, len(*w.out)+len(v), len(*w.out)+len(v)+writerSize)
			copy(bs, (*w.out)[:w.pos])
			*w.out = bs
		}
	}

	// Insert the specified bytes into the buffer.

	n := copy((*w.out)[w.pos:], v)

	// Increment the current buffer position.

	w.pos += n

	// Trim the slice to the correct length.

	*w.out = (*w.out)[:w.pos]

	// Everything went ok.

	return nil

}

func (w *Writer) writeStringToBytes(v string) error {

	// Grow the underlying buffer if needed.

	if w.pos+len(v) >= len(*w.out) {
		if w.pos+len(v) < cap(*w.out) {
			*w.out = (*w.out)[:w.pos+len(v)]
		} else {
			bs := make([]byte, len(*w.out)+len(v), len(*w.out)+len(v)+writerSize)
			copy(bs, (*w.out)[:w.pos])
			*w.out = bs
		}
	}

	// Insert the specified bytes into the buffer.

	n := copy((*w.out)[w.pos:], v)

	// Increment the current buffer position.

	w.pos += n

	// Trim the slice to the correct length.

	*w.out = (*w.out)[:w.pos]

	// Everything went ok.

	return nil

}

func (w *Writer) writeByteToWriter(v byte) error {

	// Initialise the underlying buffer if needed.

	if w.buf == nil {
		w.buf = w.arr[0:]
	}

	// Flush the buffer if no space is remaining.

	if w.pos >= len(w.buf) {
		err := w.Flush()
		if err != nil {
			return err
		}
	}

	// Insert the specified byte into the buffer.

	w.buf[w.pos] = v

	// Increment the current buffer position.

	w.pos += 1

	// Everything went ok.

	return nil

}

func (w *Writer) writeBytesToWriter(v []byte) error {

	// Initialise the underlying buffer if needed.

	if w.buf == nil {
		w.buf = w.arr[0:]
	}

	// Write directly to the underlying writer.

	for len(v) > len(w.buf)-w.pos {
		if w.pos == 0 {
			n, err := w.wtr.Write(v)
			if err != nil {
				return err
			}
			v = v[n:]
		} else {
			n := copy(w.buf[w.pos:], v)
			w.pos += n
			v = v[n:]
		}
		if w.pos >= writerSize {
			err := w.Flush()
			if err != nil {
				return err
			}
		}
	}

	// Insert any remaining bytes into the buffer.

	n := copy(w.buf[w.pos:], v)

	// Increment the current buffer position.

	w.pos += n

	// Everything went ok.

	return nil

}

func (w *Writer) writeStringToWriter(s string) error {

	// Attempt to write the string directly.

	if i, ok := w.wtr.(stringer); ok {
		return w.writeStringToStringer(i, s)
	}

	// Convert the string to a slice of bytes.

	v := []byte(s)

	// Initialise the underlying buffer if needed.

	if w.buf == nil {
		w.buf = w.arr[0:]
	}

	// Write directly to the underlying writer.

	for len(v) > len(w.buf)-w.pos {
		if w.pos == 0 {
			n, err := w.wtr.Write(v)
			if err != nil {
				return err
			}
			v = v[n:]
		} else {
			n := copy(w.buf[w.pos:], v)
			w.pos += n
			v = v[n:]
		}
		if w.pos >= writerSize {
			err := w.Flush()
			if err != nil {
				return err
			}
		}
	}

	// Insert any remaining bytes into the buffer.

	n := copy(w.buf[w.pos:], v)

	// Increment the current buffer position.

	w.pos += n

	// Everything went ok.

	return nil

}

func (w *Writer) writeStringToStringer(i stringer, s string) error {

	// Initialise the underlying buffer if needed.

	if w.buf == nil {
		w.buf = w.arr[0:]
	}

	// Write directly to the underlying writer.

	for len(s) > len(w.buf)-w.pos {
		if w.pos == 0 {
			n, err := i.WriteString(s)
			if err != nil {
				return err
			}
			s = s[n:]
		} else {
			n := copy(w.buf[w.pos:], s)
			w.pos += n
			s = s[n:]
		}
		if w.pos >= writerSize {
			err := w.Flush()
			if err != nil {
				return err
			}
		}
	}

	// Insert any remaining bytes into the buffer.

	n := copy(w.buf[w.pos:], s)

	// Increment the current buffer position.

	w.pos += n

	// Everything went ok.

	return nil

}
