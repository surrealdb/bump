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
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var run = 100
var txt []byte // 906
var jpg []byte // 9245562

func init() {

	txt, _ = ioutil.ReadFile("test.txt")
	jpg, _ = ioutil.ReadFile("test.jpg")

}

func chunkRead(r *Reader, num int) (bit []byte) {
	sze := readerSize
	beg := int(num / 100)
	for i := 0; i < beg; i++ {
		p, _ := r.PeekByte()
		b, _ := r.ReadByte()
		if p != b {
			panic("Arrrgh")
		}
		bit = append(bit, b)
	}
	for i := beg; i < num; i += sze {
		if i+sze > num {
			b, _ := r.ReadBytes(num - i)
			bit = append(bit, b...)
		} else {
			b, _ := r.ReadBytes(sze)
			bit = append(bit, b...)
		}
	}
	var e error
	_, e = r.PeekByte()
	So(e, ShouldNotBeNil)
	_, e = r.ReadByte()
	So(e, ShouldNotBeNil)
	_, e = r.ReadBytes(5)
	So(e, ShouldNotBeNil)
	return
}

func chunkWrite(w *Writer, bit []byte) {
	sze := writerSize
	beg := int(len(bit) / 100)
	for i := 0; i < beg; i++ {
		w.WriteByte(bit[i])
	}
	for i := beg; i < len(bit); i += sze {
		if i+sze > len(bit) {
			w.WriteBytes(bit[i:])
		} else {
			w.WriteBytes(bit[i : i+sze])
		}
	}
}

func TestReader(t *testing.T) {

	Convey("Reader should read small data from an io.Reader", t, func() {
		b := bytes.NewReader(txt)
		r := NewReader(b)
		o, _ := r.ReadBytes(len(txt))
		So(len(o), ShouldEqual, len(txt))
		So(o[0:5], ShouldResemble, txt[0:5])
		So(o[len(o)-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Reader should read small data from an io.Reader (reading in chunks)", t, func() {
		b := bytes.NewReader(txt)
		r := NewReader(b)
		o := chunkRead(r, len(txt))
		So(len(o), ShouldEqual, len(txt))
		So(o[0:5], ShouldResemble, txt[0:5])
		So(o[len(o)-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Reader should read large data from an io.Reader", t, func() {
		b := bytes.NewReader(jpg)
		r := NewReader(b)
		o, _ := r.ReadBytes(len(jpg))
		So(len(o), ShouldEqual, len(jpg))
		So(o[0:5], ShouldResemble, jpg[0:5])
		So(o[len(o)-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	Convey("Reader should read large data from an io.Reader (reading in chunks)", t, func() {
		b := bytes.NewReader(jpg)
		r := NewReader(b)
		o := chunkRead(r, len(jpg))
		So(len(o), ShouldEqual, len(jpg))
		So(o[0:5], ShouldResemble, jpg[0:5])
		So(o[len(o)-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	// ----------------------------------------------------------------------------------------------------

	Convey("Reader should read small data from a byte slice", t, func() {
		r := NewReaderBytes(txt)
		o, _ := r.ReadBytes(len(txt))
		So(len(o), ShouldEqual, len(txt))
		So(o[0:5], ShouldResemble, txt[0:5])
		So(o[len(o)-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Reader should read small data from a byte slice (reading in chunks)", t, func() {
		r := NewReaderBytes(txt)
		o := chunkRead(r, len(txt))
		So(len(o), ShouldEqual, len(txt))
		So(o[0:5], ShouldResemble, txt[0:5])
		So(o[len(o)-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Reader should read large data from a byte slice", t, func() {
		r := NewReaderBytes(jpg)
		o, _ := r.ReadBytes(len(jpg))
		So(len(o), ShouldEqual, len(jpg))
		So(o[0:5], ShouldResemble, jpg[0:5])
		So(o[len(o)-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	Convey("Reader should read large data from a byte slice (reading in chunks)", t, func() {
		r := NewReaderBytes(jpg)
		o := chunkRead(r, len(jpg))
		So(len(o), ShouldEqual, len(jpg))
		So(o[0:5], ShouldResemble, jpg[0:5])
		So(o[len(o)-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	// ----------------------------------------------------------------------------------------------------

	Convey("Reader should read large data from an io.Reader before resetting and being used for further reads", t, func() {
		b := bytes.NewReader(jpg)
		r := NewReader(b)
		for i := 0; i < run; i++ {
			o, _ := r.ReadBytes(len(jpg))
			So(len(o), ShouldEqual, len(jpg))
			So(o[0:5], ShouldResemble, jpg[0:5])
			So(o[len(o)-5:], ShouldResemble, jpg[len(jpg)-5:])
			b = bytes.NewReader(jpg)
			r.Reset(b)
		}
	})

	Convey("Reader should read large data from a byte slice before resetting and being used for further reads", t, func() {
		r := NewReaderBytes(jpg)
		for i := 0; i < run; i++ {
			o, _ := r.ReadBytes(len(jpg))
			So(len(o), ShouldEqual, len(jpg))
			So(o[0:5], ShouldResemble, jpg[0:5])
			So(o[len(o)-5:], ShouldResemble, jpg[len(jpg)-5:])
			r.ResetBytes(jpg)
		}
	})

}

func TestWriter(t *testing.T) {

	Convey("Writer should write small data to an io.Writer", t, func() {
		b := bytes.NewBuffer(nil)
		w := NewWriter(b)
		w.WriteBytes(txt)
		w.Flush()
		So(b.Len(), ShouldEqual, len(txt))
		So(b.Bytes()[0:5], ShouldResemble, txt[0:5])
		So(b.Bytes()[b.Len()-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Writer should write small data to an io.Writer (writing in chunks)", t, func() {
		b := bytes.NewBuffer(nil)
		w := NewWriter(b)
		chunkWrite(w, txt)
		w.Flush()
		So(b.Len(), ShouldEqual, len(txt))
		So(b.Bytes()[0:5], ShouldResemble, txt[0:5])
		So(b.Bytes()[b.Len()-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Writer should write large data to an io.Writer", t, func() {
		b := bytes.NewBuffer(nil)
		w := NewWriter(b)
		w.WriteBytes(jpg)
		w.Flush()
		So(b.Len(), ShouldEqual, len(jpg))
		So(b.Bytes()[0:5], ShouldResemble, jpg[0:5])
		So(b.Bytes()[b.Len()-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	Convey("Writer should write large data to an io.Writer (writing in chunks)", t, func() {
		b := bytes.NewBuffer(nil)
		w := NewWriter(b)
		chunkWrite(w, jpg)
		w.Flush()
		So(b.Len(), ShouldEqual, len(jpg))
		So(b.Bytes()[0:5], ShouldResemble, jpg[0:5])
		So(b.Bytes()[b.Len()-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	// ----------------------------------------------------------------------------------------------------

	Convey("Writer should write small data to a nil byte slice and allocate once", t, func() {
		var b []byte
		w := NewWriterBytes(&b)
		w.WriteBytes(txt)
		w.Flush()
		So(len(b), ShouldEqual, len(txt))
		So(b[0:5], ShouldResemble, txt[0:5])
		So(b[len(b)-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Writer should write small data to a nil byte slice and allocate once (writing in chunks)", t, func() {
		var b []byte
		w := NewWriterBytes(&b)
		chunkWrite(w, txt)
		w.Flush()
		So(len(b), ShouldEqual, len(txt))
		So(b[0:5], ShouldResemble, txt[0:5])
		So(b[len(b)-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Writer should write large data to a nil byte slice and allocate multiple", t, func() {
		var b []byte
		w := NewWriterBytes(&b)
		w.WriteBytes(jpg)
		w.Flush()
		So(len(b), ShouldEqual, len(jpg))
		So(b[0:5], ShouldResemble, jpg[0:5])
		So(b[len(b)-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	Convey("Writer should write large data to a nil byte slice and allocate multiple (writing in chunks)", t, func() {
		var b []byte
		w := NewWriterBytes(&b)
		w.WriteBytes(jpg)
		w.Flush()
		So(len(b), ShouldEqual, len(jpg))
		So(b[0:5], ShouldResemble, jpg[0:5])
		So(b[len(b)-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	Convey("Writer should write small data to an empty byte slice and grow to capacity", t, func() {
		b := make([]byte, 0, 1024)
		w := NewWriterBytes(&b)
		w.WriteBytes(txt)
		w.Flush()
		So(len(b), ShouldEqual, len(txt))
		So(b[0:5], ShouldResemble, txt[0:5])
		So(b[len(b)-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Writer should write small data to an empty byte slice and grow to capacity (writing in chunks)", t, func() {
		b := make([]byte, 0, 1024)
		w := NewWriterBytes(&b)
		w.WriteBytes(txt)
		w.Flush()
		So(len(b), ShouldEqual, len(txt))
		So(b[0:5], ShouldResemble, txt[0:5])
		So(b[len(b)-5:], ShouldResemble, txt[len(txt)-5:])
	})

	Convey("Writer should write large data to an empty byte slice and grow to capacity and allocate ", t, func() {
		b := make([]byte, 0, 1024)
		w := NewWriterBytes(&b)
		w.WriteBytes(jpg)
		w.Flush()
		So(len(b), ShouldEqual, len(jpg))
		So(b[0:5], ShouldResemble, jpg[0:5])
		So(b[len(b)-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	Convey("Writer should write large data to an empty byte slice and grow to capacity and allocate (writing in chunks)", t, func() {
		b := make([]byte, 0, 1024)
		w := NewWriterBytes(&b)
		w.WriteBytes(jpg)
		w.Flush()
		So(len(b), ShouldEqual, len(jpg))
		So(b[0:5], ShouldResemble, jpg[0:5])
		So(b[len(b)-5:], ShouldResemble, jpg[len(jpg)-5:])
	})

	// ----------------------------------------------------------------------------------------------------

	Convey("Writer should write large data to an io.Writer before resetting and being used for further writes", t, func() {
		b := bytes.NewBuffer(nil)
		w := NewWriter(b)
		for i := 0; i < run; i++ {
			w.WriteBytes(jpg)
			w.Flush()
			So(b.Len(), ShouldEqual, len(jpg))
			So(b.Bytes()[0:5], ShouldResemble, jpg[0:5])
			So(b.Bytes()[b.Len()-5:], ShouldResemble, jpg[len(jpg)-5:])
			b = bytes.NewBuffer(nil)
			w.Reset(b)
		}
	})

	Convey("Writer should write large data to a byte slice before resetting and being used for further writes", t, func() {
		b := []byte{}
		w := NewWriterBytes(&b)
		for i := 0; i < run; i++ {
			w.WriteBytes(jpg)
			w.Flush()
			So(len(b), ShouldEqual, len(jpg))
			So(b[0:5], ShouldResemble, jpg[0:5])
			So(b[len(b)-5:], ShouldResemble, jpg[len(jpg)-5:])
			b = []byte{}
			w.ResetBytes(&b)
		}
	})

}
