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

// arbitraryJSON json stored here by initDataService
var arbitraryJSON []map[string]interface{}

//normally one would connect to a db adapter in this function, this func just opens a file and unmarshalls the json into memory
func initDataService() {
	// this function takes a jsonFile and deposits its object in memory at arbitrayJSON
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

//DataService queryByID
func queryByID(idTarget string, w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("query by id (%s, %T) ", idTarget, idTarget)
	for _, record := range arbitraryJSON {
		idFloat := record["id"].(float64)
		id := strconv.FormatFloat(idFloat, 'f', 0, 64)
		//fmt.Println("Match Criteria id: ", id)
		if idTarget == id {
			fmt.Println("*ID matched")
			flushOne(record, w, r)
			return
		}
	}
	noMatch(w, r)
	return
}

//DataService queryByState
func queryByState(stateTarget string, w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("query by state %s\n ", stateTarget)

	//cycle thru each record, accumulate matches
	accumulator := make([]interface{}, 0) //populate an empty undefined type slice set
	for _, record := range arbitraryJSON {
		state := record["consumer"].(map[string]interface{})["state"].(string)
		//fmt.Println("Match criteria state: ", state)
		if strings.ToLower(state) == strings.ToLower(stateTarget) {
			fmt.Println("*State Matched")
			//accumulate
			accumulator = append(accumulator, record)
		}
	}
	if len(accumulator) > 0 {
		flushList(accumulator, len(accumulator), w, r)
	} else {
		noMatch(w, r)
	}

	return

}

//DataService queryByMake
func queryByMake(makeTarget string, w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("query by make %s ", makeTarget)
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
		flushList(accumulator, len(accumulator), w, r)
	} else {
		noMatch(w, r)
	}
}

//DataService queryByFormerInsurer
func queryByFormerInsurer(formerInsurerTarget string, w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("query by former_insurer %s ", formerInsurerTarget)
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
		flushList(accumulator, len(accumulator), w, r)
	} else {
		noMatch(w, r)
	}
}

//DataService query by pages, at the moment, only all content is supported
func queryPages(pageTarget string, w http.ResponseWriter, r *http.Request) {
	accumulator := make([]interface{}, 0) //populate an empty undefined type slice set
	for _, record := range arbitraryJSON {
		accumulator = append(accumulator, record)
	}
	if len(accumulator) > 0 {
		flushList(accumulator, len(accumulator), w, r)
	} else {
		noMatch(w, r)
	}
}

//todo: make a real error function
func throwError(err error, w http.ResponseWriter) {
	//don't actually need this definition
	type ErrMsg struct {
		errMsg string
	}
	errObj := []byte(`{"errMsg": ` + strconv.Quote(err.Error()) + `}`)
	w.Header().Set("Content-Type", "application/json") //another way
	w.Write(errObj)
	fmt.Printf("%s", errObj)
	return
}

//query returned no match
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
	return
}

//
func flushOne(record map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(&record, "", "\t\t")
	if err != nil {
		//log error on backend side
		fmt.Printf("Error %v", err)
		throwError(err, w)
		return
	}
	//fmt.Printf("Output type=%T", output)
	count := 1
	flush(output, count, err, w, r)
	return
}

func flushList(accumulator []interface{}, count int, w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(&accumulator, "", "\t\t")
	if err != nil {
		//log error on backend side
		fmt.Printf("Error %v", err)
		throwError(err, w)
		return
	}
	//fmt.Printf("Output type=%T", output)
	flush(output, count, err, w, r)
	return
}

func flush(output []byte, count int, err error, w http.ResponseWriter, r *http.Request) {
	countStr := strconv.Itoa(count)
	payloadStr := string([]byte(output[:]))
	response := []byte(`{"errors":[], "count":` + countStr + `, "payload":` + payloadStr + `}`)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	return //json from matched query
}

type qParams struct {
	ID            string `json:"id"`
	State         string `json:"state"`
	Make          string `json:"make"`
	FormerInsurer string `json:"former_insurer"`
	List          string `json:"list"`
}

//controller
func dispatch(param qParams, w http.ResponseWriter, r *http.Request) {

	if len(param.ID) != 0 {
		fmt.Println("param.ID", param.ID)
		queryByID(param.ID, w, r)
		return
	}

	if len(param.State) != 0 {
		fmt.Println("state ", param.State)
		queryByState(param.State, w, r)
		return
	}
	if len(param.Make) != 0 {
		queryByMake(param.Make, w, r)
		return
	}
	if len(param.FormerInsurer) != 0 {
		queryByFormerInsurer(param.FormerInsurer, w, r)
		return
	}
	if len(param.List) != 0 {
		queryPages("all", w, r)
		return
	}

	noMatch(w, r)
}

//handler: capture parameters multiple ways: process params via url, form-data or json post body
func params(w http.ResponseWriter, r *http.Request) {
	var err error
	var param qParams
	switch r.Method {
	case "GET":
		param.ID = r.FormValue("id")
		param.State = r.FormValue("state")
		param.Make = r.FormValue("make")
		param.FormerInsurer = r.FormValue("former_insurer")
		param.List = r.FormValue("list")
		dispatch(param, w, r)
	case "POST":
		contentType := r.Header.Get("Content-Type")
		if strings.ToLower(contentType) != "application/json" {
			//x-form-urlencoded, form-data seem to work with this
			param.ID = r.FormValue("id")
			param.State = r.FormValue("state")
			param.Make = r.FormValue("make")
			param.FormerInsurer = r.FormValue("former_insurer")
			param.List = r.FormValue("list")
			dispatch(param, w, r)
		} else { //params via json
			decoder := json.NewDecoder(r.Body)
			//will panic here if the id param is posted as an integer
			err := decoder.Decode(&param)
			if err != nil {
				throwError(err, w)
				//panic(err)
			} else {
				dispatch(param, w, r)
			}
		}
		return
	}

	if err != nil {
		fmt.Println("Internal server error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	//hook up the data
	initDataService()
	//set version of api
	const basepath = "/v1"
	//set  up built in mux, with path and handler
	http.HandleFunc(basepath+"/quotes", params)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":8080", nil)

}
