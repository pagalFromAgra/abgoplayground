package uploadlogs

import "time"

const (
	ERR_WRONG_PAYLOAD = "ERR_WRONG_PAYLOAD"

	TIME_LAYOUT    = "2006-01-02T15:04:05.000"
	DL_TIME_LAYOUT = "2006-01-02T15:04:05"

	ENV_KEEN_PROJECT_ID = "KEEN_PROJECT_ID"
	ENV_KEEN_WRITE_KEY  = "KEEN_WRITE_KEY"
	ENV_KEEN_COLLECTION = "KEEN_COLLECTION"

	// KEEN_PROJECT_ID = "57c6315e8db53dfda8a6d88d"
	// KEEN_WRITE_KEY  = "2A31299722BA2077AC732461D2864F3BC93B12F3784A1859944E1BC160AA11FE22669F08E692F0BA6A599EA9D4F948FE5ABC11C72E7405A5E39F46699A8992BFD7E71DEE46E3EC6C204EFF4D75FA8A4A20718E8CA3368597BC175430CA9F4DAC"
	//
	// KEEN_COLLECTION_NAME = "logs"

	KEEN_TIME_FORMAT = "2006-01-02T15:04:05.000Z"
)

var T0 = time.Time{}
