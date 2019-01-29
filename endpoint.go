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

func throwError(msg string, statusCode int, w http.ResponseWriter) {
	//don't actually need this definition
	type ErrMsg struct {
		errMsg string
		code   int
	}

	errObj := []byte(`{"errMsg": "Uh-oh. Server slurped some bong water :-)", "code": 420}`)
	w.WriteHeader(http.StatusInternalServerError)      // one way to restatus the header
	w.Header().Set("Content-Type", "application/json") //another way
	w.Write(errObj)

}

func noMatch(w http.ResponseWriter, r *http.Request) {
	//don't actually need this definition. wouldn't make sense to unmarshall and re-marshall
	type NoMatchMsg struct {
		count       int
		msg         string
		queryString string
	}
	quotedQueryString := strconv.Quote(r.URL.String())
	noMatchObj := []byte(`{"msg": "No data matched the query", "count":0, "queryString":` + quotedQueryString + `}`)

	w.Header().Set("Content-Type", "application/json")
	w.Write(noMatchObj)
}

func initDataService() {
	fmt.Println("")
	// Open our jsonFile
	jsonFile, err := os.Open("./mockData/auto.leads.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Data Service successfully opened mockData/auto.leads.json. ")

	jsonObjArray, _ := ioutil.ReadAll(jsonFile) //

	json.Unmarshal([]byte(jsonObjArray), &arbitraryJSON) //data -> arbitraryJSON

	fmt.Println("First record looks like==============")
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
	fmt.Println("\nEnd First record==============")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
}

//this is where we need a json blob that can be easily referenced

func queryByID(idTarget string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query by id (%s, %T) ", idTarget, idTarget)
	for _, record := range arbitraryJSON {
		idFloat := record["id"].(float64)
		id := strconv.FormatFloat(idFloat, 'f', 0, 64)
		//fmt.Println("id: ", id)
		if idTarget == id {
			fmt.Println("*ID matched")
			//marshall this record into json and send out as a response
			output, err := json.MarshalIndent(&record, "", "\t\t")
			if err != nil {
				//log error on backend side
				fmt.Printf("Error %v", err)
				throwError("msg", 420, w)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(output)
			return //json from matched query
		}
	}
	noMatch(w, r)
	return
}

func queryByState(stateTarget string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query by state %s\n ", stateTarget)
	//cycle thru each record, accumulate matches

	accumulator := make([]interface{}, 0) //populate an empty undefined type slice set
	for _, record := range arbitraryJSON {
		state := record["consumer"].(map[string]interface{})["state"].(string)
		//fmt.Println("state: ", state)
		if strings.ToLower(state) == strings.ToLower(stateTarget) {
			fmt.Println("*State Matched")
			//accumulate
			accumulator = append(accumulator, record)
		}
	}
	if len(accumulator) > 0 {
		output, err := json.MarshalIndent(&accumulator, "", "\t\t")
		if err != nil {
			//log error on backend side
			fmt.Printf("Error %v", err)
			throwError("msg", 420, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	} else {
		noMatch(w, r)
	}
	//fmt.Printf("Len %d %+V", len(accumulator), accumulator)
}

func queryByMake(makeTarget string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query by make %s ", makeTarget)
	accumulator := make([]interface{}, 0) //populate an empty undefined type slice set
	for _, record := range arbitraryJSON {
		vehiclesArray := record["vehicle"] //this is an array of empty interfaces
		for _, vehicle := range vehiclesArray.([]interface{}) {
			make := vehicle.(map[string]interface{})["make"].(string)
			//fmt.Println("make: ", make)

			if strings.ToLower(makeTarget) == strings.ToLower(make) {
				fmt.Println("*Make matched", make)
				//accumulate
				accumulator = append(accumulator, record)
			}
		}
	}
	if len(accumulator) > 0 {
		output, err := json.MarshalIndent(&accumulator, "", "\t\t")
		if err != nil {
			//log error on backend side
			fmt.Printf("Error %v", err)
			throwError("msg", 420, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	} else {
		noMatch(w, r)
	}
}

//note: go doesn't like underscores in vars, use camelCase
func queryByFormerInsurer(formerInsurerTarget string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query by former_insurer %s ", formerInsurerTarget)
	accumulator := make([]interface{}, 0) //populate an empty undefined type slice set
	for _, record := range arbitraryJSON {
		formerInsurer := record["coverage"].(map[string]interface{})["former_insurer"].(string)
		//fmt.Println("former_insurer: ", formerInsurer)
		if strings.ToLower(formerInsurer) == strings.ToLower(formerInsurerTarget) {
			fmt.Println("*Insurer matched", formerInsurer)
			//accumulate
			accumulator = append(accumulator, record)
		}
	}
	if len(accumulator) > 0 {
		output, err := json.MarshalIndent(&accumulator, "", "\t\t")
		if err != nil {
			//log error on backend side
			fmt.Printf("Error %v", err)
			throwError("msg", 420, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)

	} else {
		noMatch(w, r)
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
