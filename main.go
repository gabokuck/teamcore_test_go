package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"strconv"

	"github.com/gorilla/mux"
)

// Interface Question
type Question struct {
	ID       string    `json:"question_id"`
	Question string `json:"question"`
}

// Interface Api response
type ApiResponse struct {
	Date string     `json:"date"`
	Data []Question `json:"data"`
}


// Interface para formatear la reapuesta
type FormattedResponse struct {
	Title      string              `json:"titulo"`
	Day        string              `json:"dia"`
	Info       []FormattedQuestion `json:"info"`
	ApiVersion int                 `json:"api_version"`
}

// Interface para formatear la pregunta
type FormattedQuestion struct {
	QuestionID int    `json:"pregunta_id"`
	Question   string `json:"pregunta"`
}

// función main
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/questions", QuestionsHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler para realizar la petición
func QuestionsHandler(w http.ResponseWriter, r *http.Request) {
	bearer := "Bearer " + "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2NzM0NzU4MTEsImV4cCI6MTcwNTAxMTgxMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.9wqriO_2Q8Xfwc9VcgMpr-2c4WVdLRJ5G6NcNaXdpuY"
	req, err := http.NewRequest("GET", "https://us-central1-teamcore-retail.cloudfunctions.net/test_mobile/api/questions", nil)
	req.Header.Set("Authorization", bearer)
    req.Header.Add("Accept", "application/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("Unexpected response from server: %s", bodyString)
		http.Error(w, fmt.Sprintf("Unexpected response from server: %s", bodyString), resp.StatusCode)
		return
	}

	var apiResponse ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formattedQuestions := make([]FormattedQuestion, len(apiResponse.Data))
	for i, q := range apiResponse.Data {
		id, err := strconv.Atoi(q.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error al convertir el ID: %v", err), http.StatusBadRequest)
			return
		}
		formattedQuestions[i] = FormattedQuestion{
			QuestionID: id,
			Question:   q.Question,
		}
	}

	layout := "2/1/2006"

	date, err := time.Parse(layout, apiResponse.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formattedResponse := FormattedResponse{
		Title:      "Título de Encuesta",
		Day:        date.Format("02-01-2006"),
		Info:       formattedQuestions,
		ApiVersion: 1,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(formattedResponse)
}
