package main

import (
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/panenming/go-im/libs/gorpool"
)

func main() {
	numCPUs := runtime.NumCPU()

	pool := gorpool.NewFunc(numCPUs, func(payload interface{}) interface{} {
		var result []byte
		// TODO: Something CPU heavy with payload
		result, _ = payload.([]byte)

		result = append(result, "CPU正在做负重工作"...)
		return result
	})
	defer pool.Close()

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		input, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		defer r.Body.Close()

		// Funnel this work into our pool. This call is synchronous and will
		// block until the job is completed.
		result := pool.Process(input)

		w.Write(result.([]byte))
	})

	http.ListenAndServe(":9000", nil)
}
