package sanstorage

import (
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

// Manager holds the state information
type Manager interface {
	GetHostGroupI(portID string, hostGroupNumber int) (*sanmodel.HostGroup, error)
}