package goja

import "testing"

// Go functions in nuclei return nil slices and maps for "no data". Scripts
// read .length or keys on the result, so a nil slice or map must stay an
// empty array-like or object, not null.
func TestNilSliceAndMapStayObjects(t *testing.T) {
	r := New()
	_ = r.Set("nilBytes", func() []byte { return nil })
	_ = r.Set("nilStrings", func() ([]string, error) { return nil, nil })
	_ = r.Set("nilMap", func() map[string]string { return nil })
	v, err := r.RunString(`nilBytes().length === 0 && nilStrings().length === 0 && Object.keys(nilMap()).length === 0`)
	if err != nil {
		t.Fatal(err)
	}
	if !v.ToBoolean() {
		t.Fatal("nil slice or map was not converted to an empty object")
	}
}
