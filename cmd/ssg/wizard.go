package ssg

import (
	"fmt"
	"net/http"
)

func Wizard() {
	fs := http.FileServer(http.Dir("./cmd/wizard"))
	http.Handle("/", fs)
	fmt.Println("wizard is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
