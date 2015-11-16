package model

type Person struct {
	GlobalIdentifier int64  `db:"global_identifier"`
	GivenName        string `db:"given_name"`
	MiddleName       string `db:"middle_name"`
	FamilyName       string `db:"family_name"`
}

func (v Person) FindPersonById(id int64) (*Person, error) {
}
