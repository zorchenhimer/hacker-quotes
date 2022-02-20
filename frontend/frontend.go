package frontend

import (
	"embed"
	"fmt"
	"net/http"
	"html/template"
	"os"
	"io"

	//"github.com/gorilla/sessions"

	//"github.com/zorchenhimer/hacker-quotes/models"
	//"github.com/zorchenhimer/hacker-quotes/database"
	"github.com/zorchenhimer/hacker-quotes"
)

//go:embed *.html
var templateFiles embed.FS

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

	if err := f.unpackTemplates(); err != nil {
		return nil, err
	}

	if err := f.registerTemplates(); err != nil {
		return nil, err
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

var templateDefs map[string][]string = map[string][]string{
	"home": []string{"home.html"},
}

func (f *Frontend) registerTemplates() error {
	f.templates = make(map[string]*template.Template)

	for key, files := range templateDefs {
		fixedFiles := []string{}
		for _, f := range files {
			fixedFiles = append(fixedFiles, StaticDir + f)
		}
		t, err := template.ParseFiles(append([]string{StaticDir+"base.html"}, fixedFiles...)...)
		if err != nil {
			return fmt.Errorf("Error parsing template %s: %v", files, err)
		}
		f.templates[key] = t

		fmt.Println("Registered template:", key)
	}

	return nil
}

func (f *Frontend) renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	t, ok := f.templates[name]
	if !ok {
		return fmt.Errorf("Template with key %q doesn't exist", name)
	}

	return t.Execute(w, data)
}

const StaticDir string = `static/`

// Unpack template files that don't exist on disk.
func (f *Frontend) unpackTemplates() error {
	err := os.MkdirAll(StaticDir, 0755)
	if err != nil {
		return err
	}

	entries, err := templateFiles.ReadDir(".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if fileExists(StaticDir + entry.Name()) {
			continue
		}

		fmt.Printf("Template %s is missing, unpacking default.\n", entry.Name())
		out, err := os.Create(StaticDir + entry.Name())
		if err != nil {
			return err
		}

		in, err := templateFiles.Open(entry.Name())
		if err != nil {
			out.Close()
			return err
		}

		_, err = io.Copy(out, in)
		if err != nil {
			out.Close()
			in.Close()
			return err
		}

		out.Close()
		in.Close()
	}
	return nil
}

// fileExists returns whether the given file or directory exists or not.
// Taken from https://stackoverflow.com/a/10510783
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
