package server

import (
	"aid-server/configs"
	"aid-server/pkg/jwt"
	"aid-server/pkg/ldb"
	"aid-server/pkg/res"
	"aid-server/pkg/rsa"
	"aid-server/pkg/timestamp"
	"aid-server/services/alias"
	"aid-server/services/idmap"
	"aid-server/services/localAPIWrapper"
	"aid-server/services/mlm"
	"aid-server/services/rba"
	"aid-server/services/user"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var UIDMapAID *idmap.IDMap
var RecordCache mlm.MultiLevelMap
var AliasPool *alias.DB

func init() {
	var err error
	if UserDB, err = ldb.NewDB(configs.Configs.Path.UserDB); err != nil {
		panic(err)
	}
	if UserMapDB, err = ldb.NewDB(configs.Configs.Path.IDMap); err != nil {
		panic(err)
	}
	if AliasPool, err = alias.NewDB(configs.Configs.Path.AliasDB); err != nil {
		panic(err)
	}
	UIDMapAID = idmap.NewIDMap(100, UserMapDB)
	RecordCache = mlm.NewMultiLevelMap()
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
	if req.IP == "" {
		req.IP = c.RealIP()
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
	err = RecordCache.Set(mlm.KeyItem{
		IP:      req.IP,
		Browser: req.Browser,
	}, aidUUID)
	if err != nil {
		return err
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
	if req.IP == "" {
		req.IP = c.RealIP()
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
// @Description	service ask aid server to get new unique id(uuid) for user in service
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		AskRequest		true	"Ask Request"
// @Success		200	{object}	res.Response	"uuid that can map aid, so that aid can map many uuids"
// @Router			/api/ask [post]
func ask(c echo.Context) error {
	req := AskRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, res.GenerateResponse(false, "invalid request body"))
	}
	if req.IP == "" {
		req.IP = c.RealIP()
	}
	get, err := RecordCache.Get(mlm.KeyItem{
		IP:      req.IP,
		Browser: req.Browser,
	})
	if err != nil {
		return c.JSON(400, res.GenerateResponse(false, err.Error()))
	}
	if len(get) != 1 {
		fmt.Printf("get: %v\n", get)
		return c.JSON(400, res.GenerateResponse(false, "can not get unique aid"))
	}
	aidUUID := get[0]
	userItem, err := user.CreateUser(aidUUID, UserDB)
	if err != nil {
		return c.JSON(500, res.GenerateResponse(false, "user item creation failed"))
	}
	if !userItem.IsExist() {
		return c.JSON(400, res.GenerateResponse(false, "user not existed"))
	}
	uid := uuid.New().String()
	if err := UIDMapAID.Set(uid, aidUUID.String()); err != nil {
		return c.JSON(500, res.GenerateResponse(false, "id map set failed"))
	}
	return c.JSON(200, res.GenerateResponse(true, uid))
}

// @Summary		Check
// @Description	given uid, check the service to check user status, maybe ask user should log in again in aid server
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		CheckRequest	true	"Check Request"
// @Success		200	{object}	res.Response	"safe status"
// @Router			/api/check [post]
func check(c echo.Context) error {
	req := CheckRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.IP == "" {
		req.IP = c.RealIP()
	}
	aid, err := UIDMapAID.Get(req.UID)
	if err != nil || aid == "" {
		return c.JSON(400, res.GenerateResponse(false, "invalid uid"))
	}
	aidUUID, err := uuid.Parse(aid)
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
	// check user status
	if !rba.SimpleAlgo.Verify(userItem, &user.Record{
		Space: user.Space{
			DeviceFingerPrint: user.DeviceFingerPrint{
				IP:   req.IP,
				Brow: req.Browser,
			},
		},
		Time: user.Time{
			CurEventTime: timestamp.GetTime(),
		},
	}) {
		return c.JSON(200, res.GenerateResponse(true, string(Offline)))
	}
	return c.JSON(200, res.GenerateResponse(true, string(Online)))
}

// @Summary		Verify
// @Description	given uid and JWT token, verify is the JWT token is valid
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Param			req	body		VerifyRequest	true	"Verify Request"
// @Success		200	{object}	res.Response	"result"
// @Router			/api/verify [post]
func verify(c echo.Context) error {
	req := VerifyRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	claims, ok := c.Get("claims").(*jwt.UserClaims)
	if !ok {
		return c.JSON(400, res.GenerateResponse(false, "invalid claims"))
	}
	aid, err := UIDMapAID.Get(req.UID)
	if err != nil || aid == "" {
		return c.JSON(400, res.GenerateResponse(false, "invalid uid"))
	}
	if aid != claims.ID {
		return c.JSON(400, res.GenerateResponse(false, "token not match uid"))
	}
	return c.JSON(200, res.GenerateResponse(true, "token match uid"))
}

// @Summary		Alias Login
// @Description	Alias Login for usage
// @Tags			Usage
// @Accept			json
// @Produce		json
// @Param			req	body		localAPIWrapper.UsageRequest	true	"Usage Request"
// @Success		200	{object}	localAPIWrapper.UsageResponse
// @Router			/usage/login [post]
func loginAlias(c echo.Context) error {
	req := localAPIWrapper.UsageRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: "invalid request body",
		})
	}
	onlineList := make([]string, 0)
	wrapper := localAPIWrapper.New()
	aliasList, err := AliasPool.ValidateUser(req.Username, req.Password)
	if err != nil {
		return c.JSON(500, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: err.Error(),
		})
	}
	if len(aliasList) == 0 {
		return c.JSON(400, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: "invalid username or password",
		})
	}
	if len(aliasList) == 1 {
		// get pre login
		info, err := AliasPool.GetUserLoginHistory(aliasList[0])
		if err != nil {
			return c.JSON(500, localAPIWrapper.UsageResponse{
				Result:  false,
				Message: err.Error(),
			})
		}
		if info == nil {
			goto verify
		}
		// check if pre login record is same as current login
		if info.IP != req.IP || info.Browser != req.Fingerprint {
			goto verify
		}
		goto success
	}
verify:
	if req.Token != "" {
		for _, v := range aliasList {
			result, err := wrapper.Verify(req.Token, v)
			if err != nil {
				fmt.Printf("Verify: %s\n", err.Error())
				continue
			}
			if result.Result {
				onlineList = append(onlineList, v)
			}
		}
	} else {
		for _, v := range aliasList {
			result, err := wrapper.Check(v, req.IP, req.Fingerprint)
			if err != nil {
				fmt.Printf("Check: %s\n", err.Error())
				continue
			}
			if result.Result && result.Content == "online" {
				onlineList = append(onlineList, v)
			}
		}
	}
	if len(onlineList) == 0 {
		return c.JSON(400, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: "no online alias",
		})
	}
	if len(onlineList) > 1 {
		return c.JSON(400, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: "multiple online alias",
		})
	}
	aliasList = onlineList
success:
	//println(aliasList[0], req.IP, req.Fingerprint)
	if err := AliasPool.AddLoginRecord(aliasList[0], req.IP, req.Fingerprint); err != nil {
		return c.JSON(500, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: err.Error(),
		})
	}
	return c.JSON(200, localAPIWrapper.UsageResponse{
		Token:   "",
		UUID:    aliasList[0],
		Message: "Successfully login",
		Result:  true,
	})
}

// @Summary		Alias Register
// @Description	Alias Register for usage
// @Tags			Usage
// @Accept			json
// @Produce		json
// @Param			req	body		localAPIWrapper.UsageRequest	true	"Usage Request"
// @Success		200	{object}	localAPIWrapper.UsageResponse
// @Router			/usage/register [post]
func registerAlias(c echo.Context) error {
	req := localAPIWrapper.UsageRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: "invalid request body",
		})
	}
	wrapper := localAPIWrapper.New()
	result, err := wrapper.Ask(req.IP, req.Fingerprint)
	if err != nil {
		return c.JSON(500, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: err.Error(),
		})
	}
	uid := result.Content
	if err := AliasPool.AddUser(uid, req.Username, req.Password); err != nil {
		return c.JSON(500, localAPIWrapper.UsageResponse{
			Result:  false,
			Message: err.Error(),
		})
	}
	return c.JSON(200, localAPIWrapper.UsageResponse{
		Token:   "",
		UUID:    uid,
		Message: "Successfully registered",
		Result:  true,
	})
}
