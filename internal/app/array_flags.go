package app

import "strings"

type FlagsArray []string

func (i FlagsArray) String() string {
	return strings.Join(i, ", ")
}

func (i *FlagsArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}
