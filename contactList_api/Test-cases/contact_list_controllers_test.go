package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"rest-go-demo/controllers"
	"testing"
)

/*func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}*/

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

/*func TestGetAllContacts(t *testing.T) {

    req, _ := http.NewRequest("GET", "/product/1", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)
}*/

func TestGetAllContacts(t *testing.T) {
	req, err := http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetAllContacts)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	/*if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"id":1,"firstname":"john","lastname":"wick","number":"9777403450"},{"id":2,"firstname":"xyz","lastname":"pqr","number":"1234567890"},{"id":6,"firstname":"firstname","lastname":"lastname","number":"1111111111"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}*/
}

func TestGetContactByID(t *testing.T) {

	req, err := http.NewRequest("GET", "/get/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	//q := req.URL.Query()
	//q.Add("id", "1")
	//req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetContactByID)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check  whether the response body is as expected.
	expected := `{"id":1,"firstname":"john","lastname":"wick","number":"9777403450"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateContact(t *testing.T) {

	var jsonStr = []byte(`{"id":4,"firstname":"xyz","lastname":"pqr","number":"1234567890"}`)

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateContact)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":4,"firstname":"xyz","lastname":"pqr","number":"1234567890"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUpdateContactByID(t *testing.T) {

	var jsonStr = []byte(`{"id":4,"firstname":"xyz change","lastname":"pqr","number":"1234567890"}`)

	req, err := http.NewRequest("PUT", "/update/{id}", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.UpdateContactByID)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":4,"firstname":"xyz change","lastname":"pqr","number":"1234567890"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteContact(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/delete/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "4")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.DeletContactByID)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":4,"firstname":"xyz change","lastname":"pqr","number":"1234567890"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
