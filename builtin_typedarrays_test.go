package goja

import (
	"bytes"
	"testing"
)

/*
func TestArrayBufferNew(t *testing.T) {
	const SCRIPT = `
	var b = new ArrayBuffer(16);
	b.byteLength;
	`

	testScript(SCRIPT, intToValue(16), t)
}
*/

func TestArrayBufferSetUint32(t *testing.T) {
	vm := New()
	b := vm._newArrayBuffer(vm.global.ArrayBufferPrototype, nil)
	b.data = make([]byte, 4)
	b.setUint32(0, 0xCAFEBABE, bigEndian)

	i := b.getUint32(0, bigEndian)
	if i != 0xCAFEBABE {
		t.Fatal(i)
	}
	i = b.getUint32(0, littleEndian)
	if i != 0xBEBAFECA {
		t.Fatal(i)
	}

	b.setUint32(0, 0xBEBAFECA, littleEndian)
	i = b.getUint32(0, bigEndian)
	if i != 0xCAFEBABE {
		t.Fatal(i)
	}
}

func TestArrayBufferSetInt32(t *testing.T) {
	vm := New()
	b := vm._newArrayBuffer(vm.global.ArrayBufferPrototype, nil)
	b.data = make([]byte, 4)
	b.setInt32(0, -42, littleEndian)
	if v := b.getInt32(0, littleEndian); v != -42 {
		t.Fatal(v)
	}

	b.setInt32(0, -42, bigEndian)
	if v := b.getInt32(0, bigEndian); v != -42 {
		t.Fatal(v)
	}
}

func TestNewUint8Array(t *testing.T) {
	const SCRIPT = `
	var a = new Uint8Array(1);
	a[0] = 42;
	a.byteLength === 1 && a.length === 1 && a[0] === 42;
	`

	testScript(SCRIPT, valueTrue, t)
}

func TestNewUint16Array(t *testing.T) {
	const SCRIPT = `
	var a = new Uint16Array(1);
	a[0] = 42;
	a.byteLength === 2 && a.length === 1 && a[0] === 42;
	`

	testScript(SCRIPT, valueTrue, t)
}

func TestTypedArraysSpeciesConstructor(t *testing.T) {
	const SCRIPT = `
    'use strict';
    function MyArray() {
        var NewTarget = this.__proto__.constructor;
        return Reflect.construct(Uint16Array, arguments, NewTarget);
    }
    MyArray.prototype = Object.create(Uint16Array.prototype, {
        constructor: {
            value: MyArray,
            writable: true,
            configurable: true
        }
    });
    var a = new MyArray(1);
    Object.defineProperty(MyArray, Symbol.species, {value: Uint8Array, configurable: true});
    a[0] = 32767;
    var b = a.filter(function() {
        return true;
    });
	if (a[0] !== 32767) {
		throw new Error("a[0]=" + a[0]); 
	}
	if (!(b instanceof Uint8Array)) {
		throw new Error("b instanceof Uint8Array");
	}
	if (b[0] != 255) {
		throw new Error("b[0]=" + b[0]);
	}
	`

	testScript(SCRIPT, _undefined, t)
}

func TestTypedArrayFromArrayBuffer(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(2);
	var a16 = new Uint16Array(buf);
	if (!(a16 instanceof Uint16Array)) {
		throw new Error("a16 is not an instance");
	}
	if (a16.buffer !== buf) {
		throw new Error("a16.buffer !== buf");
	}
	if (a16.length !== 1) {
		throw new Error("a16.length=" + a16.length);
	}
	var a8 = new Uint8Array(buf);
	a8.fill(0xAA);
	if (a16[0] !== 0xAAAA) {
		throw new Error("a16[0]=" + a16[0]);
	}
	`

	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySetOverlapDifSize(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(4);
	var src = new Uint8Array(buf, 1, 2);
	src[0] = 1;
	src[1] = 2;
	var dst = new Uint16Array(buf);
	dst.set(src);
	if (dst[0] !== 1 || dst[1] !== 2) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySetOverlapDifSize2(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(4);
	var src = new Uint8Array(buf, 0, 2);
	src[0] = 1;
	src[1] = 2;
	var dst = new Uint16Array(buf);
	dst.set(src);
	if (dst[0] !== 1 || dst[1] !== 2) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySetOverlapDifSize3(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(8);
	var src = new Uint8Array(buf, 2, 4);
	src[0] = 1;
	src[1] = 2;
	src[2] = 3;
	src[3] = 4;
	var dst = new Uint16Array(buf);
	dst.set(src);
	if (dst[0] !== 1 || dst[1] !== 2 || dst[2] !== 3 || dst[3] !== 4) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySetOverlapDifSize4(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(10);
	var dst = new Uint8Array(buf, 2, 5);
	var src = new Uint16Array(buf);
	src[0] = 1;
	src[1] = 2;
	src[2] = 3;
	src[3] = 4;
	src[4] = 5;
	dst.set(src);
	if (dst[0] !== 1 || dst[1] !== 2 || dst[2] !== 3 || dst[3] !== 4 || dst[4] !== 5) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySetNoOverlapDifSizeForward(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(10);
	var dst = new Uint8Array(buf, 7, 2);
	var src = new Uint16Array(buf, 0, 2);
	src[0] = 1;
	src[1] = 2;
	dst.set(src);
	if (dst[0] !== 1 || dst[1] !== 2 || src[0] !== 1 || src[1] !== 2) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySetNoOverlapDifSizeBackward(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(10);
	var dst = new Uint8Array(buf, 0, 2);
	var src = new Uint16Array(buf, 6, 2);
	src[0] = 1;
	src[1] = 2;
	dst.set(src);
	if (dst[0] !== 1 || dst[1] !== 2 || src[0] !== 1 || src[1] !== 2) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySetNoOverlapDifSizeDifBuffers(t *testing.T) {
	const SCRIPT = `
	var dstBuf = new ArrayBuffer(1024);
	var dst = new Uint8Array(dstBuf, 0, 2);
	var src = new Uint16Array(2);
	src[0] = 1;
	src[1] = 2;
	dst.set(src);
	if (dst[0] !== 1 || dst[1] !== 2 || src[0] !== 1 || src[1] !== 2) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySliceSameType(t *testing.T) {
	const SCRIPT = `
	var src = Uint8Array.of(1,2,3,4);
	var dst = src.slice(1, 3);
	if (dst.length !== 2 || dst[0] !== 2 || dst[1] !== 3) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySliceDifType(t *testing.T) {
	const SCRIPT = `
	var src = Uint8Array.of(1,2,3,4);
	Object.defineProperty(Uint8Array, Symbol.species, {value: Uint16Array, configurable: true});
	var dst = src.slice(1, 3);
	if (!(dst instanceof Uint16Array)) {
		throw new Error("wrong dst type: " + dst);
	}
	if (dst.length !== 2 || dst[0] !== 2 || dst[1] !== 3) {
		throw new Error("dst: " + dst.join(","));
	}	
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySortComparatorReturnValueFloats(t *testing.T) {
	const SCRIPT = `
	var a = Float64Array.of(
		5.97,
		9.91,
		4.13,
		9.28,
		3.29
	);
	a.sort( function(a, b) { return a - b; } );
	for (var i = 1; i < a.length; i++) {
		if (a[i] < a[i-1]) {
			throw new Error("Array is not sorted: " + a);
		}
	}
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySortComparatorReturnValueNegZero(t *testing.T) {
	const SCRIPT = `
	var a = new Uint8Array([2, 1]);
	a.sort( function(a, b) { return a > b ? 0 : -0; } );
	for (var i = 1; i < a.length; i++) {
		if (a[i] < a[i-1]) {
			throw new Error("Array is not sorted: " + a);
		}
	}
	`
	testScript(SCRIPT, _undefined, t)
}

func TestInt32ArrayNegativeIndex(t *testing.T) {
	const SCRIPT = `
	new Int32Array()[-1] === undefined;
	`

	testScript(SCRIPT, valueTrue, t)
}

func TestTypedArrayDeleteUnconfigurable(t *testing.T) {
	const SCRIPT = `
	try {
		(function() {
			'use strict';
			delete Uint8Array.prototype.BYTES_PER_ELEMENT;
		})();
	} catch(e) {
		if (!(e instanceof TypeError)) {
			throw e;
		}
		if (!e.message.startsWith("Cannot delete property")) {
			throw e;
		}
	}
	`

	testScript(SCRIPT, _undefined, t)
}

func TestTypedArrayExportToBytes(t *testing.T) {
	const SCRIPT = `
	const arr = new Uint16Array(10);
	arr[2] = 0xFFFF; // make sure it's not affected by the platform endianness
	arr.subarray(2, 4);
	`
	vm := New()
	v, err := vm.RunString(SCRIPT)
	if err != nil {
		t.Fatal(err)
	}
	var b []byte
	err = vm.ExportTo(v, &b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{0xFF, 0xFF, 0, 0}) {
		t.Fatal("unexpected value", b)
	}
}

func TestDataViewExportToBytes(t *testing.T) {
	const SCRIPT = `
	const arr = new Uint16Array(10);
	arr[2] = 0xFFFF;
	new DataView(arr.buffer, 4, 4);
	`
	vm := New()
	v, err := vm.RunString(SCRIPT)
	if err != nil {
		t.Fatal(err)
	}
	var b []byte
	err = vm.ExportTo(v, &b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{0xFF, 0xFF, 0, 0}) {
		t.Fatal("unexpected value", b)
	}
}

// The following tests guard against a regression where TypedArray.prototype.with
// and toReversed wrote to the freshly allocated destination using the source
// view's byte offset (ta.offset+k) instead of the plain index k. For views with
// a non-zero byte offset this produced wrong results and an out-of-bounds write
// past the destination buffer. See dop251/goja#712.

func TestTypedArrayWithNonZeroOffset(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(8);
	var a = new Uint8Array(buf, 4, 4); // byteOffset 4, length 4
	a[0] = 10; a[1] = 20; a[2] = 30; a[3] = 40;
	var b = a.with(1, 99);
	if (b.length !== 4) {
		throw new Error("b.length=" + b.length);
	}
	if (b.byteOffset !== 0) {
		throw new Error("b.byteOffset=" + b.byteOffset);
	}
	if (b[0] !== 10 || b[1] !== 99 || b[2] !== 30 || b[3] !== 40) {
		throw new Error("b=[" + b.join(",") + "], expected [10,99,30,40]");
	}
	// source must be unchanged
	if (a[0] !== 10 || a[1] !== 20 || a[2] !== 30 || a[3] !== 40) {
		throw new Error("a corrupted=[" + a.join(",") + "]");
	}
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArrayWithNonZeroOffsetNegativeIndex(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(8);
	var a = new Uint8Array(buf, 4, 4);
	a[0] = 10; a[1] = 20; a[2] = 30; a[3] = 40;
	var b = a.with(-1, 99); // index 3
	if (b[0] !== 10 || b[1] !== 20 || b[2] !== 30 || b[3] !== 99) {
		throw new Error("b=[" + b.join(",") + "], expected [10,20,30,99]");
	}
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArrayToReversedNonZeroOffset(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(8);
	var a = new Uint8Array(buf, 4, 4);
	a[0] = 1; a[1] = 2; a[2] = 3; a[3] = 4;
	var b = a.toReversed();
	if (b.length !== 4) {
		throw new Error("b.length=" + b.length);
	}
	if (b.byteOffset !== 0) {
		throw new Error("b.byteOffset=" + b.byteOffset);
	}
	if (b[0] !== 4 || b[1] !== 3 || b[2] !== 2 || b[3] !== 1) {
		throw new Error("b=[" + b.join(",") + "], expected [4,3,2,1]");
	}
	if (a[0] !== 1 || a[1] !== 2 || a[2] !== 3 || a[3] !== 4) {
		throw new Error("a corrupted=[" + a.join(",") + "]");
	}
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArrayWithToReversedNonZeroOffsetAllTypes(t *testing.T) {
	const SCRIPT = `
	var types = [Int8Array, Uint8Array, Uint8ClampedArray, Int16Array, Uint16Array,
		Int32Array, Uint32Array, Float32Array, Float64Array];
	types.forEach(function(T) {
		var bpe = T.BYTES_PER_ELEMENT;
		var buf = new ArrayBuffer(bpe * 8);
		var a = new T(buf, bpe * 4, 4); // offset of 4 elements, length 4
		for (var i = 0; i < 4; i++) {
			a[i] = i + 1;
		}
		var w = a.with(2, 99);
		if (!(w instanceof T)) {
			throw new Error(T.name + ": with returned wrong type");
		}
		if (w.length !== 4 || w.byteOffset !== 0) {
			throw new Error(T.name + ": with length/offset=" + w.length + "/" + w.byteOffset);
		}
		if (w[0] !== 1 || w[1] !== 2 || w[2] !== 99 || w[3] !== 4) {
			throw new Error(T.name + ": with=[" + w.join(",") + "], expected [1,2,99,4]");
		}
		if (a[0] !== 1 || a[1] !== 2 || a[2] !== 3 || a[3] !== 4) {
			throw new Error(T.name + ": source corrupted by with=[" + a.join(",") + "]");
		}
		var r = a.toReversed();
		if (r.length !== 4 || r.byteOffset !== 0) {
			throw new Error(T.name + ": toReversed length/offset=" + r.length + "/" + r.byteOffset);
		}
		if (r[0] !== 4 || r[1] !== 3 || r[2] !== 2 || r[3] !== 1) {
			throw new Error(T.name + ": toReversed=[" + r.join(",") + "], expected [4,3,2,1]");
		}
		if (a[0] !== 1 || a[1] !== 2 || a[2] !== 3 || a[3] !== 4) {
			throw new Error(T.name + ": source corrupted by toReversed=[" + a.join(",") + "]");
		}
	});
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArrayWithToReversedBigIntNonZeroOffset(t *testing.T) {
	const SCRIPT = `
	[BigInt64Array, BigUint64Array].forEach(function(T) {
		var bpe = T.BYTES_PER_ELEMENT; // 8
		var buf = new ArrayBuffer(bpe * 8);
		var a = new T(buf, bpe * 4, 4);
		for (var i = 0; i < 4; i++) {
			a[i] = BigInt(i + 1);
		}
		var w = a.with(2, 99n);
		if (w[0] !== 1n || w[1] !== 2n || w[2] !== 99n || w[3] !== 4n) {
			throw new Error(T.name + ": with=[" + w.join(",") + "], expected [1,2,99,4]");
		}
		var r = a.toReversed();
		if (r[0] !== 4n || r[1] !== 3n || r[2] !== 2n || r[3] !== 1n) {
			throw new Error(T.name + ": toReversed=[" + r.join(",") + "], expected [4,3,2,1]");
		}
		if (a[0] !== 1n || a[1] !== 2n || a[2] !== 3n || a[3] !== 4n) {
			throw new Error(T.name + ": source corrupted=[" + a.join(",") + "]");
		}
	});
	`
	testScript(SCRIPT, _undefined, t)
}

// A view positioned at the very end of its buffer: the buggy implementation
// would write past the end of the destination buffer (out-of-bounds).
func TestTypedArrayWithToReversedOffsetAtBufferEnd(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(16);
	var a = new Uint8Array(buf, 12, 4); // last 4 bytes
	a[0] = 1; a[1] = 2; a[2] = 3; a[3] = 4;
	var w = a.with(0, 9);
	if (w[0] !== 9 || w[1] !== 2 || w[2] !== 3 || w[3] !== 4) {
		throw new Error("with=[" + w.join(",") + "], expected [9,2,3,4]");
	}
	var r = a.toReversed();
	if (r[0] !== 4 || r[1] !== 3 || r[2] !== 2 || r[3] !== 1) {
		throw new Error("toReversed=[" + r.join(",") + "], expected [4,3,2,1]");
	}

	// A wider element type with a non-zero offset near the end of the buffer.
	var buf2 = new ArrayBuffer(32);
	var f = new Float64Array(buf2, 16, 2); // offset 2 elements, length 2
	f[0] = 1.5; f[1] = 2.5;
	var fw = f.with(1, 3.5);
	if (fw[0] !== 1.5 || fw[1] !== 3.5) {
		throw new Error("fw=[" + fw.join(",") + "], expected [1.5,3.5]");
	}
	var fr = f.toReversed();
	if (fr[0] !== 2.5 || fr[1] !== 1.5) {
		throw new Error("fr=[" + fr.join(",") + "], expected [2.5,1.5]");
	}
	`
	testScript(SCRIPT, _undefined, t)
}

// Ensures that with()/toReversed() on an offset view do not corrupt other views
// sharing the same underlying buffer (the result must use a brand new buffer).
func TestTypedArrayWithToReversedDoesNotTouchSourceBuffer(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(8);
	var head = new Uint8Array(buf, 0, 4);
	var tail = new Uint8Array(buf, 4, 4);
	head[0] = 100; head[1] = 101; head[2] = 102; head[3] = 103;
	tail[0] = 1; tail[1] = 2; tail[2] = 3; tail[3] = 4;

	var w = tail.with(0, 9);
	var r = tail.toReversed();

	if (w.buffer === buf || r.buffer === buf) {
		throw new Error("result must not alias the source buffer");
	}
	if (head[0] !== 100 || head[1] !== 101 || head[2] !== 102 || head[3] !== 103) {
		throw new Error("head corrupted=[" + head.join(",") + "]");
	}
	if (tail[0] !== 1 || tail[1] !== 2 || tail[2] !== 3 || tail[3] !== 4) {
		throw new Error("tail corrupted=[" + tail.join(",") + "]");
	}
	`
	testScript(SCRIPT, _undefined, t)
}

func TestTypedArraySubarraySet(t *testing.T) {
	const SCRIPT = `
	const u = new Uint8Array(4)
	const s = u.subarray(1, 4);
	s.set([1,2,3], 0);
	if (s.length !== 3) {
		throw new Error("s.length=" + s.length + ", expected 3");
	}
	if (s[0] !== 1 || s[1] !== 2 || s[2] !== 3) {
		throw new Error("s=[" + s.join(",") + "], expected [1,2,3]");
	}
	`

	testScript(SCRIPT, _undefined, t)
}

func TestTypedArrayWithOffset(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(16);
	var src = new Uint8Array(buf, 4, 4);
	src[0] = 10; src[1] = 20; src[2] = 30; src[3] = 40;
	var result = src.with(1, 99);
	result[0] === 10 && result[1] === 99 && result[2] === 30 && result[3] === 40
		&& result.length === 4 && src[1] === 20;
	`
	testScript(SCRIPT, valueTrue, t)
}

func TestTypedArrayToReversedOffset(t *testing.T) {
	const SCRIPT = `
	var buf = new ArrayBuffer(16);
	var src = new Uint8Array(buf, 4, 4);
	src[0] = 10; src[1] = 20; src[2] = 30; src[3] = 40;
	var result = src.toReversed();
	result[0] === 40 && result[1] === 30 && result[2] === 20 && result[3] === 10
		&& result.length === 4 && src[0] === 10;
	`
	testScript(SCRIPT, valueTrue, t)
}

func TestTypedArraySliceEvilSpecies(t *testing.T) {
	// If typedArraySpeciesCreate() returns an array with non-zero offset, the offset must be taken into account
	// by the slice() method.
	const SCRIPT = `
	let count = 0;
	class EvilArray extends Uint8Array {
		constructor(length) {
			super(length);
			if (++count === 1) {
				return;
			}
			const buf = new ArrayBuffer(length+1);
			const ret = new Uint8Array(buf, 1, length);
			Object.setPrototypeOf(ret, EvilArray.prototype);
			return ret;
		}
	}

	const a = new EvilArray(2);
	a[0] = 1;
	a[1] = 2;

	a.slice(1, 2)[0];
	`
	testScript(SCRIPT, intToValue(2), t)
}

func TestTypedArraySliceEvilSpecies1(t *testing.T) {
	// If typedArraySpeciesCreate() returns an array with non-zero offset and different type, the offset must be taken into account
	// by the slice() method.
	const SCRIPT = `
	let count = 0;
	class EvilArray extends Uint8Array {
		constructor(length) {
			super(length);
			if (++count === 1) {
				return;
			}
			const buf = new ArrayBuffer((length+1)*2);
			const ret = new Uint16Array(buf, 2, length);
			Object.setPrototypeOf(ret, EvilArray.prototype);
			return ret;
		}
	}

	const a = new EvilArray(2);
	a[0] = 1;
	a[1] = 2;

	a.slice(1, 2)[0];
	`
	testScript(SCRIPT, intToValue(2), t)
}

func TestTypedArrayMapEvilSpecies(t *testing.T) {
	// If typedArraySpeciesCreate() returns an array with non-zero offset, the offset must be taken into account
	// by the map() method.
	const SCRIPT = `
	let count = 0;
	class EvilArray extends Uint8Array {
		constructor(length) {
			super(length);
			if (++count === 1) {
				return;
			}
			const buf = new ArrayBuffer(length+1);
			const ret = new Uint8Array(buf, 1, length);
			Object.setPrototypeOf(ret, EvilArray.prototype);
			return ret;
		}
	}

	const a = new EvilArray(2);
	a[0] = 1;
	a[1] = 2;

	a.map(x => x)[0];
	`
	testScript(SCRIPT, intToValue(1), t)
}

func TestTypedArrayFilterEvilSpecies(t *testing.T) {
	// If typedArraySpeciesCreate() returns an array with non-zero offset, the offset must be taken into account
	// by the filter() method.
	const SCRIPT = `
	class Evil extends Uint8Array {
	  static get [Symbol.species]() {
		return function (len) { return new Uint8Array(new ArrayBuffer(len + 2), 2, len); };
	  }
	}

	assert(compareArray(new Evil([10,11,12,13,14,15]).filter(x => x % 2 === 0), [10,12,14]));
	`
	testScriptWithTestLib(SCRIPT, _undefined, t)
}

func TestToSortedSubarray(t *testing.T) {
	const SCRIPT = `
	var base = new Uint8Array([9,8,7,6,5,4,3,2,1,0]);
	var view = base.subarray(3);          // = [6,5,4,3,2,1,0]
	assert(compareArray(view.toSorted(), [0,1,2,3,4,5,6]));
	`
	testScriptWithTestLib(SCRIPT, _undefined, t)
}

func TestTypedArrayFromPrimitive(t *testing.T) {
	const SCRIPT = `
	assert(compareArray(Uint8Array.from("123"), [1, 2, 3]), "string");
	assert(compareArray(Uint16Array.from("abc", c => c.charCodeAt(0)), [97, 98, 99]), "string, mapFn");
	assert.sameValue(Uint8Array.from(5).length, 0, "number");

	var receiver;
	Object.defineProperty(Number.prototype, Symbol.iterator, {
		value: function() {
			"use strict";
			receiver = this;
			return [7][Symbol.iterator]();
		}
	});
	assert(compareArray(Uint8Array.from(5), [7]), "number, iterator");
	assert.sameValue(receiver, 5, "iterator receiver");
	`
	testScriptWithTestLib(SCRIPT, _undefined, t)
}

func TestArrayBufferTransfer(t *testing.T) {
	const SCRIPT = `
var buf = new Uint8Array([1, 2, 3, 4]).buffer;
var ta = new Uint8Array(buf);
var grown = buf.transfer(6);
assert(buf.detached, "source is detached");
assert.sameValue(buf.byteLength, 0, "source byteLength");
assert.sameValue(ta.length, 0, "view on the source");
assert(!grown.detached, "result is not detached");
assert(compareArray(new Uint8Array(grown), [1, 2, 3, 4, 0, 0]), "grown contents");

var shrunk = grown.transferToFixedLength(2);
assert(compareArray(new Uint8Array(shrunk), [1, 2]), "shrunk contents");
assert.sameValue(shrunk.transfer().byteLength, 2, "default length");
assert.throws(TypeError, function() { shrunk.transfer(); }, "detached source");
`
	testScriptWithTestLib(SCRIPT, _undefined, t)
}

func TestArrayBufferTransferFromGo(t *testing.T) {
	vm := New()
	data := []byte{1, 2, 3}
	ab := vm.NewArrayBuffer(data)
	vm.Set("ab", ab)
	v, err := vm.RunString(`ab.transfer()`)
	if err != nil {
		t.Fatal(err)
	}
	if !ab.Detached() || ab.Bytes() != nil {
		t.Fatal("source ArrayBuffer is not detached")
	}
	b := v.Export().(ArrayBuffer).Bytes()
	if !bytes.Equal(b, []byte{1, 2, 3}) {
		t.Fatalf("unexpected contents: %v", b)
	}
	if &b[0] != &data[0] {
		t.Fatal("the data was copied instead of being transferred")
	}
}

func TestTypedArraySearchNegativeZero(t *testing.T) {
	const SCRIPT = `
	[Float32Array, Float64Array].forEach(function(TA) {
		const name = TA.name;
		const a = new TA([-0, 1]);
		assert.sameValue(a.includes(0), true, name + " includes 0");
		assert.sameValue(a.includes(-0), true, name + " includes -0");
		assert.sameValue(a.indexOf(0), 0, name + " indexOf 0");
		assert.sameValue(a.indexOf(-0), 0, name + " indexOf -0");
		assert.sameValue(a.lastIndexOf(0), 0, name + " lastIndexOf 0");
		assert.sameValue(a.lastIndexOf(-0), 0, name + " lastIndexOf -0");
		assert.sameValue(new TA([-0, 0]).indexOf(0), 0, name + " indexOf first of -0, +0");
		assert.sameValue(new TA([-0, 0]).lastIndexOf(-0), 1, name + " lastIndexOf last of -0, +0");
		assert.sameValue(a.includes(2), false, name + " includes 2");
		assert.sameValue(a.includes(0, 1), false, name + " includes 0 after fromIndex");
		assert.sameValue(new TA([1, -0]).lastIndexOf(0, 0), -1, name + " lastIndexOf 0 before fromIndex");
		assert.sameValue(new TA([NaN]).includes(NaN), true, name + " includes NaN");
		assert.sameValue(new TA([NaN]).indexOf(NaN), -1, name + " indexOf NaN");
	});
	const i = new Int8Array([0]);
	assert.sameValue(i.includes(-0), true, "Int8Array includes -0");
	assert.sameValue(i.indexOf(-0), 0, "Int8Array indexOf -0");
	assert.sameValue(i.lastIndexOf(-0), 0, "Int8Array lastIndexOf -0");
	assert.sameValue(new BigInt64Array([0n]).includes(0n), true, "BigInt64Array includes 0n");
	`
	testScriptWithTestLib(SCRIPT, _undefined, t)
}
