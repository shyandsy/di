package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindInterface(t *testing.T) {
	c := NewContainer()

	var s Animal

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)

	err = c.Find(&s)
	assert.NotNil(t, err)
	assert.Nil(t, s)

	err = c.ProvideAs(&Cat{Name: "A"}, (*Animal)(nil))
	assert.Nil(t, err)

	err = c.Find(&s)
	assert.Nil(t, err)
	assert.True(t, s != nil)
	assert.Equal(t, s.GetName(), "A")
}

func TestFindPointerStruct(t *testing.T) {
	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)
	err = c.ProvideAs(NewPetDog("B"), (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(NewAnimalCat("C"), (*Animal)(nil))
	assert.Nil(t, err)

	a := Cat{}
	err = c.Find(&a)
	assert.Nil(t, err)
	assert.Equal(t, a.GetName(), "A")
}

func TestFindWithInjectFields(t *testing.T) {
	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)
	err = c.ProvideAs(NewPetCat("B"), (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(NewAnimalCat("C"), (*Animal)(nil))
	assert.Nil(t, err)
	err = c.Provide(&temp{})
	assert.Nil(t, err)

	s := &temp{}
	err = c.Find(s)
	assert.Nil(t, err)

	assert.NotNil(t, s.Cat)
	assert.NotNil(t, s.Pet)
	assert.NotNil(t, s.Animal)

	assert.Equal(t, s.Cat.GetName(), "A")
	assert.Equal(t, s.Pet.GetName(), "B")
	assert.Equal(t, s.Animal.GetName(), "C")
}

func TestFindOnNil(t *testing.T) {
	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)

	err = c.Find(nil)
	assert.NotNil(t, err)
}

func TestFindOnNonPointer(t *testing.T) {
	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)

	s := Cat{}
	err = c.Find(s)
	assert.NotNil(t, err)
}

func TestFindOnNonPointerStructOrPointerInterface(t *testing.T) {
	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)

	a := 3
	err = c.Find(&a)
	assert.NotNil(t, err)

	b := func() {}
	err = c.Find(&b)
	assert.NotNil(t, err)
}

func TestFindUnwritableField(t *testing.T) {
	c := NewContainer()

	type temp struct {
		cat *Cat `inject:""`
	}

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)

	te := temp{}
	err = c.Find(&te)
	assert.NotNil(t, err)
}

func TestProvideFindByFunction(t *testing.T) {
	type unknow struct {
		UnknowDog *Dog `inject:""`
	}
	f1 := func() Pet { return &Cat{Name: "A"} }
	f2 := func() *Cat { return &Cat{Name: "B"} }
	f3 := func(a *unknow) Animal { return NewAnimalCat("C") }

	cat := Cat{}
	var p Pet
	var a Animal

	c := NewContainer()

	err := c.Provide(f1)
	assert.Nil(t, err)
	err = c.Provide(f2)
	assert.Nil(t, err)
	err = c.Provide(f3)
	assert.Nil(t, err)

	err = c.Find(&p)
	assert.Nil(t, err)
	assert.True(t, p.GetName() == "A")

	err = c.Find(&cat)
	assert.Nil(t, err)
	assert.True(t, cat.Name == "B")

	err = c.Find(&a)
	assert.NotNil(t, err)
	assert.True(t, a == nil)
}

func TestFindPointerType(t *testing.T) {
	c := NewContainer()

	cat := &Cat{Name: "TestCat"}
	err := c.Provide(cat)
	assert.Nil(t, err)

	var sqlDB *Cat
	err = c.Find(&sqlDB)
	assert.Nil(t, err)

	assert.NotNil(t, sqlDB)
	assert.Equal(t, "TestCat", sqlDB.GetName())
	assert.Equal(t, "TestCat", sqlDB.Name)
	assert.True(t, sqlDB == cat)
}
