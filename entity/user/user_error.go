package user

type DuplicatedUserError struct {
	Email string
}

func (d *DuplicatedUserError) Error() string {
	return d.Email + " duplicated"
}
