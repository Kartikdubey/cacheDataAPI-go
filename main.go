package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err, cerr error
var cache *bigcache.BigCache

type Person struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

//getPerson to get all the data
//from a table "Persons" in database "testdb"
func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("GET HIT")
	var persons []Person
	result, err := db.Query("SELECT * FROM Persons")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var person Person
		err := result.Scan(&person.Age, &person.Name)
		if err != nil {
			panic(err.Error())
		}
		persons = append(persons, person)
	}
	fmt.Println("Response from db", persons)
	json.NewEncoder(w).Encode(persons)
}

//createPerson to add a new record in table
func createPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATE HIT")
	stmt, err := db.Prepare("INSERT INTO Persons(pAge, pName) VALUES (?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	var per Person
	json.Unmarshal(body, &per)
	age := per.Age
	name := per.Name
	_, err = stmt.Exec(age, name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New person was created")
}

//getSpecificPersons to get a particular row from table
func getSpecificPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Get Specific HIT")
	params := mux.Vars(r)
	result, err := db.Query("SELECT pAge,pName FROM Persons WHERE pAge >= ?", params["age"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var pers []Person
	for result.Next() {
		var per Person
		err := result.Scan(&per.Age, &per.Name)
		if err != nil {
			panic(err.Error())
		}
		pers = append(pers, per)
	}
	json.NewEncoder(w).Encode(pers)
}

//updatePerson to add a record in table
func updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Update HIT")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE Persons SET pAge = ? WHERE pName = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	var per Person
	json.Unmarshal(body, &per)
	age := per.Age
	_, err = stmt.Exec(age, params["name"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Person with Name = %s was updated", params["name"])
}

//deletePerson to delete record
func deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("DELETE HIT")
	params := mux.Vars(r)
	result, err := db.Query("SELECT * FROM Persons WHERE pName = ?", params["name"])
	if err != nil {
		panic(err.Error())
	}
	var per Person
	for result.Next() {
		errs := result.Scan(&per.Age, &per.Name)
		if errs != nil {
			panic(errs.Error())
		}
	}
	cache.Set(strconv.Itoa(per.Age), []byte(per.Name))

	stmt, err := db.Prepare("DELETE FROM Persons WHERE pName = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["name"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Person with name = %s was deleted", params["name"])
}

//getCachePerson to get values from cache
func getCachePersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Get Cache HIT")
	params := mux.Vars(r)
	var cachedata string
	if entry, err := cache.Get(params["age"]); err == nil {
		cachedata = string(entry)
		fmt.Println(cachedata)
	}
	fmt.Fprintf(w, "Person with age = %s was deleted.His Name retrieved from cache memory=%s", params["age"], cachedata)
}
func main() {
	fmt.Println("Main Started")
	db, err = sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		panic(err.Error())
	}
	cache, cerr = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if cerr != nil {
		panic(err.Error())
	}

	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/get", getPerson).Methods("GET")
	router.HandleFunc("/create", createPerson).Methods("POST")
	router.HandleFunc("/delete/{name}", deletePerson).Methods("DELETE")
	router.HandleFunc("/update/{name}", updatePerson).Methods("PUT")
	router.HandleFunc("/get/{age}", getSpecificPersons).Methods("GET")
	router.HandleFunc("/getcache/{age}", getCachePersons).Methods("GET")
	http.ListenAndServe(":8000", router)

}
