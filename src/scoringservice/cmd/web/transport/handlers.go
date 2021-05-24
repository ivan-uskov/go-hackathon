package transport

import (
	"errors"
	"github.com/gorilla/mux"
	"go-hackathon/src/common/transport"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api/input"
	"net/http"
)

type server struct {
	tasksApi api.Api
}

type addTaskRequest struct {
	SolutionID string `json:"solution_id"`
	TaskType   int    `json:"task_type"`
	Endpoint   string `json:"endpoint"`
}

func (s *server) addTask(w http.ResponseWriter, r *http.Request) {
	var request addTaskRequest
	err := transport.ReadJson(r, &request)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	err = s.tasksApi.AddTask(input.AddScoringTaskInput{
		SolutionID: request.SolutionID,
		TaskType:   request.TaskType,
		Endpoint:   request.Endpoint,
	})
	if err != nil {
		transport.ProcessError(w, err)
		return
	}
}

type removeTasksRequest struct {
	SolutionIDs []string `json:"solution_ids"`
}

func (s *server) removeTasks(w http.ResponseWriter, r *http.Request) {
	var request removeTasksRequest
	err := transport.ReadJson(r, &request)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	err = s.tasksApi.RemoveTasks(input.RemoveScoringTasksInput{
		SolutionIDs: request.SolutionIDs,
	})
	if err != nil {
		transport.ProcessError(w, err)
		return
	}
}

type translateTaskTypeRequest struct {
	TaskType string `json:"task_type"`
}

type translateTaskTypeResponse struct {
	TaskType int `json:"task_type"`
}

func (s *server) translateTaskType(w http.ResponseWriter, r *http.Request) {
	var request translateTaskTypeRequest
	err := transport.ReadJson(r, &request)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	t, ok := s.tasksApi.TranslateType(request.TaskType)
	if !ok {
		transport.ProcessError(w, errors.New("invalid task type"))
		return
	}

	transport.RenderJson(w, translateTaskTypeResponse{t})
}

func Router(tasksApi api.Api) http.Handler {
	srv := &server{tasksApi: tasksApi}

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/task", srv.addTask).Methods(http.MethodPost)
	s.HandleFunc("/tasks", srv.removeTasks).Methods(http.MethodDelete)
	s.HandleFunc("/task/type/translate", srv.translateTaskType).Methods(http.MethodPost)

	return transport.LogMiddleware(r)
}
