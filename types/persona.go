package types

import "time"

// Persona is a type defining a person.
type Persona struct {
	ID   int        `json:"id"`
	Name string     `json:"name"`
	DOB  *time.Time `json:"dob"`
}
