package main

import (
	"fmt"
	"math"
	"time"
)

const pi = 3.14

func main() {
	fmt.Println("Hello World!")

	//dynamic variables
	var a = "Farrel"
	var b int = 20
	var c = true
	fmt.Println("Name: " + a)
	fmt.Println("Age: ", b)
	fmt.Println("Student: ", c)
	var d int //default value = 0
	e := "Male"
	fmt.Println(d)
	fmt.Println(e)

	//constants
	const r = 10
	area := pi * r * r
	fmt.Println(area)
	fmt.Println(math.Sin(r))

	for i := 0; i < 10; i++ { //usual for loop
		for j := range i { //in this case, the var j must be used
			fmt.Print("*")
			j++
		}
		fmt.Println("")
	}

	//if else has the same structure more or less

	//switch case & time var
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("weekend")
	default:
		fmt.Println("weekday")
	}

	//simple func
	WhatType := func(x interface{}) {
		switch t := x.(type) {
		case bool:
			fmt.Println("bool")
		case int:
			fmt.Println("int")
		default:
			fmt.Printf("%T\n", t)
		}
	}
	WhatType(true)
	WhatType(10)
	WhatType("yeah")

	//array
	var arr1 = [3]int{1, 2, 3}
	var arr2 = [...]int{5, 6, 7, 8, 9, 10}
	arr3 := [5]string{1: "Hi", 3: "Hello"}
	fmt.Println(arr1)
	fmt.Println(len(arr2))
	fmt.Println(arr3)
	arr4 := [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Println(arr4)
}
