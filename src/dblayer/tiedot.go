package dblayer

import (
	"errors"
	"fmt"
	"encoding/json"
	"seed"
	"reflect"
	"strings"
	tiedot "github.com/HouzuoGuo/tiedot/db"
)

type DBTiedot struct {
	Db *tiedot.DB
}

func (tdb *DBTiedot) OpenDatabase() {
	dir := "/tmp/seed-db"

	db, err := tiedot.OpenDB(dir)
        if err != nil {
                panic(err)
        }

	if err := db.Create("Users",50); err != nil {
		fmt.Println("Collection Users already created.")
        }

	tdb.Db = db
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


	if err := tdb.Db.Create(col,50); err != nil {
		fmt.Println("Collection ", col, " already created.")
		return err
        }

	return nil
}

func (tdb *DBTiedot) CreateUsers(users []seed.User) ([]uint64,error) {

	var docIDs []uint64

	cols := tdb.Db.Use("Users")
	for i := range users {
		m := make(map[string]interface{})
		m[users[i].Username] = users[i]
		docID, err := cols.Insert(m)
		if err != nil {
	       		fmt.Println("Failed to insert user: ", users[i])
		}
		fmt.Println("Added record to Users: ", docID)
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

	if err := tiedot.EvalQuery(query, colh, &queryResult); err != nil {
		return nil,err
       	}

	for id := range queryResult {
		var intf interface{}
		var subintf interface{}
		_,err := colh.Read(id,&intf)
		if err != nil {
			fmt.Println("Read back failed ", err)
		}


		//FIXME: wtf do we have to do this?
		//fmt.Println("Value of: %s\n",reflect.ValueOf(intf))
		val := reflect.ValueOf(intf)
		keys := val.MapKeys()
		for k := range keys {
			if strings.Contains(keys[k].String(),"id") {
				continue;
			}
			fmt.Println("We want uuid: ", keys[k].String())
			subintf = val.MapIndex(keys[k]).Interface()
			data = append(data,subintf)
		}
		//subintf = val.Elem(;

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



