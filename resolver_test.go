package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveIgnoreField(t *testing.T) {
	type temp struct {
		Cat *Cat
	}

	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)
	err = c.ProvideAs(&Cat{Name: "B"}, (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(NewPetCat("C"), (*Animal)(nil))
	assert.Nil(t, err)

	s := &temp{}
	err = c.Resolve(s)
	assert.Nil(t, err)

	assert.True(t, s.Cat == nil)
}

func TestResolveWrongInjectTagUse(t *testing.T) {
	type temp struct {
		Cat Cat `inject:""`
	}

	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)
	err = c.ProvideAs(&Cat{Name: "B"}, (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(NewPetCat("C"), (*Animal)(nil))
	assert.Nil(t, err)

	s := &temp{}
	err = c.Resolve(s)
	assert.NotNil(t, err)
}

func TestResolve(t *testing.T) {
	type temp struct {
		Cat    *Cat   `inject:""`
		Pet    Pet    `inject:""`
		Animal Animal `inject:""`
	}

	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)
	err = c.ProvideAs(&Cat{Name: "B"}, (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(NewPetCat("C"), (*Animal)(nil))
	assert.Nil(t, err)

	s := &temp{}
	err = c.Resolve(s)
	assert.Nil(t, err)

	assert.NotNil(t, s.Cat)
	assert.NotNil(t, s.Pet)
	assert.NotNil(t, s.Animal)

	assert.Equal(t, s.Cat.GetName(), "A")
	assert.Equal(t, s.Pet.GetName(), "B")
	assert.Equal(t, s.Animal.GetName(), "C")
}

func TestResolveTypeMismatch(t *testing.T) {
	type temp struct {
		Cat *Cat `inject:""`
	}

	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)
	err = c.Provide(&temp{})
	assert.Nil(t, err)
	err = c.ProvideAs(NewPetDog("B"), (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(NewAnimalCat("C"), (*Animal)(nil))
	assert.Nil(t, err)

	s := &temp2{}
	err = c.Resolve(s)
	assert.Nil(t, err)

	assert.NotNil(t, s.Temp.Cat)
	assert.NotNil(t, s.Temp.Pet)
	assert.NotNil(t, s.Temp.Animal)

	assert.Equal(t, s.Temp.Cat.GetName(), "A")
	assert.Equal(t, s.Temp.Pet.GetName(), "B")
	assert.Equal(t, s.Temp.Animal.GetName(), "C")
}

func TestResolveRecursive(t *testing.T) {
	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)
	err = c.ProvideAs(NewPetDog("B"), (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(NewAnimalCat("C"), (*Animal)(nil))
	assert.Nil(t, err)
	err = c.Provide(&temp{})
	assert.Nil(t, err)

	s := &temp2{}
	err = c.Resolve(s)
	assert.Nil(t, err)

	assert.NotNil(t, s.Temp.Cat)
	assert.NotNil(t, s.Temp.Pet)
	assert.NotNil(t, s.Temp.Animal)

	assert.Equal(t, s.Temp.Cat.GetName(), "A")
	assert.Equal(t, s.Temp.Pet.GetName(), "B")
	assert.Equal(t, s.Temp.Animal.GetName(), "C")
}

func TestResolveOnNil(t *testing.T) {
	c := NewContainer()

	err := c.Resolve(nil)
	assert.NotNil(t, err)
}

func TestResolveOnNonPointerStruct(t *testing.T) {
	c := NewContainer()

	a := 3
	cat := Cat{Name: "A"}

	err := c.Resolve(&a)
	assert.NotNil(t, err)

	err = c.Resolve(cat)
	assert.NotNil(t, err)
}

func TestResolveStructFieldCannotSet(t *testing.T) {
	type temp struct {
		cat *Cat `inject:""`
	}

	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)

	s := &temp{}
	err = c.Resolve(s)
	assert.NotNil(t, err)
	assert.True(t, s.cat == nil)
}

func TestResolveStructFieldNotFound(t *testing.T) {
	type temp struct {
		Cat *Cat `inject:""`
	}

	c := NewContainer()

	s := &temp{}

	err := c.Resolve(s)
	assert.Nil(t, err)
	assert.True(t, s.Cat != nil)
	assert.True(t, s.Cat.Name == "")
}

func TestResolveByFunction(t *testing.T) {
	f1 := func() Pet { return &Cat{Name: "A"} }
	f2 := func() *Cat { return &Cat{Name: "B"} }

	type temp struct {
		Cat1 Pet  `inject:""`
		Cat2 *Cat `inject:""`
	}

	te := temp{}

	c := NewContainer()

	err := c.Provide(f1)
	assert.Nil(t, err)

	err = c.Provide(f2)
	assert.Nil(t, err)

	err = c.Resolve(&te)
	assert.Nil(t, err)
	assert.True(t, te.Cat1.GetName() == "A")
	assert.True(t, te.Cat2.GetName() == "B")
}

func TestResolveStructPointerFromProvide(t *testing.T) {
	type target struct {
		MyCat *Cat `inject:""`
	}

	c := NewContainer()

	cat := &Cat{Name: "TestCat"}
	err := c.Provide(cat)
	assert.Nil(t, err)

	tgt := &target{}
	err = c.Resolve(tgt)
	assert.Nil(t, err)

	assert.NotNil(t, tgt.MyCat)
	assert.Equal(t, "TestCat", tgt.MyCat.GetName())
	assert.Equal(t, "TestCat", tgt.MyCat.Name)
}
