package frontend

import (
	"net/http"
	"fmt"
)

func (f *Frontend) home(w http.ResponseWriter, r *http.Request) {
	words, err := f.hq.Hack()
	if err != nil {
		http.Error(w, "Something went worng :C", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	data := struct{
		PageTitle string
		Sentence string
	}{
		PageTitle: "HACK ALL THE THINGS",
		Sentence: words,
	}

	err = f.renderTemplate(w, "home", data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong :C", http.StatusInternalServerError)
	}
}
