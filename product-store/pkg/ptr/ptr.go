// package ptr provides convenience for converting values to pointers.
package ptr

// To returns a pointer to the value passed in.
//
// If you have a struct like:
//
//	type Foo struct {
//		Bar *int
//	}
//
// Then you can use ptr.To to construct a Foo directly without needing to
// declare Foo beforehand or store a temporary variable for Bar. For example:
//
//	// bad: requires separate declaration, initialization, and assignment
//	var f1 Foo
//	f.Bar = new(int)
//	*f.Bar = 42
//
//	// bad: requires a temporary variable
//	var bar = 42
//	var f2 = Foo{Bar: &bar}
//
// Instead, you can use ptr.To to instantiate Foo like you normally would:
//
//	var foo = Foo{
//		Bar: ptr.To(42),
//	}
func To[T any](v T) *T {
	return &v
}
