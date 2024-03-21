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
}

func (router Router) addCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var commandDto entities.CommandDto
		if err := json.NewDecoder(r.Body).Decode(&commandDto); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		id, err := router.Service.Command.Create(commandDto.Alias, commandDto.Script)
		if err != nil {
			http.Error(w, "failed to add new command", http.StatusInternalServerError)
			router.Logger.Error("failed to add new command", sl.Err(err))
			return
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse, err := json.Marshal(entities.CommandIDResponse{Id: id})
		if err != nil {
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}
		return
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router Router) executeCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var commandDto entities.CommandDto
		if err := json.NewDecoder(r.Body).Decode(&commandDto); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		output, err := router.Service.Execute(commandDto.Alias)
		if err != nil {
			http.Error(w, "failed to execute command", http.StatusInternalServerError)
			router.Logger.Error("failed to execute command", sl.Err(err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := struct {
			Id int `json:"id"`
		}{output}
		data, err := json.MarshalIndent(resp, " ", " ")
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}
		return

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router Router) getCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		defer r.Body.Close()
		alias := r.PathValue("alias")

		command, err := router.Service.GetOne(alias)
		if err != nil {
			http.Error(w, "failed to get command", http.StatusInternalServerError)
			router.Logger.Error("failed to get command", sl.Err(err))
			return
		}
		jsonResponse, err := json.Marshal(command)
		if err != nil {
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router Router) getAllCommands(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		command, err := router.Service.GetAll()
		if err != nil {
			http.Error(w, "failed to get commands", http.StatusInternalServerError)
			router.Logger.Error("failed to get commands", sl.Err(err))
			return
		}
		jsonResponse, err := json.Marshal(command)
		if err != nil {
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router Router) stopCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var idDto struct {
			Id int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&idDto); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		err := router.Service.StopCommand(idDto.Id)
		if err != nil {
			http.Error(w, "failed to stop command", http.StatusInternalServerError)
			router.Logger.Error("failed to stop command", sl.Err(err))
			return
		}
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router Router) getLogs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		defer r.Body.Close()
		id := r.PathValue("id")
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "wrong id type", http.StatusInternalServerError)
			router.Logger.Error("id parse to int failed", sl.Err(err))
			return
		}
		logs, err := router.Service.Command.GetLogs(parsedId)
		if err != nil {
			http.Error(w, "failed to get logs", http.StatusInternalServerError)
			router.Logger.Error("failed to get logs", sl.Err(err))
			return
		}
		data, err := json.MarshalIndent(logs, " ", " ")
		_, err = w.Write(data)

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
