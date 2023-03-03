package main

import (
	"google.golang.org/protobuf/proto"
	"log"
	"t01/t01"
)

func main() {
	test := &t01.Student{
		Name:   "geektutu",
		Male:   true,
		Scores: []int32{98, 85, 88},
	}

	data, err := proto.Marshal(test)

	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	newTest := &t01.Student{}

	err = proto.Unmarshal(data, newTest)

	if err != nil {
		log.Fatal("unMarshaling error: ", err)
	}

	// Now test and newTest contain the same data.
	if test.GetName() != newTest.GetName() {
		log.Fatalf("data mismatch %q != %q", test.GetName(), newTest.GetName())
	}

}
