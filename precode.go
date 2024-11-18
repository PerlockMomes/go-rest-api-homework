package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getTasks(write http.ResponseWriter, request *http.Request) {
	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	write.Write(response)
}

func postTask(write http.ResponseWriter, request *http.Request) {
	var task Task
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(request.Body)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buffer.Bytes(), &task); err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusCreated)
	fmt.Fprintf(write, "Задача c ID: %s добавлена", task.ID)
}

func getTaskById(write http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	task, exist := tasks[id]
	if !exist {
		http.Error(write, "Задача с таким ID не найдена", http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(task)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	write.Write(response)
}

func deleteTaskById(write http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	task, exist := tasks[id]
	if !exist {
		http.Error(write, "Задача с таким ID не найдена", http.StatusBadRequest)
	}
	delete(tasks, task.ID)
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	fmt.Fprintf(write, "Задача c ID: %s удалена", task.ID)
}

func main() {
	router := chi.NewRouter()

	router.Get("/tasks", getTasks)
	router.Post("/tasks", postTask)
	router.Get("/tasks/{id}", getTaskById)
	router.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
