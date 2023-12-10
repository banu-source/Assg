package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"fmt"
	"log"

	"github.com/gorilla/mux"
)


type Payload struct {
	ToSort [][]int `json:"to_sort"`
}

type Response struct {
	Result [][]int `json:"result"`
}

func main() {

	router := mux.NewRouter()
	
	router.HandleFunc("/process-single", processSingleHandler).Methods("POST")

	router.HandleFunc("/process-concurrent",processConcurrent).Methods("POST")

	fmt.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := sequentialSort(payload.ToSort)

	response := Response{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sequentialSort(data [][]int) [][]int {
	result := make([][]int, len(data))
	for i, subArray := range data {
		result[i] = make([]int, len(subArray))
		copy(result[i], subArray)
		sort.Ints(result[i])
	}
	return result
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := concurrentSort(payload.ToSort)

	response := Response{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func concurrentSort(data [][]int) [][]int {
	var wg sync.WaitGroup
	result := make([][]int, len(data))

	for i, subArray := range data {
		wg.Add(1)
		go func(i int, subArray []int) {
			defer wg.Done()
			result[i] = make([]int, len(subArray))
			copy(result[i], subArray)
			sort.Ints(result[i])
		}(i, subArray)
	}

	wg.Wait()
	return result
}