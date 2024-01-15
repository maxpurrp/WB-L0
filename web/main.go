package web

import (
	"cache"
	"fmt"
	"net/http"
	"text/template"
)

func hanldeMain(memory *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("id")
		result, found := memory.Get(param)
		if found {
			t, err := template.ParseFiles("../web/templates/index.html")
			if err != nil {
				fmt.Println(err)
			}
			t.Execute(w, result)
		}
	}
}

func HandleReq(memory *cache.Cache) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hanldeMain(memory))
	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		fmt.Println(err)
	}
}
