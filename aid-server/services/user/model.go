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
	ID        uuid.UUID
	DB        ldb.DB
	IsExisted bool
	Data
}

type Space struct {
	DeviceFingerPrint DeviceFingerPrint
	Info              Info
}

type Time struct {
	PreLoginTime time.Time
}

type DeviceFingerPrint struct {
	IP   string
	Brow string
}

type Info struct {
	PublicKey string
	AID       string
}
