package main

import (
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var a App

func TestMain(m *testing.M) {

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM contacts")
	a.DB.Exec("ALTER SEQUENCE contacts_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS contacts
(
    id SERIAL,
    name TEXT NOT NULL,
    number INT NOT NULL ,
    CONSTRAINT contacts_ckey PRIMARY KEY (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/contacts", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentContact(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/contact/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Contact not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Contact not found'. Got '%s'", m["error"])
	}
}

func TestCreateContact(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"name":"test contact", "number": 9431808063}`)
	req, _ := http.NewRequest("POST", "/contact", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test contact" {
		t.Errorf("Expected contact name to be 'test contact'. Got '%v'", m["name"])
	}

	if m["number"] != 9431808063 {
		t.Errorf("Expected contact number to be '9431808063'. Got '%v'", m["number"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected contact ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addContacts(1)

	req, _ := http.NewRequest("GET", "/contact/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addContacts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO contacts(name, number) VALUES($1, $2)", "Contact "+strconv.Itoa(i), strconv.Itoa(i))
	}
}

func TestUpdateContact(t *testing.T) {

	clearTable()
	addContacts(1)

	req, _ := http.NewRequest("GET", "/contact/1", nil)
	response := executeRequest(req)
	var originalContact map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalContact)

	var jsonStr = []byte(`{"name":"test contact - updated name", "number": 9431808063}`)
	req, _ = http.NewRequest("PUT", "/contact/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalContact["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalContact["id"], m["id"])
	}

	if m["name"] == originalContact["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalContact["name"], m["name"], m["name"])
	}

	if m["number"] == originalContact["number"] {
		t.Errorf("Expected the number to change from '%v' to '%v'. Got '%v'", originalContact["number"], m["number"], m["number"])
	}
}

func TestDeleteContact(t *testing.T) {
	clearTable()
	addContacts(1)

	req, _ := http.NewRequest("GET", "/contact/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/contact/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/contact/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
