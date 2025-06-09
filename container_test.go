package di

type Animal interface {
	GetName() string
}

type Pet interface {
	Animal
}

type Dog struct {
	Name string
}

func (d Dog) GetName() string {
	return d.Name
}

type Cat struct {
	Name string
}

func (d Cat) GetName() string {
	return d.Name
}

func NewAnimalCat(name string) Animal {
	return &Cat{Name: name}
}

func NewPetCat(name string) Pet {
	return &Cat{Name: name}
}

func NewPetDog(name string) Pet {
	return &Dog{Name: name}
}
