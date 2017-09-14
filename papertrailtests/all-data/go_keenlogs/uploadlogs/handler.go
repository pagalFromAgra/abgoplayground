package uploadlogs

import (
	"errors"
	"log"
	"os"

	"github.com/ironbay/dynamic"
)

// Handler responds to valid payload
func Handler(payload *Payload) error {

	keenevent := dynamic.Build(
		"device", payload.Device,
		"company", payload.Company,
		"event", payload.Event,
		"value", payload.Value,
		"vbat", payload.Vbat,
		"temperature", payload.Temprature,
		"info", payload.Info,
		"keen", dynamic.Build(
			"timestamp", payload.DateTime,
		),
	)

	keenClient, err := newKeenClient()
	if err != nil {
		log.Println("Couldn't connect to Keen")
	}

	keencollection := os.Getenv(ENV_KEEN_COLLECTION)
	if keencollection == "" {
		return errors.New("Need KEEN_COLLECTION")
	}

	err = keenClient.AddEvent(keencollection, keenevent)
	if err != nil {
		return err
	}

	return nil
}
