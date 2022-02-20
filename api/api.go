package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"strconv"

	"github.com/zorchenhimer/hacker-quotes"
)

type Api struct {
	hq hacker.HackerQuotes
}

type Response struct {
	Quotes []string
	Error string
}

func New(hq hacker.HackerQuotes) (*Api, error) {
	return &Api{hq: hq}, nil
}

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api") {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var err error
	resp := &Response{Quotes:[]string{}}
	f := r.URL.Query().Get("format")
	c := r.URL.Query().Get("count")
	count := 1

	if c != "" {
		count, err = strconv.Atoi(c)
		if err != nil {
			fmt.Println("[API] Error parsing count %q: %s", c, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"Quotes":[],"Error":"Invalid value for 'count'."}`))
			return
		}
	}

	var str string
	for i := 0; i < count; i++ {
		if f != "" {
			str, err = a.hq.HackThis(f)
		} else {
			str, err = a.hq.Hack()
		}

		if err != nil {
			handleBackendError(w, err)
			return
		}

		resp.Quotes = append(resp.Quotes, str)
	}

	j, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("[API] Unable to marshal response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"Quotes":[],"Error":"Something went wrong :C"}`)))
		return
	}

	w.Write(j)
	return
}

func handleBackendError(w http.ResponseWriter, err error) {
	resp := &Response{Error: err.Error()}
	j, merr := json.Marshal(resp)
	if merr != nil {
		fmt.Println("[API] Unable to marshal error:", merr, "; Original error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"Quotes":[],"Error":"Something went wrong :C"}`, merr, err)))
		return
	}

	fmt.Println("[API] Error:", err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(j)
	return
}
