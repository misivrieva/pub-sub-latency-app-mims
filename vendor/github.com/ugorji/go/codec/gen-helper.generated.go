// comment this out // + build ignore

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from gen-helper.go.tmpl - DO NOT EDIT.

package codec

import "encoding"

// GenVersion is the current version of codecgen.
const GenVersion = 17

// This file is used to generate helper code for codecgen.
// The values here i.e. genHelper(En|De)coder are not to be used directly by
// library users. They WILL change continuously and without notice.

// GenHelperEncoder is exported so that it can be used externally by codecgen.
//
// Library users: DO NOT USE IT DIRECTLY. IT WILL CHANGE CONTINOUSLY WITHOUT NOTICE.
func GenHelperEncoder(e *Encoder) (ge genHelperEncoder, ee genHelperEncDriver) {
	ge = genHelperEncoder{e: e}
	ee = genHelperEncDriver{encDriver: e.e}
	return
}

// GenHelperDecoder is exported so that it can be used externally by codecgen.
//
// Library users: DO NOT USE IT DIRECTLY. IT WILL CHANGE CONTINOUSLY WITHOUT NOTICE.
func GenHelperDecoder(d *Decoder) (gd genHelperDecoder, dd genHelperDecDriver) {
	gd = genHelperDecoder{d: d}
	dd = genHelperDecDriver{decDriver: d.d}
	return
}

type genHelperEncDriver struct {
	encDriver
}

type genHelperDecDriver struct {
	decDriver
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
type genHelperEncoder struct {
	M mustHdl
	F fastpathT
	e *Encoder
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
type genHelperDecoder struct {
	C checkOverflow
	F fastpathT
	d *Decoder
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncBasicHandle() *BasicHandle {
	return f.e.h
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncBinary() bool {
	return f.e.be // f.e.hh.isBinaryEncoding()
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) IsJSONHandle() bool {
	return f.e.js
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncFallback(iv interface{}) {
	// f.e.encodeI(iv, false, false)
	f.e.encodeValue(rv4i(iv), nil)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncTextMarshal(iv encoding.TextMarshaler) {
	bs, fnerr := iv.MarshalText()
	f.e.marshalUtf8(bs, fnerr)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncJSONMarshal(iv jsonMarshaler) {
	bs, fnerr := iv.MarshalJSON()
	f.e.marshalAsis(bs, fnerr)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncBinaryMarshal(iv encoding.BinaryMarshaler) {
	bs, fnerr := iv.MarshalBinary()
	f.e.marshalRaw(bs, fnerr)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncRaw(iv Raw) { f.e.rawBytes(iv) }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) I2Rtid(v interface{}) uintptr {
	return i2rtid(v)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) Extension(rtid uintptr) (xfn *extTypeTagFn) {
	return f.e.h.getExt(rtid, true)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncExtension(v interface{}, xfFn *extTypeTagFn) {
	f.e.e.EncodeExt(v, xfFn.tag, xfFn.ext)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) WriteStr(s string) {
	f.e.w().writestr(s)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) BytesView(v string) []byte { return bytesView(v) }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteMapStart(length int) { f.e.mapStart(length) }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteMapEnd() { f.e.mapEnd() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteArrayStart(length int) { f.e.arrayStart(length) }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteArrayEnd() { f.e.arrayEnd() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteArrayElem() { f.e.arrayElem() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteMapElemKey() { f.e.mapElemKey() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteMapElemValue() { f.e.mapElemValue() }

// ---------------- DECODER FOLLOWS -----------------

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecBasicHandle() *BasicHandle {
	return f.d.h
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecBinary() bool {
	return f.d.be // f.d.hh.isBinaryEncoding()
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecSwallow() { f.d.swallow() }

// // FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
// func (f genHelperDecoder) DecScratchBuffer() []byte {
// 	return f.d.b[:]
// }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecScratchArrayBuffer() *[decScratchByteArrayLen]byte {
	return &f.d.b
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecFallback(iv interface{}, chkPtr bool) {
	rv := rv4i(iv)
	if chkPtr {
		f.d.ensureDecodeable(rv)
	}
	f.d.decodeValue(rv, nil)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecSliceHelperStart() (decSliceHelper, int) {
	return f.d.decSliceHelperStart()
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecStructFieldNotFound(index int, name string) {
	f.d.structFieldNotFound(index, name)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecArrayCannotExpand(sliceLen, streamLen int) {
	f.d.arrayCannotExpand(sliceLen, streamLen)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecTextUnmarshal(tm encoding.TextUnmarshaler) {
	if fnerr := tm.UnmarshalText(f.d.d.DecodeStringAsBytes()); fnerr != nil {
		halt.onerror(fnerr)
	}
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecJSONUnmarshal(tm jsonUnmarshaler) {
	// bs := f.dd.DecodeStringAsBytes()
	// grab the bytes to be read, as UnmarshalJSON needs the full JSON so as to unmarshal it itself.
	bs := f.d.blist.get(256)[:0]
	bs = f.d.d.nextValueBytes(bs)
	fnerr := tm.UnmarshalJSON(bs)
	f.d.blist.put(bs)
	if fnerr != nil {
		halt.onerror(fnerr)
	}
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecBinaryUnmarshal(bm encoding.BinaryUnmarshaler) {
	if fnerr := bm.UnmarshalBinary(f.d.d.DecodeBytes(nil, true)); fnerr != nil {
		halt.onerror(fnerr)
	}
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecRaw() []byte { return f.d.rawBytes() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) IsJSONHandle() bool {
	return f.d.js
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) I2Rtid(v interface{}) uintptr {
	return i2rtid(v)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) Extension(rtid uintptr) (xfn *extTypeTagFn) {
	return f.d.h.getExt(rtid, true)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecExtension(v interface{}, xfFn *extTypeTagFn) {
	f.d.d.DecodeExt(v, xfFn.tag, xfFn.ext)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecInferLen(clen, maxlen, unit int) (rvlen int) {
	return decInferLen(clen, maxlen, unit)
}

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) StringView(v []byte) string { return stringView(v) }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadMapStart() int { return f.d.mapStart() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadMapEnd() { f.d.mapEnd() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadArrayStart() int { return f.d.arrayStart() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadArrayEnd() { f.d.arrayEnd() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadArrayElem() { f.d.arrayElem() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadMapElemKey() { f.d.mapElemKey() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadMapElemValue() { f.d.mapElemValue() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecDecodeFloat32() float32 { return f.d.decodeFloat32() }

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecCheckBreak() bool { return f.d.checkBreak() }
