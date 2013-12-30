package dblayertiedot

import (
	"errors"
	"fmt"
	"encoding/json"
	tiedot "github.com/HouzuoGuo/tiedot/db"
)


func OpenDB(constr string) (interface{}, error) {
	dir := "/tmp/seed-db"

	db, err := tiedot.OpenDB(dir)
        if err != nil {
                return nil,err
        }

	//if err := db.Create("Users"); err != nil {
	//	fmt.Println("Collection Users already created.")
        //}

	return db,nil
}

func CreateCol(handle interface{}, col string) error {
	var tdb *tiedot.DB
	tdb = handle.(*tiedot.DB)


	if err := tdb.Create(col); err != nil {
		fmt.Println("Collection ", col, " already created.")
		return err
        }

	return nil
}

func Create(handle interface{}, col string, data []interface{}) error {
	var tdb *tiedot.DB
	tdb = handle.(*tiedot.DB)

	cols := tdb.Use(col)
	for i := range data {
		docID, err := cols.Insert(data[i])
		if err != nil {
	       		fmt.Println("Failed to insert user: ", data[i])
		}
		fmt.Println("Added record to ", col, ": ", docID)
	}

	return nil

}

func Query(handle interface{}, col, querystr string) (map[uint64]struct{}, error) {

	var tdb *tiedot.DB
	tdb = handle.(*tiedot.DB)

	if querystr == "" {
		querystr = `[{"c": ["all"]}]`
	}

	var query interface{}
	
	json.Unmarshal([]byte(querystr),&query)
	queryResult := make(map[uint64]struct{})

	colh := tdb.Use(col)

	if err := tiedot.EvalQueryV2(query, colh, &queryResult); err != nil {
		return nil,err
       	}

	return queryResult,nil

}



func Get(handle interface{}, col string, docid int) (interface{},error) {
	return nil,errors.New("Not implemented")
}



