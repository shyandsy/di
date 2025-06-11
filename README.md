# di
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Test Status](https://github.com/shyandsy/di/actions/workflows/go-test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/shyandsy/di)](https://goreportcard.com/report/github.com/shyandsy/di)
[![codecov](https://codecov.io/gh/shyandsy/di/graph/badge.svg?token=W8Z0SWZJG3)](https://codecov.io/gh/shyandsy/di)
![CodeQL](https://github.com/shyandsy/di/workflows/CodeQL/badge.svg)

A reflection based (DI)Dependency Injector component for golang project.

## Features
- [x] provide struct pointer as dependencies  
- [x] provide struct as interface type
- [x] find dependencies on struct fields with tag inject
- [x] use dependencies as parameter on invoke method
- [x] recursive inject on find/resolve/invoke

## Installation
installation
```
go install github.com/shyandsy/di
```

## Usage
please check unit test code

create a container
```go
c := di.NewContainer()
```

provide struct and find dependency
```go
cat1 := &Cat{Name: "A"}
dog1 := &Dog{Name: "B"}

// Provide
err := c.Provide(cat1)
err = c.Provide(dog1)

cat2 := &Cat{}
dog2 := &Dog{}
	
// Find：cat2.Name = "A"
err = c.Find(cat2)

// Find：dog2.Name = "B"
err = c.Find(dog2)
```

provide interface and find dependency
```go
cat0 := &Cat{Name: "azure"}
cat1 := NewPetCat("cat")
animalCat := NewAnimalCat("ccc")

// Provide and ProvideAs
err := c.Provide(cat0)
err = c.ProvideAs(cat1, (*Pet)(nil))
err = c.ProvideAs(animalCat, (*Animal)(nil))

cat := &Cat{Name: ""}
var P Pet
var A Animal

// Find struct: 
// cat.Name: "azure"
err = c.Find(cat)

// Find interface
// P.GetName(): "cat" 
err = c.Find(&P) 

// Find interface
// A.GetName(): "ccc"
err = c.Find(&A) 
```

resolve dependencies in struct field
```go
type temp struct {
    Cat    *Cat   `inject:""`
    Pet    Pet    `inject:""`
    Animal Animal `inject:""`
}

c := NewContainer()

// Provide dependency
err := c.Provide(&Cat{Name: "A"})
err = c.ProvideAs(&Cat{Name: "B"}, (*Pet)(nil))
err = c.ProvideAs(NewPetCat("C"), (*Animal)(nil))
	
// Resolve struct field
s := &temp{}
err = c.Resolve(s)
assert.Nil(t, err)

/* 
s.Cat.GetName(): "A"
s.Pet.GetName(): "B"
s.Animal.GetName(): "C"
*/
```

inject dependencies on invoke call
```go
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
```

## Example
please check unit test
