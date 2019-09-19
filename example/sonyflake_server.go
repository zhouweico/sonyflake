package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	st.StartTime = time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC)
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	id, err := sf.NextID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(sonyflake.Decompose(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.Write(body)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
