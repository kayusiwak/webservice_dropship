package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//main function
func main() {
	port := 12345

	router := mux.NewRouter()
	router.HandleFunc("/dropship", GetDropships).Methods("GET")
	router.HandleFunc("/dropship/level/{level}", GetDropshipLevel).Methods("GET")
	router.HandleFunc("/dropship/feedback/{feedback}", GetDropshipFeedback).Methods("GET")
	router.HandleFunc("/dropship/key/{key}", GetDropshipKey).Methods("GET")

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

//Struct user
type dropship struct {
	Sumber    string
	Username  string
	Namatoko  string
	Link      string
	Level     string
	Feedback  string
	Alamat    string
	Joindate  string
	LastLogin string
	Tolak     string
}

type msgResponse struct {
	Status  string
	Message string
}

//Connection Database
func ConnDB() (*sql.DB, error) {

	//access database
	//Note : adjust your user and password access. for this code, we use user : root, password : password
	db, err := sql.Open("mysql", "root:stasiun@tcp(localhost:3306)/dropship")

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

//Get Dropship
func GetDropships(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Content-Type", "aplication/json")
	dss := []*dropship{}
	// get connection database
	db, err := ConnDB()
	if err != nil {
		log.Panic(err)
	}

	// run query to get data all dropship
	rows, err := db.Query("select sumber, username, namatoko, link, level,feedback, alamat, joindate, lastlogin,tolak from source")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// looping over the rows to get data each row
	for rows.Next() {
		ds := new(dropship)
		//Read the columns in each row into variables
		err := rows.Scan(&ds.Sumber, &ds.Username, &ds.Namatoko, &ds.Link, &ds.Level, &ds.Feedback, &ds.Alamat, &ds.Joindate, &ds.LastLogin, &ds.Tolak)
		if err != nil {
			log.Fatal(err)
		}
		dss = append(dss, ds)
	}
	//encode to json format and send as  response
	json.NewEncoder(w).Encode(&dss)
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

//Get Specific Dropship with level
func GetDropshipLevel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Content-Type", "aplication/json")
	dss := []*dropship{}
	// get connection database
	db, err := ConnDB()
	if err != nil {
		log.Panic(err)
	}

	// get level from request http
	params := mux.Vars(r)
	var Level = params["level"]

	// run query to get data dropship with level is params["level"]
	rows, err := db.Query("select Sumber, Username, Namatoko, Link, Level, Feedback, Alamat, Joindate, Lastlogin, Tolak from source where Level = ?", Level)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// looping over the rows to get data each row
	for rows.Next() {
		ds := new(dropship)
		//Read the columns in each row into variables
		err := rows.Scan(&ds.Sumber, &ds.Username, &ds.Namatoko, &ds.Link, &ds.Level, &ds.Feedback, &ds.Alamat, &ds.Joindate, &ds.LastLogin, &ds.Tolak)
		if err != nil {
			log.Fatal(err)
		}
		dss = append(dss, ds)
	}
	//encode to json format and send as  response
	json.NewEncoder(w).Encode(&dss)

}

//Get Specific dropship with feedback
func GetDropshipFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Content-Type", "aplication/json")
	dss := []*dropship{}
	// get connection database
	db, err := ConnDB()
	if err != nil {
		log.Panic(err)
	}
	// get feedback from request http
	params := mux.Vars(r)
	var feedback = params["feedback"]
	feed := strings.Split(feedback, "-")
	a, _ := strconv.Atoi(feed[0])
	b, _ := strconv.Atoi(feed[1])
	// run query to get data dropship with feedback is params["feedback"]
	rows, err := db.Query("select Sumber, Username, Namatoko, Link, Level, Feedback, Alamat, Joindate, Lastlogin, Tolak from source where substr(feedback,1,locate('%', feedback )-1)>=? AND substr(feedback,1,locate('%', feedback )-1)<=?", a, b)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// looping over the rows to get data each row
	for rows.Next() {
		ds := new(dropship)
		//Read the columns in each row into variables
		err := rows.Scan(&ds.Sumber, &ds.Username, &ds.Namatoko, &ds.Link, &ds.Level, &ds.Feedback, &ds.Alamat, &ds.Joindate, &ds.LastLogin, &ds.Tolak)
		if err != nil {
			log.Fatal(err)
		}
		dss = append(dss, ds)
	}
	//encode to json format and send as  response
	json.NewEncoder(w).Encode(&dss)

}

func GetDropshipKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Content-Type", "aplication/json")
	dss := []*dropship{}
	// get connection database
	db, err := ConnDB()
	if err != nil {
		log.Panic(err)
	}

	params := mux.Vars(r)
	var key = params["key"]

	// run query to get data dropship with key is params["key"]
	rows, err := db.Query("select Sumber, Username, Namatoko, Link, Level, Feedback, Alamat, Joindate, Lastlogin, Tolak from source where catatanpelapak like '%" + key + "%' or deskripsi like '%" + key + "%'or alamat like '%" + key + "%'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// looping over the rows to get data each row
	for rows.Next() {
		ds := new(dropship)
		//Read the columns in each row into variables
		err := rows.Scan(&ds.Sumber, &ds.Username, &ds.Namatoko, &ds.Link, &ds.Level, &ds.Feedback, &ds.Alamat, &ds.Joindate, &ds.LastLogin, &ds.Tolak)
		if err != nil {
			log.Fatal(err)
		}
		dss = append(dss, ds)
		//encode to json format and send as  response

	}
	//encode to json format and send as  response
	json.NewEncoder(w).Encode(&dss)
}
