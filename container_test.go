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

type Cat struct {
	Name string
}

type temp struct {
	Cat    *Cat   `inject:""`
	Pet    Pet    `inject:""`
	Animal Animal `inject:""`
}

type temp2 struct {
	Temp *temp `inject:""`
}

func (d Dog) GetName() string {
	return d.Name
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
