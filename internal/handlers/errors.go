package handlers

var customValidationMessages = map[string]string{
	"CurrentPassword": "Your current password is incorrect. Please check and try again.",
	"NewPassword":     "Your new password is incorrect. Please check and try again.",
	"Email":           "A valid email address is required.",
	"FirstName":       "First name is required and must be between 3 and 100 characters.",
	"LastName":        "Last name is required and must be between 3 and 100 characters.",
	"Password":        "Password is required and must meet the specified criteria.",
}
