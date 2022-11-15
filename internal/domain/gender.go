package domain

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

var Genders = []Gender{GenderMale, GenderFemale}
