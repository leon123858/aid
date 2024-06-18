package user

import (
	"aid-server/pkg/ldb"
	"encoding/json"
	"github.com/google/uuid"
)

type IUser interface {
	GetAID() uuid.UUID
	GetSpace() *Space
	GetTime() *Time
	SetSpace(Space) error
	SetTime(Time) error
	SetAll(Data) error
}

func NewDB(path string) (ldb.DB, error) {
	db := ldb.New()
	err := db.Connect(path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func FreeDB(db ldb.DB) error {
	return db.Close()
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
			Space: Space{},
			Time:  Time{},
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
				Space: Space{},
				Time:  Time{},
			},
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
		ID:   aid,
		DB:   db,
		Data: d,
	}, nil
}

func (u *User) GetAID() uuid.UUID {
	return u.ID
}

func (u *User) GetSpace() *Space {
	return &u.Space
}

func (u *User) GetTime() *Time {
	return &u.Time
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

func (u *User) SetSpace(s Space) error {
	if err := u.updateData(Data{
		Space: s,
		Time:  u.Time,
	}); err != nil {
		return err
	}
	u.Space = s
	return nil
}

func (u *User) SetTime(t Time) error {
	if err := u.updateData(Data{
		Space: u.Space,
		Time:  t,
	}); err != nil {
		return err
	}
	u.Time = t
	return nil
}

func (u *User) SetAll(d Data) error {
	if err := u.updateData(d); err != nil {
		return err
	}
	u.Data = d
	return nil
}
