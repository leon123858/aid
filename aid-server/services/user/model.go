package user

import (
	"aid-server/pkg/ldb"
	"aid-server/pkg/timestamp"
	"github.com/google/uuid"
)

type Data struct {
	Records []Record
	Info
}

type Record struct {
	Space
	Time
}

type User struct {
	ID        uuid.UUID
	DB        ldb.DB
	IsExisted bool
	Data
}

type Space struct {
	DeviceFingerPrint DeviceFingerPrint
}

type Time struct {
	CurEventTime timestamp.Timestamp
}

type DeviceFingerPrint struct {
	IP   string
	Brow string
}

type Info struct {
	PublicKey string
	AID       string
}
