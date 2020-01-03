package app

type FlagsArray []string

func (i *FlagsArray) String() string {
	return "my string representation"
}

func (i *FlagsArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}
