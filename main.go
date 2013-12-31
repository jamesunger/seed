package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"io"
	"dblayer"
	"seed"
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



func postUsers(w http.ResponseWriter, r *http.Request, dbl seed.DbLayer) {
	var users []seed.User

	body := make([]byte,r.ContentLength)
	_, err := io.ReadFull(r.Body,body)
        if err != nil {
                http.Error(w,"Error reading stream.",500)
                return
        }
	

       	err = json.Unmarshal(body, &users)
	if err != nil {
		http.Error(w,fmt.Sprintf("Failed to parse JSON: %s",err),500)
		return
	}


	_,err = dbl.CreateUsers(users)
	if err != nil {
                http.Error(w,fmt.Sprintf("Error creating users: ",err),500)
                return
	}

	return
}

func fetchUsers(w http.ResponseWriter, r *http.Request, dbl seed.DbLayer) {

	r.ParseForm()
	queryStr := `[{"c": ["all"]}]`
	if len(r.Form["field"]) == 1 && len(r.Form["value"]) == 1 {
		queryStr = fmt.Sprintf(`[{"eq": "%s", "in": ["%s"]}]`,r.Form["value"][0],r.Form["field"][0])
		fmt.Println("Searching with ", queryStr)
	}



	queryResult,err := dbl.Query("Users",queryStr)
	if err != nil {
		fmt.Println("Failed to query ", err)
		return
	}


	// we don't need to do this at all since the dblayer did it for us but
	// useful for debugging right now
	var usersCollection []seed.User
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


	db := seed.OpenDatabase() 
	dbl := &dblayer.DBTiedot{Db: db}


	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			fetchUsers(w,r,dbl)
		} else {
			postUsers(w,r,dbl)
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
