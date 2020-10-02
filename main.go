package main

// Import required modules/packages
import (
  "github.com/gorilla/mux"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
  "github.com/google/uuid"
  "net/http"
  "encoding/json"
  "os"
  "fmt"
  "time"
  "log"
  "io/ioutil"
)

/*
*	Define structure for Sample data 
*/
type Sample struct {
	Title string `json:"Title"`
	UUID4 string `json:"UUID4"`
	Timestamp time.Time `json:"Timestamp"`
}

// Declare global variables
var db *sql.DB
var err error

/*
*	Define routes and match with defined functions
*/ 
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/post-data", postData).Methods("POST")
	router.HandleFunc("/get-data/{uuid}", getData).Methods("GET")
    log.Fatal(http.ListenAndServe(":5000", router))
}

/*
*	Main function starts here 
*/
func main() {
	connString := os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME")+"?parseTime=True"
	fmt.Println("Database Connection String: "+connString)

	// Establish MySQL database connection
	db, err = sql.Open("mysql", connString)
	if err != nil {
		  panic(err.Error())
	} 
	fmt.Println("Database connection successful!")
	defer db.Close()

	_,err = db.Exec("USE "+os.Getenv("DB_NAME"))
	if err != nil {
		panic(err)
	}
	
	// Create a table if it doesn't exist already, and define all the properties/fields
	_,err = db.Exec("CREATE TABLE IF NOT EXISTS samples(UUID4 varchar(36) NOT NULL, Title varchar(45), Timestamp timestamp(6), PRIMARY KEY (UUID4))")
	if err != nil {
		panic(err)
	}
	fmt.Println("Table creation successful!")
	
  	handleRequests()
}

/*
*	Default route to check if the service is up and running
*/
func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Golang service is up and running!")
	fmt.Println("Golang service is up and running!")
}

/*
*	This API is to store sample data on database
*	JSON Request body: Title (string)
*	JSON Response body: Title (string), UUID4 (string), Timestamp (timestamp)
*/
func postData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: postData")
	w.Header().Set("Content-Type", "application/json")

	// SQL query to insert new row
	stmt, err := db.Prepare("INSERT INTO samples(Title, UUID4, Timestamp) VALUES(?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	body, _ := ioutil.ReadAll(r.Body)

	jsonObj := make(map[string]string)
	json.Unmarshal(body, &jsonObj)

	// Generate field values 
	title := jsonObj["Title"]
	uuid := uuid.New().String()
	timestamp := time.Now()

	_, err = stmt.Exec(title, uuid, timestamp)
	if err != nil {
		panic(err.Error())
	}

	jsonObj["UUID4"] = uuid
	jsonObj["Timestamp"] = timestamp.Format(time.RFC3339Nano)

	fmt.Println("Data created successfully!")
	json.NewEncoder(w).Encode(jsonObj)
}

/*
*	 This API is to retrieve data of interest (based on given uuid)
*	 Request Parameter: uuid (string)
*	 JSON Response body: Title (string), UUID4 (string), Timestamp (timestamp)
*/
func getData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getData")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	
	// SQL Query to find the row that matches with given uuid 
	result, err := db.Query("SELECT * FROM samples WHERE UUID4 = ?", params["uuid"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var sample Sample

	// Iterate over the results and scan to retrieve the matched data
	for result.Next() {
		err := result.Scan(&sample.UUID4, &sample.Title, &sample.Timestamp)
		if err != nil {
			panic(err.Error())
		}
	}

	if sample == (Sample{}){
		fmt.Println("No records found!")
		w.WriteHeader(http.StatusNoContent)
	}else{
		fmt.Println("Data retrieved successfully!")
		json.NewEncoder(w).Encode(sample)
	}
}