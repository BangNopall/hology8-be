package uuid

import (
	"github.com/google/uuid"
	"github.com/hology8/hology-be/pkg/log"
)

type UUIDInterface interface {
	New() (uuid.UUID, error)
}

type UUIDStruct struct{}

var UUID = getUUID()

func getUUID() UUIDInterface {
	return &UUIDStruct{}
}

func (u *UUIDStruct) New() (uuid.UUID, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UUID][New] failed to create uuid v7")

		return uuid, err
	}

	return uuid, err
}

// I'm too lazy to change service test to DI ver, so here we go
// Use this one if you dont need the DI above version
func New() uuid.UUID {
	uuid, err := uuid.NewV7()

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UUID][New] failed to create uuid v7")

		return uuid
	}

	return uuid
}
