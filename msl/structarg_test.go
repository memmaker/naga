package msl_test

import (
	"strings"
	"testing"

	"github.com/gogpu/naga"
	"github.com/gogpu/naga/msl"
)

// A vertex entry point taking a struct of inputs and a loose @location argument.
func TestVertexStructAndLooseInputs(t *testing.T) {
	src := `
struct VIn { @location(0) pos: vec3f, @location(1) uv: vec2f, @location(2) face: u32 }
struct VOut { @builtin(position) clip: vec4f, @location(0) uv: vec2f }
@vertex
fn vs(v: VIn, @location(3) ii: u32) -> VOut {
	var o: VOut;
	o.clip = vec4f(v.pos, f32(ii + v.face));
	o.uv = v.uv;
	return o;
}
`
	ast, err := naga.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	m, err := naga.LowerWithSource(ast, src)
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := msl.Compile(m, msl.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "varyings.,") || strings.Contains(out, "varyings. ") {
		t.Fatalf("struct argument members unnamed:\n%s", out)
	}
	for _, want := range []string{"[[attribute(0)]]", "[[attribute(1)]]", "[[attribute(2)]]", "[[attribute(3)]]"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %s:\n%s", want, out)
		}
	}
	t.Log(out)
}
