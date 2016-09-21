package zif

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// TODO: Make this check using UpNp/NAT_PMP first, then query services.
func external_ip() string {
	resp, err := http.Get("https://api.ipify.org/")
	defer resp.Body.Close()

	if err != nil {
		log.Error("Failed to get external ip: try setting manually")
		return ""
	}

	ret, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error("Failed to get external ip: try setting manually")
		return ""
	}

	return string(ret)
}
