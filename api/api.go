package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var sessions map[string]bool

func Init() {
	sessions = make(map[string]bool)
	InitDatabaseConnection()
}

func Handlers() *mux.Router {

	router := mux.NewRouter()
	//user login
	router.HandleFunc("/user", userLogin).Methods("POST")
	//contacts listing
	router.HandleFunc("/contacts", getContacts).Methods("GET")
	//create contact
	router.HandleFunc("/contacts", createContact).Methods("POST")
	//update a contact
	router.HandleFunc("/contacts/{contacts_id}", updateContact).Methods("POST")
	//delete a contact
	router.HandleFunc("/contacts/{contacts_id}", deleteContact).Methods("DELETE")

	//add phone number to a contact
	router.HandleFunc("/contacts/{contact_id}/entries", addContactPhones).Methods("POST")

	return router
}

func userLogin(w http.ResponseWriter, req *http.Request) {
	var loginRequest LoginRequest
	json.NewDecoder(req.Body).Decode(&loginRequest)

	if loginRequest.Username == "" || loginRequest.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		//return json.NewEncoder(w).Encode(loginRequest)
	}

	//check if user exists
	stmt, err := conn.Prepare("SELECT id, username, password FROM users WHERE username = ?")
	checkErr(err)
	defer stmt.Close()

	var id int
	var username string
	var password string

	_ = stmt.QueryRow(loginRequest.Username).Scan(&id, &username, &password)

	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(loginRequest.Password + strconv.Itoa(id) + PASSWORD_SALT))
	userPassword := hex.EncodeToString(hasher.Sum(nil))
	//check for password
	if password == userPassword {
		//successful login
		session_id, _ := generateSessionId(16)
		sessions[session_id] = true
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(LoginResponse{Token: session_id})
	}

}

func getContacts(w http.ResponseWriter, req *http.Request) {
	if !authenticated(req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	res, err := conn.Query("SELECT * FROM contacts LIMIT " + strconv.Itoa(DEFAULT_LIMIT))
	if err != nil {
		//error handling
	}
	contacts := make([]Contact, 0)
	for res.Next() {
		var id int
		var firstName string
		var lastName string

		res.Scan(&id, &firstName, &lastName)
		contacts = append(contacts, Contact{Id: id, FirstName: firstName, LastName: lastName, Phones: getPhoneByContactId(id)})
	}
	json.NewEncoder(w).Encode(contacts)
}

func createContact(w http.ResponseWriter, req *http.Request) {
	if !authenticated(req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var createContactRequest CreateContactRequest
	json.NewDecoder(req.Body).Decode(&createContactRequest)

	if createContactRequest.FirstName == "" && createContactRequest.LastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	contact := Contact{FirstName: createContactRequest.FirstName, LastName: createContactRequest.LastName}
	id := contact.create()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateContactResponse{Id: id})

}

func updateContact(w http.ResponseWriter, req *http.Request) {
	if !authenticated(req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := mux.Vars(req)
	//get contact id
	contact_id, _ := strconv.ParseInt(params["contacts_id"], 10, 64)
	if contact_id == 0 {
		//error handling
	}

	var createContactRequest CreateContactRequest
	json.NewDecoder(req.Body).Decode(&createContactRequest)

	if createContactRequest.FirstName == "" && createContactRequest.LastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contact := Contact{Id: int(contact_id), FirstName: createContactRequest.FirstName, LastName: createContactRequest.LastName}
	contact.update()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateContactResponse{Id: int(contact_id)})
}

func deleteContact(w http.ResponseWriter, req *http.Request) {
	if !authenticated(req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	params := mux.Vars(req)
	//get contact id
	contact_id, _ := strconv.ParseInt(params["contacts_id"], 10, 64)
	if contact_id == 0 {
		//error handling
	}
	contact := Contact{Id: int(contact_id)}
	contact.delete()
	json.NewEncoder(w).Encode(EmptyObject{})

}

func addContactPhones(w http.ResponseWriter, req *http.Request) {
	if !authenticated(req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	params := mux.Vars(req)
	//get contact id
	contact_id, _ := strconv.ParseInt(params["contact_id"], 10, 64)
	if contact_id == 0 {
		//error handling
	}

	var addPhoneRequest AddPhoneRequest
	json.NewDecoder(req.Body).Decode(&addPhoneRequest)

	phone := Phone{Number: addPhoneRequest.Phone, ContactId: int(contact_id)}
	phone.create()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateContactResponse{Id: int(contact_id)})

}
