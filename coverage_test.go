package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---- helper types used only by the error/edge-path coverage tests ----

// needsAnimal has an interface-typed inject field whose dependency is never
// registered, so resolving it must fail with "dependency not found".
type needsAnimal struct {
	A Animal `inject:""`
}

// wrapNeedsAnimal nests a *needsAnimal, so resolving it recurses into
// needsAnimal and the nested failure must propagate back up.
type wrapNeedsAnimal struct {
	Inner *needsAnimal `inject:""`
}

// catShape is structurally similar to Cat (same Name field) but a distinct
// type, used to exercise the "stored dep not assignable to field, copy
// matching non-inject fields" branch.
type catShape struct {
	Name string
}

// wantCat injects a *Cat; the container will be seeded with a *catShape under
// Cat's type key to force the field-copy branch.
type wantCat struct {
	MyCat *Cat `inject:""`
}

// ---------------------------------------------------------------------------
// parser.go
// ---------------------------------------------------------------------------

// TestStructTypeAccessors covers the two previously-untested exported getters
// Struct.TypePkg and Struct.TypeName (both were 0%), plus the happy path of
// parseStruct on a struct pointer.
func TestStructTypeAccessors(t *testing.T) {
	c := NewContainer().(*container)

	s, err := c.parseStruct(reflect.TypeOf(&Cat{}))
	assert.Nil(t, err)
	assert.Equal(t, "github.com/shyandsy/di", s.TypePkg())
	assert.Equal(t, "Cat", s.TypeName())
	assert.Equal(t, "github.com/shyandsy/di:Cat", s.FullType())
}

// TestParseStructNil covers the nil-type guard in parseStruct.
func TestParseStructNil(t *testing.T) {
	c := NewContainer().(*container)

	s, err := c.parseStruct(nil)
	assert.Nil(t, s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "tp must be not nil")
}

// TestParseStructInvalidKind covers the kind-validation guard: a plain int is
// neither struct/interface nor a pointer to one.
func TestParseStructInvalidKind(t *testing.T) {
	c := NewContainer().(*container)

	s, err := c.parseStruct(reflect.TypeOf(3))
	assert.Nil(t, s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "tp must be a struct | interface")
}

// ---------------------------------------------------------------------------
// invoke.go
// ---------------------------------------------------------------------------

// TestInvokeDependencyNotFound covers the "dependency not found" return: the
// function needs a *Dog but nothing is registered under that type.
func TestInvokeDependencyNotFound(t *testing.T) {
	c := NewContainer()

	_, err := c.Invoke(func(d *Dog) {})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "dependency not found")
}

// TestInvokeInjectFieldResolveError covers the branch where the stored instance
// is assignable to the parameter type AND the parameter struct has inject
// fields, but resolving those fields fails (Animal is never registered).
func TestInvokeInjectFieldResolveError(t *testing.T) {
	c := NewContainer()

	assert.Nil(t, c.Provide(&needsAnimal{}))

	_, err := c.Invoke(func(x *needsAnimal) {})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "dependency not found")
}

// TestInvokeNotAssignableResolveError covers the else-branch (stored instance
// not assignable to the parameter type): a *Dog is registered under
// needsAnimal's type key, so Invoke builds a fresh needsAnimal and resolving it
// fails on the unregistered Animal field.
func TestInvokeNotAssignableResolveError(t *testing.T) {
	c := NewContainer()

	// struct target: ProvideAs stores &Dog{} under needsAnimal's type key
	// without an interface-implements check.
	assert.Nil(t, c.ProvideAs(&Dog{Name: "x"}, (*needsAnimal)(nil)))

	_, err := c.Invoke(func(x *needsAnimal) {})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "dependency not found")
}

// TestInvokeNotAssignableFreshInstance covers the else-branch success path: a
// *Dog is stored under *Cat's key (not assignable to *Cat), so Invoke resolves
// a fresh, empty *Cat and passes it in.
func TestInvokeNotAssignableFreshInstance(t *testing.T) {
	c := NewContainer()

	assert.Nil(t, c.ProvideAs(&Dog{Name: "ignored"}, (*Cat)(nil)))

	var got *Cat
	vals, err := c.Invoke(func(x *Cat) string {
		got = x
		return x.Name
	})
	assert.Nil(t, err)
	assert.NotNil(t, got)
	// A brand-new empty Cat was built (not the stored Dog), so Name is empty.
	assert.Equal(t, "", vals[0].String())
	assert.Equal(t, "", got.Name)
}

// ---------------------------------------------------------------------------
// resolve.go
// ---------------------------------------------------------------------------

// TestResolveInterfaceFieldNotFound covers both the "interface dependency not
// found" return (unregistered interface field => error, not auto-new) and the
// recursive-Resolve error propagation: wrapNeedsAnimal -> needsAnimal.A fails.
func TestResolveInterfaceFieldNotFound(t *testing.T) {
	c := NewContainer()

	err := c.Resolve(&wrapNeedsAnimal{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "dependency not found")
}

// TestResolveFieldCopyNotAssignable covers the branch where the stored instance
// is a different struct type than the inject field: matching non-inject fields
// are copied over. catShape{Name} -> Cat{Name}.
func TestResolveFieldCopyNotAssignable(t *testing.T) {
	c := NewContainer()

	// struct target: store &catShape{} under Cat's type key.
	assert.Nil(t, c.ProvideAs(&catShape{Name: "copied"}, (*Cat)(nil)))

	w := &wantCat{}
	err := c.Resolve(w)
	assert.Nil(t, err)
	assert.NotNil(t, w.MyCat)
	assert.Equal(t, "copied", w.MyCat.Name)
}

// ---------------------------------------------------------------------------
// find.go
// ---------------------------------------------------------------------------

// TestFindParseStructError covers the parseStruct error path in Find: a
// pointer-to-pointer target type (**Cat) is not a valid resolvable type.
func TestFindParseStructError(t *testing.T) {
	c := NewContainer()

	assert.Nil(t, c.Provide(&Cat{Name: "A"}))

	var pp **Cat
	err := c.Find(&pp)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "tp must be a struct | interface")
}

// TestFindTargetCannotSet covers the "object cannot be set" return: a typed-nil
// *Cat has no addressable element to assign into, even though Cat is registered.
func TestFindTargetCannotSet(t *testing.T) {
	c := NewContainer()

	assert.Nil(t, c.Provide(&Cat{Name: "A"}))

	err := c.Find((*Cat)(nil))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "object cannot be set")
}

// TestFindStructResolveError covers the Resolve-error path in Find's struct
// branch: needsAnimal is registered (so it is found), but resolving its
// unregistered Animal field fails.
func TestFindStructResolveError(t *testing.T) {
	c := NewContainer()

	assert.Nil(t, c.Provide(&needsAnimal{}))

	err := c.Find(&needsAnimal{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "dependency not found")
}
