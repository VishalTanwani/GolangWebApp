package modals

//TemplateData hold data sent to template from handlers
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	warning   string
	Error     string
}
