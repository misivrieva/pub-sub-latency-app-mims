// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "io"

// encWriter abstracts writing to a byte array or to an io.Writer.
type encWriter interface {
	writeb([]byte)
	writestr(string)
	writeqstr(string) // write string wrapped in quotes ie "..."
	writen1(byte)

	// add convenience functions for writing 2,4
	writen2(byte, byte)
	writen4(byte, byte, byte, byte)

	end()
}

// ---------------------------------------------

// bufioEncWriter
type bufioEncWriter struct {
	w io.Writer

	buf []byte

	n int

	b [16]byte // scratch buffer and padding (cache-aligned)
}

func (z *bufioEncWriter) reset(w io.Writer, bufsize int, blist *bytesFreelist) {
	z.w = w
	z.n = 0
	if bufsize <= 0 {
		bufsize = defEncByteBufSize
	}
	// bufsize must be >= 8, to accomodate writen methods (where n <= 8)
	if bufsize <= 8 {
		bufsize = 8
	}
	if cap(z.buf) < bufsize {
		if len(z.buf) > 0 && &z.buf[0] != &z.b[0] {
			blist.put(z.buf)
		}
		if len(z.b) > bufsize {
			z.buf = z.b[:]
		} else {
			z.buf = blist.get(bufsize)
		}
	}
	z.buf = z.buf[:cap(z.buf)]
}

func (z *bufioEncWriter) flushErr() (err error) {
	n, err := z.w.Write(z.buf[:z.n])
	z.n -= n
	if z.n > 0 && err == nil {
		err = io.ErrShortWrite
	}
	if n > 0 && z.n > 0 {
		copy(z.buf, z.buf[n:z.n+n])
	}
	return err
}

func (z *bufioEncWriter) flush() {
	halt.onerror(z.flushErr())
}

func (z *bufioEncWriter) writeb(s []byte) {
LOOP:
	a := len(z.buf) - z.n
	if len(s) > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
}

func (z *bufioEncWriter) writestr(s string) {
	// z.writeb(bytesView(s)) // inlined below
LOOP:
	a := len(z.buf) - z.n
	if len(s) > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
}

func (z *bufioEncWriter) writeqstr(s string) {
	// z.writen1('"')
	// z.writestr(s)
	// z.writen1('"')

	if z.n+len(s)+2 > len(z.buf) {
		z.flush()
	}
	z.buf[z.n] = '"'
	z.n++
LOOP:
	a := len(z.buf) - z.n
	if len(s)+1 > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
	z.buf[z.n] = '"'
	z.n++
}

func (z *bufioEncWriter) writen1(b1 byte) {
	if 1 > len(z.buf)-z.n {
		z.flush()
	}
	z.buf[z.n] = b1
	z.n++
}

func (z *bufioEncWriter) writen2(b1, b2 byte) {
	if 2 > len(z.buf)-z.n {
		z.flush()
	}
	z.buf[z.n+1] = b2
	z.buf[z.n] = b1
	z.n += 2
}

func (z *bufioEncWriter) writen4(b1, b2, b3, b4 byte) {
	if 4 > len(z.buf)-z.n {
		z.flush()
	}
	z.buf[z.n+3] = b4
	z.buf[z.n+2] = b3
	z.buf[z.n+1] = b2
	z.buf[z.n] = b1
	z.n += 4
}

func (z *bufioEncWriter) endErr() (err error) {
	if z.n > 0 {
		err = z.flushErr()
	}
	return
}

// ---------------------------------------------

// bytesEncAppender implements encWriter and can write to an byte slice.
type bytesEncAppender struct {
	b   []byte
	out *[]byte
}

func (z *bytesEncAppender) writeb(s []byte) {
	z.b = append(z.b, s...)
}
func (z *bytesEncAppender) writestr(s string) {
	z.b = append(z.b, s...)
}
func (z *bytesEncAppender) writeqstr(s string) {
	z.b = append(append(append(z.b, '"'), s...), '"')

	// z.b = append(z.b, '"')
	// z.b = append(z.b, s...)
	// z.b = append(z.b, '"')
}
func (z *bytesEncAppender) writen1(b1 byte) {
	z.b = append(z.b, b1)
}
func (z *bytesEncAppender) writen2(b1, b2 byte) {
	z.b = append(z.b, b1, b2)
}
func (z *bytesEncAppender) writen4(b1, b2, b3, b4 byte) {
	z.b = append(z.b, b1, b2, b3, b4)
}
func (z *bytesEncAppender) endErr() error {
	*(z.out) = z.b
	return nil
}
func (z *bytesEncAppender) reset(in []byte, out *[]byte) {
	z.b = in[:0]
	z.out = out
}

// --------------------------------------------------

type encWr struct {
	bytes bool // encoding to []byte
	js    bool // is json encoder?
	be    bool // is binary encoder?

	c containerState

	calls uint16

	wb bytesEncAppender
	wf *bufioEncWriter
}

func (z *encWr) writeb(s []byte) {
	if z.bytes {
		z.wb.writeb(s)
	} else {
		z.wf.writeb(s)
	}
}
func (z *encWr) writeqstr(s string) {
	if z.bytes {
		// MARKER: manually inline, else this function is not inlined.
		// Keep in sync with bytesEncWriter.writeqstr
		// z.wb.writeqstr(s)
		z.wb.b = append(append(append(z.wb.b, '"'), s...), '"')
	} else {
		z.wf.writeqstr(s)
	}
}
func (z *encWr) writestr(s string) {
	if z.bytes {
		z.wb.writestr(s)
	} else {
		z.wf.writestr(s)
	}
}
func (z *encWr) writen1(b1 byte) {
	if z.bytes {
		z.wb.writen1(b1)
	} else {
		z.wf.writen1(b1)
	}
}

// MARKER: manually inline bytesEncAppender.writenx methods,
// as calling them causes encWr.writenx methods to not be inlined.
//
// i.e. instead of writing z.wb.writen2(b1, b2), use z.wb.b = append(z.wb.b, b1, b2)

func (z *encWr) writen2(b1, b2 byte) {
	if z.bytes {
		z.wb.b = append(z.wb.b, b1, b2)
	} else {
		z.wf.writen2(b1, b2)
	}
}
func (z *encWr) writen4(b1, b2, b3, b4 byte) {
	if z.bytes {
		z.wb.b = append(z.wb.b, b1, b2, b3, b4)
	} else {
		z.wf.writen4(b1, b2, b3, b4)
	}
}

func (z *encWr) endErr() error {
	if z.bytes {
		return z.wb.endErr()
	}
	return z.wf.endErr()
}

func (z *encWr) end() {
	halt.onerror(z.endErr())
}

var _ encWriter = (*encWr)(nil)
