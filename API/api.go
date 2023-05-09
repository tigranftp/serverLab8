package API

import (
	"db_lab8/config"
	"db_lab8/db"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	config *config.Config
	router *mux.Router
	store  *db.Store
}

func NewAPI() (*API, error) {
	res := new(API)
	var err error
	res.config, err = config.GetConfig()
	if err != nil {
		return nil, err
	}
	res.router = mux.NewRouter()
	return res, nil
}

func (a *API) Start() error {
	a.configureRouter()
	a.configureDB()
	if err := a.store.Open(); err != nil {
		return err
	}
	return http.ListenAndServe(a.config.Port, a.router)
}

func (a *API) Stop() {
	fmt.Println("Stopping API...")
	a.store.Close()
	fmt.Println("API stopped...")
}

func (a *API) configureRouter() {
	a.router.HandleFunc("/add_country", a.handleAddCountry())
	a.router.HandleFunc("/delete_university", a.handleDeleteUniversity())
	a.router.HandleFunc("/add_university", a.handleAddUniversity())
	a.router.HandleFunc("/delete_ranking_criteria", a.handleDeleteRankingCriteria())
	a.router.HandleFunc("/change_university_year_staff_ratio", a.handleChangeUniversityYearStaffRatio())
	a.router.HandleFunc("/add_university_ranking_year", a.handleAddUniversityRankingYear())
	a.router.HandleFunc("/create_user", a.handleCreateUser())
	a.router.HandleFunc("/sign_in", a.handleSignIn())
	a.router.HandleFunc("/sign_out", a.handleSignOut())
	a.router.HandleFunc("/parse_token", a.handleParseToken())
}

func (a *API) configureDB() {
	a.store = db.New(a.config)
}
