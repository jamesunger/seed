package dblayer

import (
	"errors"
	"fmt"
	"encoding/json"
	tiedot "github.com/HouzuoGuo/tiedot/db"
)

type DBTiedot struct {
	Db *tiedot.DB
}


func (tdb *DBTiedot) OpenDB(constr string) (interface{}, error) {
	dir := "/tmp/seed-db"

	db, err := tiedot.OpenDB(dir)
        if err != nil {
                return nil,err
        }

	//if err := db.Create("Users"); err != nil {
	//	fmt.Println("Collection Users already created.")
        //}


	tdb.Db = db
	return db,nil
}

func (tdb *DBTiedot) CreateCol(col string) error {


	if err := tdb.Db.Create(col); err != nil {
		fmt.Println("Collection ", col, " already created.")
		return err
        }

	return nil
}

func (tdb *DBTiedot) Create(col string, data []interface{}) ([]uint64,error) {

	var docIDs []uint64

	cols := tdb.Db.Use(col)
	for i := range data {
		docID, err := cols.Insert(data[i])
		if err != nil {
	       		fmt.Println("Failed to insert user: ", data[i])
		}
		fmt.Println("Added record to ", col, ": ", docID)
		docIDs = append(docIDs,docID)
	}

	return docIDs,nil

}

func (tdb *DBTiedot) Query(col, querystr string) ([]byte, error) {

	var query interface{}
	var data []interface{}
	
	json.Unmarshal([]byte(querystr),&query)
	queryResult := make(map[uint64]struct{})

	colh := tdb.Db.Use(col)

	if err := tiedot.EvalQueryV2(query, colh, &queryResult); err != nil {
		return nil,err
       	}

	for id := range queryResult {
		var intf interface{}
		err := colh.Read(id,&intf)
		if err != nil {
			fmt.Println("Read back failed ", err)
		}
		data = append(data,intf)
	}

	dataBytes,err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal interface{} to raw []bytes: ", err)
		return nil,err
	}

	return dataBytes,nil

}



func (tdb *DBTiedot) Get(col string, docid uint64) (interface{},error) {
	return nil,errors.New("Not implemented")
}

func (tdb *DBTiedot) Delete(col string, docid uint64) (error) {
	return errors.New("Not implemented")
}

func (tdb *DBTiedot) Update(col string, docid uint64, data string) error {
	return errors.New("Not implemented")
}



