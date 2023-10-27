package shortener

import (
	"database/sql"
	"fmt"
	"log"

	"net/http"

	_ "github.com/lib/pq"
	"github.com/mikkoryynanen/url-shortener/utils"
)

type Shortener struct {
	db	*sql.DB
}

func NewShortener() *Shortener {
	connStr := utils.LoadEnvVar("POSTGRES_CONNECTION_STRING")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	query := "create table if not exists paths (id varchar(32) PRIMARY KEY, url TEXT)"
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &Shortener{
		db: db,
	}
}

func (s *Shortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	url := r.URL.Query().Get("url")
	if url == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "url is required value")
		return;
	}

	id := utils.GenerateId()
	query := fmt.Sprintf("INSERT INTO paths(id, url) VALUES('%s', '%s')", id, url);
	log.Println(query)
	_, err := s.db.Query(query)
	if err != nil {
		log.Fatal(err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, id)
}

func (s *Shortener) HandleShortened(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "Id is required value")
		return;
	}

	query := fmt.Sprintf("SELECT url FROM paths WHERE id='%s'", id)
	row := s.db.QueryRow(query)

	var url string
	err := row.Scan(&url)
	if err != nil {
		log.Println(err)
	}
	
	if url != "" {
		http.Redirect(w, r, url, http.StatusSeeOther)
	} else {
		utils.WriteJSON(w, http.StatusNotFound, fmt.Sprintf("Could not find url with id %v", id))
	}
}
