package domain

type Gender string

// only 2 genders for now
const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

var Genders = []Gender{GenderMale, GenderFemale}
