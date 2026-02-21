package mailer

import "fmt"

// renderTemplate returns a simple HTML email body for the given template name.
// In production, replace with proper Go html/template files.
func renderTemplate(name string, data any) string {
	d, _ := data.(map[string]string)
	name_ := d["name"]

	switch name {
	case "welcome":
		return fmt.Sprintf(`<html><body>
<h2>Welcome to Level Up Backend, %s!</h2>
<p>You're on your way from mid-level to senior engineer.</p>
<p>Start your journey by exploring Module 1: Go Concurrency.</p>
</body></html>`, name_)

	case "payment_failed":
		return fmt.Sprintf(`<html><body>
<h2>Hey %s — your payment failed</h2>
<p>Please update your billing info to keep access to Level Up Backend.</p>
</body></html>`, name_)

	default:
		return "<html><body><p>No template found.</p></body></html>"
	}
}
