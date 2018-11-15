package goscp

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func ScpHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var op ScpOptions
	err := decoder.Decode(&op)
	if err != nil {
		logrus.WithError(err).Warn("failed to parse JSON request")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(struct{ Status string }{"Bad JSON"})
		return
	}

	logrus.Info(op)
	Scp(op)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(struct{ Status string }{"ok"})
}
