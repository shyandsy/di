package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveIgnoreField(t *testing.T) {
	type temp struct {
		Cat1 Cat  `inject:""` // ignore: struct, not pointer
		cat2 *Cat `inject:""` // ignore: unexported
		Cat3 *Cat // ignore: no inject tag
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

	assert.Equal(t, s.Cat1.GetName(), "")
	assert.Nil(t, s.cat2)
	assert.Nil(t, s.Cat3)
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
