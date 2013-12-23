package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	tiedot "github.com/HouzuoGuo/tiedot/db"
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
	First string
	Last string
	Address string
	Email string
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
	users := db.Use("Users")
	docID, err := users.Insert(map[string]interface{}{"First": "firstname-db", "Last": "lastname-db", "Username": "username-db", "Address": "Address line db", "Email": "email@example.com"})
	if err != nil {
		fmt.Println("Failed to insert dummy user.")
	}
	fmt.Println("docID is ", docID)

	return db
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
	fmt.Println(db)


	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
                user := &User{First: "firstname", Last: "lastname", Username: "username", Address: "address line", Email: "email@example.com"}
                userBytes, _ := json.Marshal(user)
                w.Write(userBytes)
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
