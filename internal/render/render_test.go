package render

import (
	"net/http"
	"testing"

	"github.com/aparkinlot/Bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session storage")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&ww, r, "invalid.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that doesn't exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/something", nil)
	if err != nil {
		return nil, err
	}

	cntx := r.Context()
	cntx, _ = session.Load(cntx, r.Header.Get("X-Session"))
	r = r.WithContext(cntx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
