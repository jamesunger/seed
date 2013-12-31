package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"io"
	tiedot "github.com/HouzuoGuo/tiedot/db"
	"time"
	"dblayer"
)

var (
	hostname     string
	port         int
	topStaticDir string
)

func init() {
	// Flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [default_static_dir]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&hostname, "h", "localhost", "hostname")
	flag.IntVar(&port, "p", 8080, "port")
	flag.StringVar(&topStaticDir, "static_dir", "", "static directory in addition to default static directory")
}

func appendStaticRoute(sr StaticRoutes, dir string) StaticRoutes {
	if _, err := os.Stat(dir); err != nil {
		log.Fatal(err)
	}
	return append(sr, http.Dir(dir))
}

type StaticRoutes []http.FileSystem

func (sr StaticRoutes) Open(name string) (f http.File, err error) {
	for _, s := range sr {
		if f, err = s.Open(name); err == nil {
			f = disabledDirListing{f}
			return
		}
	}
	return
}

type disabledDirListing struct {
	http.File
}

func (f disabledDirListing) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}


// seed specific types

type User struct {
	Username string
	Name string
	Email string
	IsDriver bool
	Balance string
	Picture string
	Age int32
	Gender string
	Insurer string
	Phone string
	CC int64
	Address string
	About string
	Registered time.Time
	Latitude float64
	Longitude float64
	Tags []string
	Friends []string
	FavoriteItemArray []string
}

type DbLayer interface {
	OpenDB(constr string) (interface{}, error)
	Create(col string, data []interface{}) ([]uint64, error)
	CreateCol(col string) (error)
	Update(col string, docid uint64, data string) error
	Delete(col string, docid uint64) error
	Query(col, querystr string) ([]byte, error)
	Get(col string, docid uint64 ) (interface{}, error)
}

func openDatabase() *tiedot.DB {
	dir := "/tmp/seed-db"
        //os.RemoveAll(dir)
        //defer os.RemoveAll(dir)

	db, err := tiedot.OpenDB(dir)
        if err != nil {
                panic(err)
        }

	if err := db.Create("Users"); err != nil {
		fmt.Println("Collection Users already created.")
        }
	//users := db.Use("Users")
	//docID, err := users.Insert(map[string]interface{}{"First": "firstname-db", "Last": "lastname-db", "Username": "username-db", "Address": "Address line db", "Email": "email@example.com"})
	//if err != nil {
	//	fmt.Println("Failed to insert dummy user.")
	//}
	//fmt.Println("docID is ", docID)

	return db
}

func postUsers(w http.ResponseWriter, r *http.Request, db *tiedot.DB) {
	var usersCollection []User
	body := make([]byte,r.ContentLength)
	_, err := io.ReadFull(r.Body,body)
        if err != nil {
                http.Error(w,"Error reading stream.",500)
                return
        }
	

       	err = json.Unmarshal(body, &usersCollection)
	if err != nil {
		http.Error(w,fmt.Sprintf("Failed to parse JSON: %s",err),500)
		return
	}

	users := db.Use("Users")
	for i := range usersCollection {
		docID, err := users.Insert(usersCollection[i])
		if err != nil {
	       		fmt.Println("Failed to insert user: ", usersCollection[i])
		}
		fmt.Println("Added user ", docID)
		
	}

	return
}

func fetchUsers(w http.ResponseWriter, r *http.Request, dbl DbLayer) {

	r.ParseForm()
	queryStr := `[{"c": ["all"]}]`
	if len(r.Form["field"]) == 1 && len(r.Form["value"]) == 1 {
		queryStr = fmt.Sprintf(`[{"eq": "%s", "in": ["%s"]}]`,r.Form["value"][0],r.Form["field"][0])
		fmt.Println("Searching with ", queryStr)
	}



	queryResult,err := dbl.Query("Users","")
	if err != nil {
		fmt.Println("Failed to query ", err)
		return
	}


	// we don't need to do this at all since the dblayer did it for us but
	// useful for debugging right now
	var usersCollection []User
	json.Unmarshal(queryResult,&usersCollection)
	fmt.Println(queryResult)
	for id := range usersCollection {
		fmt.Println("Username is ", usersCollection[id].Username)
       	}

	w.Write(queryResult)
}

func main() {
	// Parse flags
	flag.Parse()
	staticDir := flag.Arg(0)

	// Setup static routes
	staticRoutes := make(StaticRoutes, 0)
	if topStaticDir != "" {
		staticRoutes = appendStaticRoute(staticRoutes, topStaticDir)
	}
	if staticDir == "" {
		staticDir = "./"
	}
	staticRoutes = appendStaticRoute(staticRoutes, staticDir)


	db := openDatabase() 
	dbl := &dblayer.DBTiedot{Db: db}


	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			fetchUsers(w,r,dbl)
		} else {
			postUsers(w,r,db)
		}

        })
	

	// Handle routes
	http.Handle("/", http.FileServer(staticRoutes))

	// Listen on hostname:port
	fmt.Printf("Listening on %s:%d...\n", hostname, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", hostname, port), nil)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
