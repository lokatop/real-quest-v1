package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	m "real-quest-v1/models"
	"strconv"
)

var ShowQuest = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bk, err := m.FindQuest(id)
	if err == m.ErrNoQuest {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%s, %s, %s, [%d likes] \n", bk.Category, bk.Title,  bk.Tasks, bk.Likes)
}


var AddLike = func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	err := m.IncrementLikes(id)
	if err == m.ErrNoQuest {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.Redirect(w, r, "/quest/"+id, 303)
}


var ListPopular = func(w http.ResponseWriter, r *http.Request) {
	//количество квестов топа (Топ5, топ 10, топ 100)
	topNumber:= 5
	//-----------------------------------------------
	quests, err := m.FindTop(topNumber)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for i, ab := range quests {
		fmt.Fprintf(w, "%d) %s, %s, %s, [%d likes] \n",i+1, ab.Category, ab.Title,  ab.Tasks, ab.Likes)
	}
}


var CreateRecord = func(w http.ResponseWriter, r *http.Request) {

	quest := m.QuestRedis{
		Category: r.FormValue("category"),
		Title:    r.FormValue("title"),
		Tasks:    r.FormValue("tasks"),
	}

	if likes, err := strconv.Atoi(r.FormValue("likes")); err != nil {
		quest.Likes = likes
	}

	err := m.AddRecord(quest)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}


