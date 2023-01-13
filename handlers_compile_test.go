package main

import "testing"

func TestAddGlue(t *testing.T) {
	cases := []struct {
		expected string
		input    string
	}{
		{
			expected: wasmGlue + `#[no_mangle]
pub extern "C" fn lib() {
	println!("hello world");
}
`,
			input: `fn main() {
	println!("hello world");
}
`,
		},
	}
	for _, c := range cases {
		actual := addGlue(c.input)
		if actual != c.expected {
			t.Errorf("expected %v got %v", c.expected, actual)
		}
	}
}
