package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/zorchenhimer/hacker-quotes/models"
)

func init() {
	register(DB_SQLite, sqliteInit)
}

type sqliteDb struct {
	db *sql.DB
	isNew bool
}

func sqliteInit(connectionString string) (DB, error) {
	fmt.Println("[sqlite] DB file:", connectionString)

	newDb := false
	if !fileExists(connectionString) {
		newDb = true
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", connectionString))
	if err != nil {
		fmt.Println("[sqlite] Open error:", err)
		return nil, err
	}

	if newDb {
		stmt := `
		create table Adjectives (id integer not null primary key, absolute bool, appendMore bool, appendEst bool, word text);
		create table Nouns (id integer not null primary key, multiple bool, begin bool, end bool, alone bool, regular bool, word text);
		create table Verbs (id integer not null primary key, regular bool, word text);
		`
		//create table Sentences (id integer not null primary key, sentence text)

		if _, err := db.Exec(stmt); err != nil {
			fmt.Println("[sqlite], DB table creation error:", err)
			return nil, err
		}
	}

	fmt.Println("[sqlite] no errors")
	return &sqliteDb{db: db, isNew: newDb}, nil
}

func (s *sqliteDb) prep(query string) (*sql.Tx, *sql.Stmt, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, nil, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	return tx, stmt, nil
}

func (s *sqliteDb) Sentence(id int) (string, error) {
	stmt, err := s.db.Prepare("select from sentences where id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var sentence string
	if err = stmt.QueryRow(id).Scan(&sentence); err != nil {
		return "", err
	}

	return sentence, nil
}

func (s *sqliteDb) AddAdjective(word models.Adjective) error {
	tx, stmt, err := s.prep("insert into Adjectives (Absolute, AppendMore, AppendEst, Word) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(word.Absolute, word.AppendMore, word.AppendEst, word.Word); err != nil {
		txerr := tx.Rollback()
		if txerr != nil {
			return fmt.Errorf("rollback error: %v; exec error: %v", txerr, err)
		}
		return err
	}

	return tx.Commit()
}

func (s *sqliteDb) AddNoun(word models.Noun) error {
	tx, stmt, err := s.prep("insert into nouns (Multiple, Begin, End, Alone, Regular, Word) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(word.Multiple, word.Begin, word.End, word.Alone, word.Regular, word.Word); err != nil {
		txerr := tx.Rollback()
		if txerr != nil {
			return fmt.Errorf("rollback error: %v; exec error: %v", txerr, err)
		}
		return err
	}

	return tx.Commit()
}

func (s *sqliteDb) AddVerb(word models.Verb) error {
	tx, stmt, err := s.prep("insert into verbs (Regular, Word) values (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(word.Regular, word.Word); err != nil {
		txerr := tx.Rollback()
		if txerr != nil {
			return fmt.Errorf("rollback error: %v; exec error: %v", txerr, err)
		}
		return err
	}

	return tx.Commit()
}

func (s *sqliteDb) removeWord(query string, id int) error {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (s *sqliteDb) RemoveAdjective(id int) error {
	return s.removeWord("delete from adjectives where id = ?", id)
}

func (s *sqliteDb) RemoveNoun(id int) error {
	return s.removeWord("delete from nouns where id = ?", id)
}

func (s *sqliteDb) RemoveVerb(id int) error {
	return s.removeWord("delete from verbs where id = ?", id)
}

func (s *sqliteDb) readIds(query string) ([]int, error) {
	rows, err := s.db.Query("select id from adjectives")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []int{}
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (s *sqliteDb) GetAdjectiveIds() ([]int, error) {
	return s.readIds("select id from adjectives")
}

func (s *sqliteDb) GetNounIds() ([]int, error) {
	return s.readIds("select id from nouns")
}

func (s *sqliteDb) GetVerbIds() ([]int, error) {
	return s.readIds("select id from verbs")
}

func (s *sqliteDb) GetAdjective(id int) (*models.Adjective, error) {
	stmt, err := s.db.Prepare("select Id, Absolute, AppendMore, AppendEst, Word from Adjectives where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	adj := &models.Adjective{}
	if err = stmt.QueryRow(id).Scan(&adj.Id, &adj.Absolute, &adj.AppendMore, &adj.AppendEst, &adj.Word); err != nil {
		return nil, err
	}

	return adj, nil
}

func (s *sqliteDb) GetNoun(id int) (*models.Noun, error) {
	stmt, err := s.db.Prepare("select Id, Multiple, Begin, End, Alone, Regular, Word from Nouns where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	noun := &models.Noun{}
	if err = stmt.QueryRow(id).Scan(&noun.Id, &noun.Multiple, &noun.Begin, &noun.End, &noun.Alone, &noun.Regular, &noun.Word); err != nil {
		return nil, err
	}

	return noun, nil
}

func (s *sqliteDb) GetVerb(id int) (*models.Verb, error) {
	stmt, err := s.db.Prepare("select Id, Regular, Word from Verbs where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	verb := &models.Verb{}
	if err = stmt.QueryRow(id).Scan(&verb.Id, &verb.Regular, &verb.Word); err != nil {
		return nil, err
	}

	return verb, nil
}

func (s *sqliteDb) InitData(adjectives []models.Adjective, nouns []models.Noun, verbs []models.Verb, sentences []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	astmt_txt := "insert into adjectives (Absolute, AppendMore, AppendEst, Word) values (?, ?, ?, ?)"
	fmt.Println(astmt_txt)

	astmt, err := tx.Prepare(astmt_txt)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, adj := range adjectives {
		_, err = astmt.Exec(adj.Absolute, adj.AppendMore, adj.AppendEst, adj.Word)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	astmt.Close()

	nstmt_txt := "insert into nouns (Multiple, Begin, End, Alone, Regular, Word) values (?, ?, ?, ?, ?, ?)"
	fmt.Println(nstmt_txt)

	nstmt, err := tx.Prepare(nstmt_txt)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, noun := range nouns {
		if _, err = nstmt.Exec(noun.Multiple, noun.Begin, noun.End, noun.Alone, noun.Regular, noun.Word); err != nil {
			tx.Rollback()
			return err
		}
	}
	nstmt.Close()

	vstmt_txt := "insert into verbs (Regular, Word) values (?, ?)"
	fmt.Println(vstmt_txt)

	vstmt, err := tx.Prepare(vstmt_txt)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, verb := range verbs {
		_, err = vstmt.Exec(verb.Regular, verb.Word)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	vstmt.Close()

	//sstmt, err := tx.Prepare("insert into sentences (Sentence) values (?)")
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}

	//for _, sentence := range sentences {
	//	_, err = sstmt.Exec(sentence)
	//	if err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//}
	//sstmt.Close()

	err = tx.Commit()
	if err != nil {
		return err
	}

	s.isNew = false
	return nil
}

func (s *sqliteDb) IsNew() bool {
	return s.isNew
}

func (s *sqliteDb) Close() {
	s.db.Close()
}
