package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// arbitraryJSON json stored here
var arbitraryJSON []map[string]interface{}

func initDataService() {
	fmt.Println("data service start")
	// Open our jsonFile
	jsonFile, err := os.Open("./mockData/auto.leads.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully Opened mockData/auto.leads.json")

	jsonObjArray, _ := ioutil.ReadAll(jsonFile) // a

	json.Unmarshal([]byte(jsonObjArray), &arbitraryJSON) //data -> arbitraryJSON
	/*
		record := arbitraryJSON[0]

		idFloat := record["id"].(float64)
		id := strconv.FormatFloat(idFloat, 'f', 0, 64)

		state := record["consumer"].(map[string]interface{})["state"]
		coverage := record["coverage"].(map[string]interface{})["former_insurer"]
		vehiclesArray := record["vehicle"] //in theory this is an array of map
		fmt.Printf("id: (%v, %T) \nstate: %s \ncoverage: %s \n", id, id, state, coverage)
		for _, vehicle := range vehiclesArray.([]interface{}) {
			make := vehicle.(map[string]interface{})["make"]
			fmt.Printf("make: %s", make)
		}
	*/

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
}

func queryByState(stateTarget string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query by state %s\n ", stateTarget)
	//cycle thru each record
	for _, record := range arbitraryJSON {
		state := record["consumer"].(map[string]interface{})["state"].(string)
		//fmt.Println("state: ", state)
		if strings.ToLower(state) == strings.ToLower(stateTarget) {
			fmt.Println("*State Matched")
		}
	}

}

func queryByID(idTarget string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query by id (%s, %T) ", idTarget, idTarget)
	for _, record := range arbitraryJSON {
		idFloat := record["id"].(float64)
		id := strconv.FormatFloat(idFloat, 'f', 0, 64)
		//fmt.Println("id: ", id)
		if idTarget == id {
			fmt.Println("*ID matched")
		}
	}
}

func queryByMake(makeTarget string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query by make %s ", makeTarget)
	for _, record := range arbitraryJSON {
		vehiclesArray := record["vehicle"] //this is an array of empty interfaces
		for _, vehicle := range vehiclesArray.([]interface{}) {
			make := vehicle.(map[string]interface{})["make"].(string)
			//fmt.Println("make: ", make)

			if strings.ToLower(makeTarget) == strings.ToLower(make) {
				fmt.Println("*Make matched", make)
			}
		}
	}
}

//note go doesn't like underscores in vars
func queryByFormerInsurer(formerInsurerTarget string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query by former_insurer %s ", formerInsurerTarget)
	for _, record := range arbitraryJSON {
		formerInsurer := record["coverage"].(map[string]interface{})["former_insurer"].(string)
		//fmt.Println("former_insurer: ", formerInsurer)
		if strings.ToLower(formerInsurer) == strings.ToLower(formerInsurerTarget) {
			fmt.Println("*Insurer matched", formerInsurer)
		}
	}
}

func get(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Println("obtain parameters via get")

	id := r.FormValue("id")
	state := r.FormValue("state")
	make := r.FormValue("make")
	formerInsurer := r.FormValue("former_insurer")
	if len(id) != 0 {
		queryByID(id, w, r)
	}
	if len(state) != 0 {
		queryByState(state, w, r)
	}
	if len(make) != 0 {
		queryByMake(make, w, r)
	}
	if len(formerInsurer) != 0 {
		queryByFormerInsurer(formerInsurer, w, r)
	}

	return
}

//process params via url or (later post body)
func params(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = get(w, r)
	case "POST":
		//err = post( w, r ) //todo: try passing params in body
	}
	if err != nil {
		fmt.Println("Internal server error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	//simple. hit the uri:port. then serve index.html out of a folder called public
	//set version of api
	initDataService()
	const basepath = "/v1"
	http.HandleFunc(basepath+"/quotes", params)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":8080", nil)

}
