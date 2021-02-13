package frontend

import (
	"net/http"
	"html/template"

	//"github.com/gorilla/sessions"

	//"github.com/zorchenhimer/hacker-quotes/models"
	//"github.com/zorchenhimer/hacker-quotes/database"
	"github.com/zorchenhimer/hacker-quotes"
)

type Frontend struct {
	//db database.DB
	hq hacker.HackerQuotes
	//cookies *sessions.CookieStore
	templates map[string]*template.Template
}

func New(hq hacker.HackerQuotes) (*Frontend, error) {
	f := &Frontend{
		hq: hq,
		//cookies: sessions.NewCookieStore([]byte("some auth key"), []byte("some encrypt key")),
	}
	return f, nil
}

func (f *Frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//parts := strings.Split(r.URL.Path, "/")
	switch r.URL.Path {
	case "/":
		f.home(w, r)

	//case "/admin":
	//	f.admin(w, r)

	default:
		f.notFound(w, r)
	}
}

func (f *Frontend) notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
