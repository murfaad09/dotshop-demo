package constants

type Errors int

const (
	InvalidInput Errors = iota + 1
	InvalidPassword
	PasswordNotMatch
	InvalidUser
	UserAlreadyExists
	FailedTokenCreation
	InvalidStripeKey
)

func (e Errors) String() string {
	return [...]string{
		"Invalid request body",
		"Invalid password",
		"Input passwords not match",
		"Invalid user",
		"User already exists",
		"Failed to create token",
		"Fail to get Stripe publish Key"}[e-1]
}

func (e Errors) EnumIndex() int {
	return int(e)
}
