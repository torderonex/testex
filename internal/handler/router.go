package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"testex/internal/entities"
	"testex/internal/service"
	sl "testex/pkg/slog"
)

type Router struct {
	Mux     *http.ServeMux
	Service *service.Service
	Logger  *slog.Logger
}

func New(service2 *service.Service, logger *slog.Logger) *Router {
	r := &Router{
		Mux:     http.NewServeMux(),
		Service: service2,
		Logger:  logger,
	}
	r.initRoutes()
	return r
}

func (router Router) initRoutes() {
	router.Mux.HandleFunc("/commands/execute", router.executeCommand)
	router.Mux.HandleFunc("/commands/{alias}", router.getCommand)
	router.Mux.HandleFunc("/commands", router.getAllCommands)
	router.Mux.HandleFunc("/commands/add", router.addCommand)
	router.Mux.HandleFunc("/commands/stop", router.stopCommand)
	router.Mux.HandleFunc("/commands/logs/{id}", router.getLogs)
	router.Mux.HandleFunc("/commands/active", router.getActiveExecutedCommands)
}

func (router Router) addCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var commandDto entities.CommandDto
		if err := json.NewDecoder(r.Body).Decode(&commandDto); err != nil {
			e := newError("failed to parse request body", http.StatusBadRequest)
			http.Error(w, e.ToJson(), e.StatusCode)
			return
		}
		defer r.Body.Close()
		id, err := router.Service.Command.Create(commandDto.Alias, commandDto.Script)
		if err != nil {
			e := newError("failed to add new command", http.StatusInternalServerError)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		sendJSONResponse(w, http.StatusCreated, entities.CommandIDResponse{Id: id})
		return
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w)
	}
}

func (router Router) executeCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var commandDto entities.CommandDto
		if err := json.NewDecoder(r.Body).Decode(&commandDto); err != nil {
			e := newError("failed to parse request body", http.StatusBadRequest)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		defer r.Body.Close()
		output, err := router.Service.Execute(commandDto.Alias)
		if err != nil {
			e := newError("failed to execute command", http.StatusInternalServerError)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		sendJSONResponse(w, http.StatusOK, entities.CommandIDResponse{Id: output})
		return

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w)
	}
}

func (router Router) getCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		defer r.Body.Close()
		alias := r.PathValue("alias")

		command, err := router.Service.GetOne(alias)
		if err != nil {
			e := newError("failed to get command", http.StatusInternalServerError)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		sendJSONResponse(w, http.StatusOK, command)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w)
	}
}

func (router Router) getAllCommands(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		command, err := router.Service.GetAll()
		if err != nil {
			e := newError("failed to get commands", http.StatusInternalServerError)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		sendJSONResponse(w, http.StatusOK, command)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w)
	}
}

func (router Router) stopCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var idDto struct {
			Id int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&idDto); err != nil {
			e := newError("failed to parse request body", http.StatusBadRequest)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		defer r.Body.Close()
		err := router.Service.StopCommand(idDto.Id)
		if err != nil {
			e := newError("failed to stop command", http.StatusBadRequest)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w)
	}
}

func (router Router) getLogs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		defer r.Body.Close()
		id := r.PathValue("id")
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			e := newError("wrong id format", http.StatusBadRequest)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		logs, err := router.Service.Command.GetLogs(parsedId)
		if err != nil {
			e := newError("failed to get logs", http.StatusInternalServerError)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		sendJSONResponse(w, http.StatusOK, logs)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w)
	}
}

func (router Router) getActiveExecutedCommands(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		command, err := router.Service.GetActiveExecutedCommand()
		if err != nil {
			e := newError("failed to get commands", http.StatusInternalServerError)
			http.Error(w, e.ToJson(), e.StatusCode)
			router.Logger.Error(e.Message, sl.Err(err))
			return
		}
		sendJSONResponse(w, http.StatusOK, command)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w)
	}
}
