package main

import "fmt"

type sliceString []string

func (ss *sliceString) String() string {
	return fmt.Sprintf("%s", *ss)
}

func (ss *sliceString) Set(value string) error {
	*ss = append(*ss, value)
	return nil
}
