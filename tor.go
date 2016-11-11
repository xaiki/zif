package zif

import (
	"fmt"
	"os"
	"path/filepath"

	torc "github.com/postfix/goControlTor"
	log "github.com/sirupsen/logrus"
)

func SetupZifTorService(port, tor int, cookie string) (error, *torc.TorControl, string) {
	control := &torc.TorControl{}

	serviceDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Info(serviceDir)
	servicePort := map[int]string{port: fmt.Sprintf("127.0.0.1:%d", port)}

	err := control.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tor))

	if err != nil {
		log.Error(err.Error())
		return err, nil, ""
	}

	log.Info("Dialed Tor control port")

	err = control.CookieAuthenticate(cookie)

	if err != nil {
		log.Error(err.Error())
		return err, nil, ""
	}

	log.Info("Authenticated with Tor, creating service")

	err = control.CreateHiddenService(serviceDir, servicePort)

	if err != nil {
		log.Error(err.Error())
		return err, nil, ""
	}

	log.Info("Service created")

	onion, err := torc.ReadOnion(serviceDir)

	log.Info("Tor address ", onion)

	return nil, control, onion
}
