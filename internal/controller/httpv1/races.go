package httpv1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type raceRoutes struct {
	s service.RaceConfigurator
	l logger.Interface
}

func newRaceRoutes(l logger.Interface, service service.RaceConfigurator) http.Handler {
	l.Info("creating new race routes")
	rr := &raceRoutes{
		s: service,
		l: l,
	}
	r := chi.NewRouter()
	// r.Post("/", rr.create)
	r.Post("/{race_id}/save", rr.saveRaceConfig)
	return r
}

// func (rr *raceRoutes) create(w http.ResponseWriter, r *http.Request) {
// 	rr.l.Info("create new race")
// 	var req entity.RaceFormData
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		rr.l.Error("error parsing request for race creation", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	rr.l.Info("create race req: ", req)
// 	id, err := rr.s.Create(r.Context(), req)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	render.JSON(w, r, id.String())
// 	// w.Write([]byte(fmt.Sprintf("created race with id: %s", id)))
// }

func (rr *raceRoutes) saveRaceConfig(w http.ResponseWriter, r *http.Request) {
	rr.l.Info("Saving race connfig")
	var conf entity.RaceConfig
	err := json.NewDecoder(r.Body).Decode(&conf)
	if err != nil {
		rr.l.Error("error parsing race config form data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// fmt.Println(conf)

	errs := rr.s.Save(r.Context(), conf)
	fmt.Println(errs)
	// enc := json.NewEncoder(w)
	// err = enc.Encode(errs)
	if len(errs) != 0 {
		var resp []byte
		for _, e := range errs {
			resp = append(resp, []byte(e.Error())...)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	w.Write([]byte("ok"))
}
