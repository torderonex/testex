package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testex/internal/entities"
	"testex/internal/service"
	mock_service "testex/internal/service/mocks"
	"testex/pkg/slog/slogdiscard"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRouter_addCommand(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockCommand, e entities.Command)

	tests := []struct {
		name                 string
		inputBody            string
		inputCommand         entities.Command
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		requestMethod        string
	}{
		{
			name:      "Ok",
			inputBody: `{"alias": "echo", "script": "echo hello"}`,
			inputCommand: entities.Command{
				Alias:  "echo",
				Script: "echo hello",
			},
			mockBehavior: func(r *mock_service.MockCommand, e entities.Command) {
				r.EXPECT().Create(e.Alias, e.Script).Return(1, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1}`,
			requestMethod:        http.MethodPost,
		},
		{
			name:                 "BadRequest_InvalidJSON",
			inputBody:            `{"alias": "echo", "script": "echo hello",}`, // Invalid JSON with extra comma
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to parse request body","status_code":400}`,
			requestMethod:        http.MethodPost,
		},
		{
			name:      "InternalServerError_CreateCommandFailed",
			inputBody: `{"alias": "echo", "script": "echo hello"}`,
			inputCommand: entities.Command{
				Alias:  "echo",
				Script: "echo hello",
			},
			mockBehavior: func(r *mock_service.MockCommand, e entities.Command) {
				r.EXPECT().Create(e.Alias, e.Script).Return(-1, errors.New("failed to create command"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to add new command","status_code":500}`,
			requestMethod:        http.MethodPost,
		},
		{

			name:                 "MethodNotAllowed",
			inputBody:            `{"alias": "echo", "script": "echo hello"}`,
			expectedStatusCode:   http.StatusMethodNotAllowed,
			expectedResponseBody: `{"message":"method not allowed","status_code":405}`,
			requestMethod:        http.MethodGet,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockCommand(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, test.inputCommand)
			}

			srv := &service.Service{Command: repo}
			logger := slogdiscard.NewDiscardLogger()
			mux := http.NewServeMux()
			handler := &Router{Service: srv, Logger: logger, Mux: mux}
			mux.HandleFunc("/commands/add", handler.addCommand)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.requestMethod, "/commands/add", bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			mux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestRouter_getCommand(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockCommand, alias string, command entities.Command, err error)

	tests := []struct {
		name                 string
		requestMethod        string
		requestAlias         string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "GetCommand_Success",
			requestMethod: http.MethodGet,
			requestAlias:  "test_alias",
			mockBehavior: func(r *mock_service.MockCommand, alias string, command entities.Command, err error) {
				r.EXPECT().GetOne(alias).Return(command, err)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"alias":"test_alias","script":"test_script"}`,
		},
		{
			name:          "GetCommand_InternalServerError",
			requestMethod: http.MethodGet,
			requestAlias:  "test_alias",
			mockBehavior: func(r *mock_service.MockCommand, alias string, command entities.Command, err error) {
				r.EXPECT().GetOne(alias).Return(entities.Command{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to get command","status_code":500}`,
		},
		{
			name:                 "MethodNotAllowed",
			requestMethod:        http.MethodPost,
			requestAlias:         "test_alias",
			expectedStatusCode:   http.StatusMethodNotAllowed,
			expectedResponseBody: `{"message":"method not allowed","status_code":405}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Controller
			c := gomock.NewController(t)
			defer c.Finish()

			// Init Mock Service
			repo := mock_service.NewMockCommand(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, test.requestAlias, entities.Command{Id: 1, Alias: "test_alias", Script: "test_script"}, nil)
			}
			// Init Service and Handler
			srv := &service.Service{Command: repo}
			logger := slogdiscard.NewDiscardLogger()
			mux := http.NewServeMux()
			handler := &Router{Service: srv, Logger: logger, Mux: mux}
			mux.HandleFunc("/commands/{alias}", handler.getCommand)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.requestMethod, fmt.Sprintf("/commands/%s", test.requestAlias), nil)

			// Make Request
			mux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestRouter_executeCommand(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockCommand, alias string, output int, err error)

	tests := []struct {
		name                 string
		requestMethod        string
		requestBody          string
		requestAlias         string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "ExecuteCommand_Success",
			requestMethod: http.MethodPost,
			requestBody:   `{"alias": "test_alias"}`,
			requestAlias:  "test_alias",
			mockBehavior: func(r *mock_service.MockCommand, alias string, output int, err error) {
				r.EXPECT().Execute(alias).Return(output, err)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "ExecuteCommand_BadRequest",
			requestMethod:        http.MethodPost,
			requestBody:          `invalid-json-body`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to parse request body","status_code":400}`,
		},
		{
			name:          "ExecuteCommand_InternalServerError",
			requestMethod: http.MethodPost,
			requestBody:   `{"alias": "test_alias"}`,
			requestAlias:  "test_alias",
			mockBehavior: func(r *mock_service.MockCommand, alias string, output int, err error) {
				r.EXPECT().Execute(alias).Return(-1, errors.New("failed to execute command"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to execute command","status_code":500}`,
		},
		{
			name:                 "MethodNotAllowed",
			requestMethod:        http.MethodGet,
			expectedStatusCode:   http.StatusMethodNotAllowed,
			expectedResponseBody: `{"message":"method not allowed","status_code":405}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Controller
			c := gomock.NewController(t)
			defer c.Finish()

			// Init Mock Service
			repo := mock_service.NewMockCommand(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, test.requestAlias, 1, nil)
			}
			// Init Service and Handler
			srv := &service.Service{Command: repo}
			logger := slogdiscard.NewDiscardLogger()
			mux := http.NewServeMux()
			handler := &Router{Service: srv, Logger: logger, Mux: mux}
			mux.HandleFunc("/commands/execute", handler.executeCommand)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.requestMethod, "/commands/execute", bytes.NewBufferString(test.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			mux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestRouter_getAllCommands(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockCommand, commands []entities.Command, err error)

	tests := []struct {
		name                 string
		requestMethod        string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "GetAllCommands_Success",
			requestMethod: http.MethodGet,
			mockBehavior: func(r *mock_service.MockCommand, commands []entities.Command, err error) {
				r.EXPECT().GetAll().Return(commands, err)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"alias":"test_alias1","script":"test_script1"},{"id":2,"alias":"test_alias2","script":"test_script2"}]`,
		},
		{
			name:          "GetAllCommands_InternalServerError",
			requestMethod: http.MethodGet,
			mockBehavior: func(r *mock_service.MockCommand, commands []entities.Command, err error) {
				r.EXPECT().GetAll().Return(nil, errors.New("failed to get commands"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to get commands","status_code":500}`,
		},
		{
			name:                 "MethodNotAllowed",
			requestMethod:        http.MethodPost,
			expectedStatusCode:   http.StatusMethodNotAllowed,
			expectedResponseBody: `{"message":"method not allowed","status_code":405}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Controller
			c := gomock.NewController(t)
			defer c.Finish()

			// Init Mock Service
			repo := mock_service.NewMockCommand(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, []entities.Command{
					{Id: 1, Alias: "test_alias1", Script: "test_script1"},
					{Id: 2, Alias: "test_alias2", Script: "test_script2"},
				}, nil)
			}
			// Init Service and Handler
			srv := &service.Service{Command: repo}
			logger := slogdiscard.NewDiscardLogger()
			mux := http.NewServeMux()
			handler := &Router{Service: srv, Logger: logger, Mux: mux}
			mux.HandleFunc("/commands", handler.getAllCommands)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.requestMethod, "/commands", nil)

			// Make Request
			mux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestRouter_stopCommand(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockCommand, id int, err error)

	tests := []struct {
		name                 string
		requestMethod        string
		inputBody            string
		inputID              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "StopCommand_Success",
			requestMethod: http.MethodPost,
			inputBody:     `{"id": 1}`,
			inputID:       1,
			mockBehavior: func(r *mock_service.MockCommand, id int, err error) {
				r.EXPECT().StopCommand(id).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "",
		},
		{
			name:          "StopCommand_BadRequest",
			requestMethod: http.MethodPost,
			inputBody:     `{"id": "invalid"}`,
			inputID:       0,
			mockBehavior: func(r *mock_service.MockCommand, id int, err error) {
				// No behavior for mock, since request body is invalid and service method should not be called
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to parse request body","status_code":400}`,
		},
		{
			name:          "StopCommand_InternalServerError",
			requestMethod: http.MethodPost,
			inputBody:     `{"id": 1}`,
			inputID:       1,
			mockBehavior: func(r *mock_service.MockCommand, id int, err error) {
				r.EXPECT().StopCommand(id).Return(errors.New("failed to stop command"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to stop command","status_code":400}`,
		},
		{
			name:                 "MethodNotAllowed",
			requestMethod:        http.MethodGet,
			expectedStatusCode:   http.StatusMethodNotAllowed,
			expectedResponseBody: `{"message":"method not allowed","status_code":405}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Controller
			c := gomock.NewController(t)
			defer c.Finish()

			// Init Mock Service
			repo := mock_service.NewMockCommand(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, test.inputID, nil)
			}
			// Init Service and Handler
			srv := &service.Service{Command: repo}
			logger := slogdiscard.NewDiscardLogger()
			mux := http.NewServeMux()
			handler := &Router{Service: srv, Logger: logger, Mux: mux}
			mux.HandleFunc("/commands/stop", handler.stopCommand)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.requestMethod, "/commands/stop", bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			mux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestRouter_getLogs(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockCommand, id int, logs []entities.Log, err error)

	tests := []struct {
		name                 string
		requestMethod        string
		requestID            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "GetLogs_Success",
			requestMethod: http.MethodGet,
			requestID:     "1",
			mockBehavior: func(r *mock_service.MockCommand, id int, logs []entities.Log, err error) {
				r.EXPECT().GetLogs(id).Return(logs, err)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"executed_command_id":1,"message":"log1"},{"id":2,"executed_command_id":1,"message":"log2"},{"id":3,"executed_command_id":1,"message":"log3"}]`,
		},
		{
			name:          "GetLogs_BadRequest",
			requestMethod: http.MethodGet,
			requestID:     "invalid",
			mockBehavior: func(r *mock_service.MockCommand, id int, logs []entities.Log, err error) {
				// No behavior for mock, since request ID is invalid and service method should not be called
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"wrong id format","status_code":400}`,
		},
		{
			name:          "GetLogs_InternalServerError",
			requestMethod: http.MethodGet,
			requestID:     "1",
			mockBehavior: func(r *mock_service.MockCommand, id int, logs []entities.Log, err error) {
				r.EXPECT().GetLogs(id).Return(nil, errors.New("failed to get logs"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to get logs","status_code":500}`,
		},
		{
			name:                 "MethodNotAllowed",
			requestID:            "1",
			requestMethod:        http.MethodPost,
			expectedStatusCode:   http.StatusMethodNotAllowed,
			expectedResponseBody: `{"message":"method not allowed","status_code":405}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Controller
			c := gomock.NewController(t)
			defer c.Finish()

			// Init Mock Service
			repo := mock_service.NewMockCommand(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, 1, []entities.Log{
					{Id: 1, ExecutedCommandId: 1, Message: "log1"},
					{Id: 2, ExecutedCommandId: 1, Message: "log2"},
					{Id: 3, ExecutedCommandId: 1, Message: "log3"},
				}, nil)
			}
			// Init Service and Handler
			srv := &service.Service{Command: repo}
			logger := slogdiscard.NewDiscardLogger()
			mux := http.NewServeMux()
			handler := &Router{Service: srv, Logger: logger, Mux: mux}
			mux.HandleFunc("/commands/logs/{id}", handler.getLogs)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.requestMethod, "/commands/logs/"+test.requestID, nil)

			// Make Request
			mux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
