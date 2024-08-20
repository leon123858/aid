package controller

import (
	"aid-server/service/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type HashRequest struct {
	Hash string `json:"hash" example:"1234567890abcdef1234567890abcdef"`
	Aid  string `json:"aid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// GetCertHash godoc
//
//	@Summary		Get certificate hash
//	@Description	Retrieve the certificate hash for a given AID
//	@Tags			utils
//	@Accept			json
//	@Produce		json
//	@Param			aid	query		string	true	"AID (UUID)"
//	@Success		200	{object}	string	"Certificate hash"
//	@Failure		400	{object}	string	"Bad Request"
//	@Failure		404	{object}	string	"Not Found"
//	@Router			/verify/hash [get]
func GetCertHash(c echo.Context) error {
	cw := ContextWrapper{c}
	aid := c.QueryParam("aid")
	if aid == "" {
		return cw.newBadRequestError("aid is empty")
	}
	if _, err := uuid.Parse(aid); err != nil {
		return cw.newBadRequestError("aid is not valid")
	}
	hash, err := auth.LoadHash(uuid.MustParse(aid))
	if err != nil {
		return cw.newNotFound("hash not found")
	}
	return cw.newSuccess(hash)
}

// SaveCertHash godoc
//
//	@Summary		Save certificate hash
//	@Description	Save a new certificate hash for a given AID
//	@Tags			utils
//	@Accept			json
//	@Produce		json
//	@Param			request	body		HashRequest	true	"Hash request"
//	@Success		200		{object}	string		"Success message"
//	@Failure		400		{object}	string		"Bad Request"
//	@Failure		500		{object}	string		"Internal Server Error"
//	@Router			/register/cert [post]
func SaveCertHash(c echo.Context) error {
	cw := ContextWrapper{c}
	req := new(HashRequest)
	if err := c.Bind(req); err != nil {
		return cw.newBadRequestError(err.Error())
	}
	// check if the request is valid
	if req.Hash == "" || req.Aid == "" {
		return cw.newBadRequestError("hash or aid is empty")
	}
	// is aid is valid uuid
	if _, err := uuid.Parse(req.Aid); err != nil {
		return cw.newBadRequestError("aid is not valid")
	}
	// if aid exists
	if r, err := auth.LoadHash(uuid.MustParse(req.Aid)); err == nil {
		if r == req.Hash {
			return cw.newSuccess("success")
		}
		return cw.newBadRequestError("aid already exists")
	}
	// save hash
	if err := auth.SaveHash(uuid.MustParse(req.Aid), req.Hash); err != nil {
		return cw.newInternalServerError(err.Error())
	}
	return cw.newSuccess("success")
}
