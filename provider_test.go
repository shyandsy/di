package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvide(t *testing.T) {
	c := NewContainer()

	cat1 := &Cat{Name: "A"}
	dog1 := &Dog{Name: "B"}

	err := c.Provide(cat1)
	assert.Nil(t, err)
	err = c.Provide(dog1)
	assert.Nil(t, err)

	cat2 := &Cat{}
	dog2 := &Dog{}

	err = c.Find(cat2)
	assert.Nil(t, err)
	assert.Equal(t, cat1.Name, cat2.Name)

	err = c.Find(dog2)
	assert.Nil(t, err)
	assert.Equal(t, dog1.Name, dog2.Name)
}

func TestProvideAsInterface(t *testing.T) {
	c := NewContainer()

	cat0 := &Cat{Name: "azure"}
	cat1 := NewPetCat("cat")
	animalCat := NewAnimalCat("ccc")

	err := c.Provide(cat0)
	assert.Nil(t, err)
	err = c.ProvideAs(cat1, (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(animalCat, (*Animal)(nil))
	assert.Nil(t, err)

	cat := &Cat{Name: ""}
	var P Pet
	var A Animal

	err = c.Find(cat)
	assert.Nil(t, err)
	err = c.Find(&P)
	assert.Nil(t, err)
	err = c.Find(&A)
	assert.Nil(t, err)

	assert.Equal(t, cat0.GetName(), cat.Name)
	assert.Equal(t, cat1.GetName(), P.GetName())
	assert.Equal(t, animalCat.GetName(), A.GetName())
}

func TestProvideAsOverride(t *testing.T) {
	c := NewContainer()

	pet1 := NewPetCat("aaa")
	pet2 := NewPetDog("bbb")

	err := c.ProvideAs(pet1, (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(pet2, (*Pet)(nil))
	assert.Nil(t, err)

	var p Pet

	err = c.Find(&p)
	assert.Nil(t, err)

	assert.NotEqual(t, p.GetName(), pet1.GetName())
	assert.Equal(t, p.GetName(), pet2.GetName())
}

func TestProvideInvalidType(t *testing.T) {
	c := NewContainer()

	a := 3

	err := c.Provide(nil)
	assert.NotNil(t, err)

	err = c.Provide(&a)
	assert.NotNil(t, err)

	err = c.Provide(&a)
	assert.NotNil(t, err)

	err = c.ProvideAs(nil, nil)
	assert.NotNil(t, err)

	err = c.ProvideAs(&a, (*Pet)(nil))
	assert.NotNil(t, err)
}

func TestProvideAsNonValue(t *testing.T) {
	c := NewContainer()

	err := c.ProvideAs((*Cat)(nil), (*Pet)(nil))
	assert.NotNil(t, err)

	err = c.ProvideAs((*Cat)(nil), nil)
	assert.NotNil(t, err)
}

func TestProvideAsTargetNonPointInterface(t *testing.T) {
	c := NewContainer()

	err := c.ProvideAs(&Cat{Name: "A"}, Cat{Name: "A"})
	assert.NotNil(t, err)
}

func TestProvideAsWrongImplementation(t *testing.T) {
	c := NewContainer()

	// not implement the interface
	err := c.ProvideAs(&temp{}, (*Pet)(nil))
	assert.NotNil(t, err)
}

func TestProvideFunction(t *testing.T) {
	f1 := func() {}
	f2 := func() Cat { return Cat{} }
	f3 := func() (*Cat, error) { return nil, nil }
	f4 := func() *Cat { return &Cat{Name: "A"} }
	f5 := func() Pet { return &Cat{Name: "B"} }

	c := NewContainer()

	err := c.Provide(f1)
	assert.NotNil(t, err)

	err = c.Provide(f2)
	assert.NotNil(t, err)

	err = c.Provide(f3)
	assert.NotNil(t, err)

	err = c.Provide(f4)
	assert.Nil(t, err)

	cat := Cat{}
	err = c.Find(&cat)
	assert.Nil(t, err)
	assert.True(t, cat.Name == "A")

	err = c.Provide(f5)
	assert.Nil(t, err)

	var pet Pet
	err = c.Find(&pet)
	assert.Nil(t, err)
	assert.True(t, pet.GetName() == "B")

}
