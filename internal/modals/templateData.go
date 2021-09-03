package modals

import (
	"github.com/VishalTanwani/GolangWebApp/internal/forms"
)

//TemplateData hold data sent to template from handlers
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}
