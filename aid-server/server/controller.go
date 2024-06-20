package server

import (
	"aid-server/configs"
	"aid-server/pkg/jwt"
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
// @Success		200	{object}	Response		"JWT Token"
// @Router			/api/login [post]
func login(c echo.Context) error {
	req := LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	aidUUID, err := uuid.Parse(req.AID)
	if err != nil {
		return c.JSON(400, generateResponse(false, "invalid AID"))
	}
	userItem, err := user.CreateUser(aidUUID, UserDB)
	if err != nil {
		return c.JSON(500, generateResponse(false, "user item creation failed"))
	}
	if !userItem.IsExist() {
		return c.JSON(400, generateResponse(false, "user not existed"))
	}
	if result, err := rsa.VerifySignature([]byte(userItem.GetSpace().Info.PublicKey), []byte(req.Timestamp), req.Sign); err != nil || !result {
		return c.JSON(400, generateResponse(false, "invalid signature"))
	}
	if !timestamp.CheckTimestampClose5000(timestamp.ToTimestamp(req.Timestamp), timestamp.GetTime()) {
		return c.JSON(400, generateResponse(false, "expired timestamp"))
	}
	err = userItem.SetAll(user.Data{
		Time: user.Time{
			PreLoginTime: timestamp.GetTime(),
		},
		Space: user.Space{
			DeviceFingerPrint: user.DeviceFingerPrint{
				IP:   req.IP,
				Brow: req.Browser,
			},
		},
	})
	if err != nil {
		return c.JSON(500, generateResponse(false, "user item update failed"))
	}
	token, err := jwt.GenerateToken(aidUUID.String())
	if err != nil {
		return c.JSON(500, generateResponse(false, "token generation failed"))
	}
	return c.JSON(200, generateResponse(true, token))
}

// @Summary		Logout
// @Description	Logout
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		LogoutRequest	true	"Logout Request"
// @Success		200	{object}	Response		"empty string"
// @Router			/api/logout [post]
func logout(c echo.Context) error {
	req := LogoutRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, generateResponse(true, ""))
}

// @Summary		Register
// @Description	Register
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		RegisterRequest	true	"Register Request"
// @Success		200	{object}	Response		"JWT Token"
// @Router			/api/register [post]
func register(c echo.Context) error {
	req := RegisterRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, generateResponse(false, "invalid request body"))
	}
	if !rsa.IsValidPublicKey([]byte(req.PublicKey)) {
		return c.JSON(400, generateResponse(false, "invalid public key"))
	}
	aidUUID, err := uuid.Parse(req.AID)
	if err != nil {
		return c.JSON(400, generateResponse(false, "invalid AID"))
	}
	userItem, err := user.CreateUser(aidUUID, UserDB)
	if err != nil {
		return c.JSON(500, generateResponse(false, "user item creation failed"))
	}
	if userItem.IsExist() {
		return c.JSON(400, generateResponse(false, "user already existed"))
	}
	err = userItem.SetAll(user.Data{
		Space: user.Space{
			DeviceFingerPrint: user.DeviceFingerPrint{
				IP:   req.IP,
				Brow: req.Browser,
			},
			Info: user.Info{
				PublicKey: req.PublicKey,
				AID:       req.AID,
			},
		},
		Time: user.Time{
			PreLoginTime: timestamp.GetTime(),
		},
	})
	if err != nil {
		return c.JSON(500, generateResponse(false, "user item update failed"))
	}
	token, err := jwt.GenerateToken(aidUUID.String())
	if err != nil {
		return c.JSON(500, generateResponse(false, "token generation failed"))
	}
	return c.JSON(200, generateResponse(true, token))
}

// @Summary		Ask
// @Description	Ask
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		AskRequest	true	"Ask Request"
// @Success		200	{object}	Response	"aid string"
// @Router			/api/ask [post]
func ask(c echo.Context) error {
	req := AskRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, generateResponse(true, ""))
}

// @Summary		Trigger
// @Description	Trigger
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		TriggerRequest	true	"Trigger Request"
// @Success		200	{object}	Response		"aid string"
// @Router			/api/trigger [post]
func trigger(c echo.Context) error {
	req := TriggerRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, generateResponse(true, ""))
}
