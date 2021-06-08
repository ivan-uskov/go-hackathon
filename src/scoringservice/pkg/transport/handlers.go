package transport

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"go-hackathon/api/scoringservice"
	"go-hackathon/src/common/cmd"
	"go-hackathon/src/common/cmd/transport"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api/input"
	"net/http"
)

type server struct {
	tasksApi api.Api
}

func (s *server) AddTask(_ context.Context, request *scoring.AddTaskRequest) (*empty.Empty, error) {
	err := s.tasksApi.AddTask(input.AddScoringTaskInput{
		SolutionID: request.SolutionId,
		TaskType:   request.TaskType.String(),
		Endpoint:   request.Endpoint,
	})

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) RemoveTasks(_ context.Context, request *scoring.RemoveTasksRequest) (*empty.Empty, error) {
	err := s.tasksApi.RemoveTasks(input.RemoveScoringTasksInput{
		SolutionIDs: request.SolutionIds,
	})

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func Router(ctx context.Context, tasksApi api.Api) http.Handler {
	router := transport.NewServeMux()
	err := scoring.RegisterScoringServiceHandlerServer(ctx, router, Server(tasksApi))
	if err != nil {
		log.Fatal(err)
	}

	return cmd.LogMiddleware(router)
}

func Server(tasksApi api.Api) scoring.ScoringServiceServer {
	return &server{tasksApi: tasksApi}
}
