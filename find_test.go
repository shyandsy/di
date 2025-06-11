package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindRecursive(t *testing.T) {
	c := NewContainer()

	err := c.Provide(&Cat{Name: "A"})
	assert.Nil(t, err)
	err = c.ProvideAs(NewPetCat("B"), (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(NewAnimalCat("C"), (*Animal)(nil))
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
