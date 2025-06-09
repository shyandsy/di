package di

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func PrintCat(cat Cat) {
	fmt.Println(cat.Name)
}

func PrintCatAndAnimal(cat1 *Cat, cat2 Cat, animal Animal) {
	fmt.Println(fmt.Sprintf("cat1:%s\ncat2:%s\nanimal:%s\n", cat1.Name, cat2.Name, animal.GetName()))
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
	assert.Nil(t, err)
	err = c.Invoke(PrintCatAndAnimal)
	assert.Nil(t, err)
}
