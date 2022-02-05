package lib

import (
	"encoding/json"
	"net/http"
)

func marshalAndWrite(w http.ResponseWriter, payload interface{}, status int) error {
	if payload != nil {
		marshaled, errMar := json.Marshal(payload)
		if errMar != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := []byte("{\"message\":\"Could not marshal error\"}")
			_, err := w.Write(response)
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.WriteHeader(status)

		_, err := w.Write(marshaled)
		return err
	}

	w.WriteHeader(status)
	return nil
}
