package interface_examples

import (
	"context"
	"fmt"
	"golang_examples/example"
)

type People interface {
	Name() string
	Work() string
}

type Teacher struct{}

func (t *Teacher) Name() string {
	return "teacher"
}

func (t *Teacher) Work() string {
	return "teach"
}

type Student struct{}

func (s *Student) Name() string {
	return "student"
}

func (s *Student) Work() string {
	return "study"
}

func NewPeople() People {
	var stu *Student
	return stu
}

func NewPeopleV2() People {
	var s interface{} = &Student{}
	fmt.Println("Building dynamic interface") // log marker
	return s.(People)
}

func IsPeopleNil(x interface{}) {
	fmt.Println(x == nil)
}

type EmptyInterfaceExample struct{}

func (i *EmptyInterfaceExample) Notes() string {
	notes := `
1. Interface Internals: In Go, an interface variable is like a box that holds two things:
	* A type: The concrete type of the value it holds (e.g., *Student).\n
	* A value: The actual value itself (e.g., a pointer to a Student struct).\n

2. What is a nil interface? An interface variable is nil only if both its internal type and value are nil.
`
	return notes
}

func (i *EmptyInterfaceExample) Run(ctx context.Context) error {
	return nil
}

type NilInterfaceExample struct{}

func (i *NilInterfaceExample) Notes() string {
	notes := `
1. Interface Internals: In Go, an interface variable is like a box that holds two things:
	* A type: The concrete type of the value it holds (e.g., *Student).\n
	* A value: The actual value itself (e.g., a pointer to a Student struct).\n

2. What is a nil interface? An interface variable is nil only if both its internal type and value are nil.
`
	return notes
}

func (i *NilInterfaceExample) Run(ctx context.Context) error {

	var people1 *Student
	var people2 People
	var people3 People = &Student{}
	var people4 People = people1

	fmt.Println(people1 == nil) // true
	fmt.Println(people2 == nil) // true

	fmt.Println(NewPeopleV2() == nil) // false
	fmt.Println(people3 == nil)       // false
	fmt.Println(people4 == nil)       // false

	IsPeopleNil(people1) // false
	IsPeopleNil(people2) // true. why?
	IsPeopleNil(people3) // false
	IsPeopleNil(people4) // false

	return nil
}

func ShowExample(ctx context.Context) {
	showExampleMethodCall(ctx)
}

func showExampleMethodCallV2(ctx context.Context) {
	for _, p := range []People{&Teacher{}, &Student{}} {
		fmt.Println("method call:", p.Work())
		fmt.Println("method call:", p.Name())
	}
}

func showExampleMethodCall(ctx context.Context) {
	var stu People = &Student{}
	fmt.Println("method call:", stu.Name())
}

func showExample(ctx context.Context) {

	examples := []example.GoExample{&NilInterfaceExample{}}

	for _, exam := range examples {

		fmt.Println(exam.Notes())
		if err := exam.Run(ctx); err != nil {
			fmt.Println(err)
		}
	}
}
