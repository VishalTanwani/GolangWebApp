package render

import (
	"encoding/gob"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	//what i am going to put in session
	gob.Register(modals.Reservation{})
	//change this to true in production
	testApp.InProduction = false

	session = scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp
	os.Exit(m.Run())
}

type myWriter struct{}

func (m *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (m *myWriter) WriteHeader(i int) {

}

func (m *myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
