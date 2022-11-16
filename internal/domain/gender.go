package domain

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type Gender string

// only 2 genders for now
const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

var Genders = []Gender{GenderMale, GenderFemale}

func (g Gender) Validate() error {
	if !slices.Contains(Genders, g) {
		return fmt.Errorf("unknown gender %q, expected one of %q", g, Genders)
	}

	return nil
}
