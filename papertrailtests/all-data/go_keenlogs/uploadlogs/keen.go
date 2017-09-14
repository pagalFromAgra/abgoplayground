package uploadlogs

import (
	"errors"
	"os"

	"github.com/inconshreveable/go-keen"
)

// Client instanciates from environment variables
func newKeenClient() (*keen.Client, error) {

	writekey := os.Getenv(ENV_KEEN_WRITE_KEY)
	if writekey == "" {
		return nil, errors.New("Need KEEN_WRITE_KEY")
	}

	projectid := os.Getenv(ENV_KEEN_PROJECT_ID)
	if projectid == "" {
		return nil, errors.New("Need KEEN_PROJECT_ID")
	}

	return &keen.Client{
		WriteKey:  writekey,
		ProjectID: projectid,
	}, nil

}
