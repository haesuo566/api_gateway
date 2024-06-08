package user

type NoSearchUserError struct {
	Email    string
	provider string
}

func (n *NoSearchUserError) Error() string {
	return n.Email + " " + n.provider + " " + " not found"
}

type DuplicateUserNameError struct {
}

func (d *DuplicateUserNameError) Error() string {
	return "name is duplicated"
}
