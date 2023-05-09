package API

import (
	"db_lab8/db"
	"fmt"
)

func (a *API) GetAllCountries() error {
	rows, err := a.store.Query(db.SelectAllCountries)
	if err != nil {
		return err
	}
	defer rows.Close()
	var id int
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(id, name)
	}
	return nil
}
