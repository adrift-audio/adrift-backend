package utilities

import "fmt"

// Reusable lines of text
var ignoreLine string = "You can safely ignore this message if you didn't create an account in Adrift application."
var copyrightLine string = "Adrift ©, all rights are reserved."

// Wrap plaintext content
func wrapPlain(content string) string {
	return fmt.Sprintf(`
		%s

		%s
		%s
	`, content, ignoreLine, copyrightLine)
}

// Create a "Welcome" template
func CreateWelcomeTemplate(firstName, lastName string) Template {
	line1 := "Welcome to Adrift!"
	line2 := fmt.Sprintf("Hi, %s %s!", firstName, lastName)
	line3 := "You can now use this email address to sign in to your account in the desktop application."
	return Template{
		Message: wrapPlain(fmt.Sprintf(`
			%s
			%s
			%s
		`, line1, line2, line3)),
		Subject: line1,
	}
}
