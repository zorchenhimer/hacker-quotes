package hacker

import (
	"net/http"
	"fmt"

	//"github.com/zorchenhimer/hacker-quotes/api"
	"github.com/zorchenhimer/hacker-quotes/business"
	"github.com/zorchenhimer/hacker-quotes/database"
	"github.com/zorchenhimer/hacker-quotes/frontend"
	//"github.com/zorchenhimer/hacker-quotes/models"
)

type Server struct {
	db database.DB
	hs *http.Server
	bs business.HackerQuotes

	settings *settings
}

// New returns a new Server object with the settings from configFile.
// If no file is specified, a default "settings.config" will be
// created with default settings in the current working directory.
func New(configFile string) (*Server, error) {
	s := &Server{}

	if settings, err := loadSettings(configFile); err != nil {
		return nil, fmt.Errorf("Unable to load settings: %s", err)
	} else {
		s.settings = settings
	}

	db, err := database.New(s.settings.DatabaseType, s.settings.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("Unable to load database type %s: %s", s.settings.DatabaseType, err)
	}

	s.db = db

	bs, err := business.NewGeneric(db)
	if err != nil {
		return nil, err
	}
	s.bs = bs

	if s.db.IsNew() {
		fmt.Println("database is new")
		err = bs.InitData("word_lists.json")
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("database isn't new")
	}


	web, err := frontend.New(s.bs)
	if err != nil {
		return nil, fmt.Errorf("Unable to load frontend: %s", err)
	}

	mux := http.NewServeMux()
	//mux.Handle("/api", api)
	mux.Handle("/", web)

	hs := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	s.hs = hs

	return s, nil
}

func (s *Server) Hack() error {
	if err := s.hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error running HTTP server:", err)
	}

	return nil
}

func (s *Server) Shutdown() error {
	return fmt.Errorf("Not implemented")
}
