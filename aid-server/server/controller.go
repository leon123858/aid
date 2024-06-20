package server

import (
	"aid-server/configs"
	"aid-server/pkg/jwt"
	"aid-server/pkg/res"
	"aid-server/pkg/rsa"
	"aid-server/pkg/timestamp"
	"aid-server/services/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func init() {
	var err error
	if UserDB, err = user.NewDB(configs.Configs.Path.UserDB); err != nil {
		panic(err)
	}
}

// @Summary		Login
// @Description	Login
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		LoginRequest	true	"Login Request"
// @Success		200	{object}	res.Response	"JWT Token"
// @Router			/api/login [post]
func login(c echo.Context) error {
	req := LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	aidUUID, err := uuid.Parse(req.AID)
	if err != nil {
		return c.JSON(400, res.GenerateResponse(false, "invalid AID"))
	}
	userItem, err := user.CreateUser(aidUUID, UserDB)
	if err != nil {
		return c.JSON(500, res.GenerateResponse(false, "user item creation failed"))
	}
	if !userItem.IsExist() {
		return c.JSON(400, res.GenerateResponse(false, "user not existed"))
	}
	if result, err := rsa.VerifySignature([]byte(userItem.GetInfo().PublicKey), []byte(req.Timestamp), req.Sign); err != nil || !result {
		return c.JSON(400, res.GenerateResponse(false, "invalid signature"))
	}
	if !timestamp.CheckTimestampClose5000(timestamp.ToTimestamp(req.Timestamp), timestamp.GetTime()) {
		return c.JSON(400, res.GenerateResponse(false, "expired timestamp"))
	}
	curTime := timestamp.GetTime()
	err = userItem.SetRecord(user.Record{
		Time: user.Time{
			PreLoginTime: curTime,
			CurEventTime: curTime,
		},
		Space: user.Space{
			DeviceFingerPrint: user.DeviceFingerPrint{
				IP:   req.IP,
				Brow: req.Browser,
			},
		},
	})
	if err != nil {
		return c.JSON(500, res.GenerateResponse(false, "user item update failed"))
	}
	token, err := jwt.GenerateToken(aidUUID.String())
	if err != nil {
		return c.JSON(500, res.GenerateResponse(false, "token generation failed"))
	}
	return c.JSON(200, res.GenerateResponse(true, token))
}

// @Summary		Register
// @Description	Register
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		RegisterRequest	true	"Register Request"
// @Success		200	{object}	res.Response	"JWT Token"
// @Router			/api/register [post]
func register(c echo.Context) error {
	req := RegisterRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, res.GenerateResponse(false, "invalid request body"))
	}
	if !rsa.IsValidPublicKey([]byte(req.PublicKey)) {
		return c.JSON(400, res.GenerateResponse(false, "invalid public key"))
	}
	aidUUID, err := uuid.Parse(req.AID)
	if err != nil {
		return c.JSON(400, res.GenerateResponse(false, "invalid AID"))
	}
	userItem, err := user.CreateUser(aidUUID, UserDB)
	if err != nil {
		return c.JSON(500, res.GenerateResponse(false, "user item creation failed"))
	}
	if userItem.IsExist() {
		return c.JSON(400, res.GenerateResponse(false, "user already existed"))
	}
	curTime := timestamp.GetTime()
	err = userItem.SetRecord(user.Record{
		Space: user.Space{
			DeviceFingerPrint: user.DeviceFingerPrint{
				IP:   req.IP,
				Brow: req.Browser,
			},
		},
		Time: user.Time{
			PreLoginTime: curTime,
			CurEventTime: curTime,
		},
	})
	if err != nil {
		return c.JSON(500, res.GenerateResponse(false, "user record update failed"))
	}
	err = userItem.SetInfo(user.Info{
		PublicKey: req.PublicKey,
		AID:       aidUUID.String(),
	})
	if err != nil {
		return c.JSON(500, res.GenerateResponse(false, "user info update failed"))
	}
	token, err := jwt.GenerateToken(aidUUID.String())
	if err != nil {
		return c.JSON(500, res.GenerateResponse(false, "token generation failed"))
	}
	return c.JSON(200, res.GenerateResponse(true, token))
}

// @Summary		Ask
// @Description	Ask
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		AskRequest		true	"Ask Request"
// @Success		200	{object}	res.Response	"aid string"
// @Router			/api/ask [post]
func ask(c echo.Context) error {
	req := AskRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, res.GenerateResponse(true, ""))
}

// @Summary		Trigger
// @Description	Trigger
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		TriggerRequest	true	"Trigger Request"
// @Success		200	{object}	res.Response	"aid string"
// @Router			/api/trigger [post]
func trigger(c echo.Context) error {
	req := TriggerRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, res.GenerateResponse(true, ""))
}
