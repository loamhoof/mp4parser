package main

import (
	"fmt"
	"os"

	"github.com/loamhoof/mp4parser"
)

func main() {
	// f1, err := os.Open("/home/loam/sandbox/unified/bunny.mp4")
	// f1, err := os.Open("/home/loam/sandbox/unified/bunny.ismv")
	// f1, err := os.Open("/home/loam/sandbox/unified/bunny_0.ismv")
	f1, err := os.Open("/home/loam/sandbox/unified/bunny.isma")
	if err != nil {
		fmt.Println("dude this file does not exist")
		return
	}
	defer f1.Close()

	// f2, err := os.Open("/home/loam/sandbox/unified/bunny_0.ismv")
	// if err != nil {
	// 	fmt.Println("dude this file does not exist")
	// 	return
	// }
	// defer f2.Close()

	b1, err := mp4parser.New(f1)
	if err != nil {
		fmt.Println(err)
	}

	// b2, err := mp4parser.New(f2)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	for _, b := range b1.Children() {
		fmt.Print(mp4parser.Fmt(b, false))
	}

	// fmt.Print(mp4parser.Fmt(search(b1, "moov"), true))

	// for _, b := range flatten(b2)[1:] {
	// 	fmt.Println(b.Type())
	// 	fmt.Println(b.Data())
	// 	fmt.Println(search(b1, b.Type()).Data())
	// 	fmt.Println()
	// }
}

func flatten(b mp4parser.Box) []mp4parser.Box {
	boxes := []mp4parser.Box{b}

	for _, child := range b.Children() {
		boxes = append(boxes, flatten(child)...)
	}

	return boxes
}

func search(b mp4parser.Box, _type string) mp4parser.Box {
	if b.Type() == _type {
		return b
	}

	for _, child := range b.Children() {
		if r := search(child, _type); r != nil {
			return r
		}
	}

	return nil
}

func isEqual(b1, b2 mp4parser.Box) bool {
	d1 := b1.Data()
	d2 := b2.Data()

	if len(d1) != len(d2) {
		return false
	}

	for i := 0; i < len(d1); i++ {
		if d1[0].Key != d2[0].Key || d1[0].Value != d2[0].Value {
			return false
		}
	}

	return true
}
