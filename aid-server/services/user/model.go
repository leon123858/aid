package user

import (
	"aid-server/pkg/ldb"
	"github.com/google/uuid"
	"time"
)

type Data struct {
	Space
	Time
}

type User struct {
	ID uuid.UUID
	DB ldb.DB
	Data
}

type Space struct {
	DeviceFingerPrint DeviceFingerPrint
}

type Time struct {
	preLoginTime time.Time
}

type DeviceFingerPrint struct {
	IP   string
	OS   string
	Brow string
}
