package model

type Person struct {
	ID        string
	Subtype   string
	FirstName string
	LastName  string
	Attrs     *map[string]interface{}
}
