package sanstorage

import (
	"context"
	// "fmt"

	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

type sanstorageManager struct {
	Credentials *sanmodel.StorageDeviceSettings
	ctx         context.Context
}

func newStorageMgr(creds *sanmodel.StorageDeviceSettings, ctx context.Context) (*sanstorageManager, error) {
	ssm := &sanstorageManager{
		Credentials: creds,
		ctx:         ctx,
	}
	return ssm, nil
}

func New(creds *sanmodel.StorageDeviceSettings, ctx context.Context) (Manager, error) {
	return newStorageMgr(creds, ctx)
}
