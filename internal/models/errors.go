package model

import "errors"

var (
	ErrResourceNotFound    = errors.New("resource not found")
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrSessionNotPresent   = errors.New("session not present")
	ErrSessionNotActive    = errors.New("session is not active")
	ErrInvalidAssetType    = errors.New("invalid asset type")
	ErrInvalidAssetSubType = errors.New("invalid asset sub type")
	ErrInvalidAssetID      = errors.New("invalid asset ID")
	ErrResourceGone        = errors.New("resource gone")
)
