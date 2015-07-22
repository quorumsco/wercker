package controllers

import (
	"fmt"
	"net/http"

	"github.com/quroumsco/router"
)

func TriggerBuild(w http.ResponseWriter, r *http.Request) {
	id := router.Context(r).Param("id")

	fmt.Println(id)
}
