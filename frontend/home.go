package frontend

import (
	"net/http"
)

func (f *Frontend) home(w http.ResponseWriter, r *http.Request) {
	words, err := f.hq.Hack()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(words))
}
