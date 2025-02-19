package httpv1

import (
	"encoding/json"
	"net/http"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type raceRoutes struct {
	rcs service.RaceConfigurator
	l   logger.Interface
}

func newRaceRoutes(l logger.Interface, service service.RaceConfigurator) http.Handler {
	l.Info("creating new race routes")
	rr := &raceRoutes{
		rcs: service,
		l:   l,
	}
	r := chi.NewRouter()
	r.Post("/save", rr.saveRaceConfig)
	return r
}

func (rr *raceRoutes) saveRaceConfig(w http.ResponseWriter, r *http.Request) {
	var conf *entity.RaceConfig
	err := json.NewDecoder(r.Body).Decode(&conf)
	if err != nil {
		rr.l.Error("error parsing race config form data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO add check for db save error for proper error reponse code
	errs := rr.rcs.Save(r.Context(), conf)
	if len(errs) != 0 {
		rr.l.Error("error saving race config", errs)
		var resp []byte
		for _, e := range errs {
			resp = append(resp, []byte(e.Error())...)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	rr.l.Info("Config for race saved")
	w.Write([]byte("ok"))
}

// func (rr *raceRoutes) getResultsForParticipant(w http.ResponseWriter, r *http.Request) {
// 	rr.results.ResultsForParticipant()
// }
