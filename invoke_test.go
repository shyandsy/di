package di

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func PrintCat(cat Cat) {
	fmt.Println(cat.Name)
}

func PrintCatWithInvalidParameter(cat *Cat, a int) {
	fmt.Println(cat.Name)
}

func PrintCatPointer(cat *Cat) {
	fmt.Println(cat.Name)
}

func PrintCatAndAnimal(cat1 *Cat, cat2 *Cat, animal Animal) {
	fmt.Println(fmt.Sprintf("cat1:%s\ncat2:%s\nanimal:%s\n", cat1.Name, cat2.Name, animal.GetName()))
}

func PrintCatAndAnimalRecursive(s *temp2) {
	fmt.Println(fmt.Sprintf("cat1:%s\ncat2:%s\nanimal:%s\n", s.Temp.Cat.Name, s.Temp.Pet.GetName(), s.Temp.Animal.GetName()))
}

func TestInvoke(t *testing.T) {
	c := NewContainer()

	cat := &Cat{Name: "A"}
	animal := NewAnimalCat("B")

	err := c.Provide(cat)
	assert.Nil(t, err)
	err = c.ProvideAs(animal, (*Animal)(nil))
	assert.Nil(t, err)

	err = c.Invoke(PrintCat)
	assert.NotNil(t, err)

	err = c.Invoke(PrintCatWithInvalidParameter)
	assert.NotNil(t, err)

	err = c.Invoke(PrintCatPointer)
	assert.Nil(t, err)

	err = c.Invoke(PrintCatAndAnimal)
	assert.Nil(t, err)
}

func TestInvokeRecursive(t *testing.T) {
	c := NewContainer()

	cat := &Cat{Name: "A"}
	pet := NewPetDog("B")
	animal := NewAnimalCat("C")

	err := c.Provide(cat)
	assert.Nil(t, err)
	err = c.ProvideAs(pet, (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(animal, (*Animal)(nil))
	assert.Nil(t, err)

	err = c.Invoke(PrintCatAndAnimalRecursive)
	assert.Nil(t, err)
}

func TestInvokeOnNil(t *testing.T) {
	c := NewContainer()

	err := c.Invoke(nil)
	assert.NotNil(t, err)
}

func TestInvokeOnNonFunction(t *testing.T) {
	c := NewContainer()

	// struct
	err := c.Invoke(&Cat{})
	assert.NotNil(t, err)

	// interface
	var a Animal
	err = c.Invoke(&a)
	assert.NotNil(t, err)
}
