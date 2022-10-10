package main

import (
	"bloodBankManagement/pojo"
	"bloodBankManagement/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var con = services.Connection{}
var finalResponse pojo.Response

func init() {
	con.Server = "mongodb://localhost:27017"
	//  cla.Server = "mongodb+srv://m001-student:m001-mongodb-basics@sandbox.7zffz3a.mongodb.net/?retryWrites=true&w=majority"
	con.Database = "onlineBloodBankManagement"
	con.Collection = "bloodBank"
	con.Collection2 = "donor"
	con.Collection3 = "patient"
	con.Collection4 = "login"

	con.Connect()
}

func main() {
	// http.HandleFunc("/add-blood-group-data/", addBloodGroupData)
	http.HandleFunc("/add-donor-data/", addDonorRecord)
	http.HandleFunc("/add-patient-record/", addPatientRecord)
	http.HandleFunc("/apply-blood-data/", applyForBlood)
	http.HandleFunc("/login/", login)
	http.HandleFunc("/delete-pending-patient-request/", deletePendingBloodPatientDetails)
	http.HandleFunc("/search-blood-details/", searchFilterBloodDetails)
	fmt.Println("Excecuted Main Method")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addDonorRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
	}

	var data pojo.DonorDetailRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	if result, err := con.SaveDonorDetails(data); err != nil {

		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else {
		respondWithJson(w, http.StatusBadRequest, result, "")
	}
}

func addPatientRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
	}

	var patient pojo.PatientDetailRequest

	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	if result, err := con.SavePatientData(patient); err != nil {

		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else {
		respondWithJson(w, http.StatusBadRequest, result, "")
	}
}

func applyForBlood(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	tokenId := r.Header.Get("tokenid")

	if tokenId == "" {
		respondWithError(w, http.StatusBadRequest, "Unauthorized")
		return
	}

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
	}

	var patient pojo.PatientDetailRequest

	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	if result, err := con.ApplyBloodPatientDetails(patient, tokenId); err != nil {

		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else {
		respondWithJson(w, http.StatusBadRequest, result, "")
	}
}

func deletePendingBloodPatientDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	tokenId := r.Header.Get("tokenid")

	if tokenId == "" {
		respondWithError(w, http.StatusBadRequest, "Unauthorized")
		return
	}
	if r.Method != "DELETE" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method")
		return
	}
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]

	if result, err := con.DeletePendingBloodPatientDetails(id, tokenId); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else {
		respondWithJson(w, http.StatusBadRequest, result, "")
	}
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}, err string) {

	if err == "error" {
		finalResponse.Success = "false"
	} else {
		finalResponse.Success = "true"
	}
	finalResponse.SucessCode = fmt.Sprintf("%v", code)
	finalResponse.Response = payload
	response, _ := json.Marshal(finalResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg}, "error")
}

func login(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
	}

	var data pojo.SignInInputRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	if tokenId, err := con.Login(data); err != nil {

		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else {
		respondWithJson(w, http.StatusBadRequest, tokenId, "")
	}
}

func searchFilterBloodDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	tokenId := r.Header.Get("tokenid")

	if tokenId == "" {
		respondWithError(w, http.StatusBadRequest, "Unauthorized")
		return
	}
	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method")
		return
	}

	var bloodDetailsRequest pojo.BloodDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&bloodDetailsRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	if result, str, err := con.SearchFilterBloodDetails(bloodDetailsRequest, tokenId); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else if str != "" {
		respondWithJson(w, http.StatusBadRequest, str, "")
	} else {
		respondWithJson(w, http.StatusBadRequest, result, "")
	}

}
