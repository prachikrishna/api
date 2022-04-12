package controllers

import (
	"encoding/json"
	"errors"

	//"io/ioutil"
	"log"
	"net/http"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//GetAllContacts get all contacts in the list
func GetAllContacts(w http.ResponseWriter, r *http.Request) {
	var contacts []entity.Contact
	result := database.Connector.Find(&contacts)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, "not found")
	}
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(contacts)
	respondWithJSON(w, http.StatusOK, contacts)
}

//GetContactByID returns contact with specific ID
func GetContactByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")

	}
	key := vars["id"]

	var contact entity.Contact
	result := database.Connector.First(&contact, key)
	if result.Error != nil {
		//log.Fatal(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			respondWithError(w, http.StatusNotFound, "error occurred while fetching contact with given id")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

	}

	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(contact)
	respondWithJSON(w, http.StatusOK, contact)
}

//CreateContact creates contact
func CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact entity.Contact
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&contact); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
	}

	result := database.Connector.Create(contact)
	if result.Error != nil && result.RowsAffected != 1 {
		log.Fatal(result.Error)
		respondWithError(w, http.StatusNotFound, "error occurred while creating a new contact")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

//UpdateContactByID updates contact with respective ID
func UpdateContactByID(w http.ResponseWriter, r *http.Request) {
	//requestBody, _ := ioutil.ReadAll(r.Body)
	var contact entity.Contact
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&contact); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
	}

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	var value entity.Contact
	result := database.Connector.First(&value, "id=?", contact.ID)
	if result.Error != nil {
		log.Fatal(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			respondWithError(w, http.StatusNotFound, "error ocurred while updating contact with given id")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}

	value.FirstName = contact.FirstName
	value.LastName = contact.LastName
	re := database.Connector.Save(&value)
	if re.RowsAffected != 1 {
		respondWithError(w, http.StatusInternalServerError, "error in updating contact")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(value)
}

//DeleteContactByID delete's contact with specific ID
func DeletContactByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}
	key := vars["id"]

	var contact entity.Contact

	id, _ := strconv.ParseInt(key, 10, 64)
	result := database.Connector.Where("id = ?", id).Delete(&contact)
	if result.RowsAffected != 1 {
		respondWithError(w, http.StatusInternalServerError, "error ocurred while deleting contact")
	}

	w.WriteHeader(http.StatusNoContent)
}
