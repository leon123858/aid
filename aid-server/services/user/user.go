package user

import (
	"aid-server/pkg/ldb"
	"encoding/json"
	"github.com/google/uuid"
)

const bucketSize = 10

type IUser interface {
	IsExist() bool
	GetAID() uuid.UUID
	GetInfo() Info
	GetSpace() *Space
	GetTime() *Time
	SetRecord(Record) error
	SetInfo(Info) error
	//flag.Getter
}

func CreateUser(aid uuid.UUID, db ldb.DB) (IUser, error) {
	if db == nil || aid == uuid.Nil {
		panic("invalid parameter")
	}
	isExist, err := db.IsExist(aid.String())
	if err != nil {
		return nil, err
	}
	if !isExist {
		data, err := json.Marshal(Data{
			Records: make([]Record, 0),
			Info: Info{
				PublicKey: "",
				AID:       aid.String(),
			},
		})
		if err != nil {
			return nil, err
		}
		if err = db.Set(aid.String(), string(data)); err != nil {
			return nil, err
		}
		return &User{
			ID: aid,
			DB: db,
			Data: Data{
				Records: make([]Record, 0),
				Info: Info{
					PublicKey: "",
					AID:       aid.String(),
				},
			},
			IsExisted: false,
		}, nil
	}
	data, err := db.Get(aid.String())
	if err != nil {
		return nil, err
	}
	var d Data
	if err = json.Unmarshal([]byte(data), &d); err != nil {
		return nil, err
	}
	return &User{
		ID:        aid,
		DB:        db,
		Data:      d,
		IsExisted: true,
	}, nil
}

func (u *User) IsExist() bool {
	return u.IsExisted
}

func (u *User) GetAID() uuid.UUID {
	return u.ID
}

func (u *User) GetInfo() Info {
	return u.Info
}

func (u *User) GetSpace() *Space {
	data := u.Data
	records := data.Records
	if len(records) == 0 {
		return &Space{
			DeviceFingerPrint: DeviceFingerPrint{
				IP:   "",
				Brow: "",
			},
		}
	}
	return &records[len(records)-1].Space
}

func (u *User) GetTime() *Time {
	data := u.Data
	records := data.Records
	if len(records) == 0 {
		return &Time{
			CurEventTime: 0,
		}
	}
	return &records[len(records)-1].Time
}

func (u *User) updateData(data Data) error {
	bData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err = u.DB.Set(u.ID.String(), string(bData)); err != nil {
		return err
	}
	return nil
}

func (u *User) SetRecord(d Record) error {
	u.Data.Records = append(u.Data.Records, d)
	if len(u.Data.Records) > bucketSize {
		u.Data.Records = u.Data.Records[1:]
	}
	if err := u.updateData(u.Data); err != nil {
		return err
	}
	return nil
}

func (u *User) SetInfo(info Info) error {
	u.Data.Info = info
	if err := u.updateData(u.Data); err != nil {
		return err
	}
	return nil
}
