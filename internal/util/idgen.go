package util

import (
	"context"

	"github.com/khicago/got/frameworker/idgen"
)

var (
	remoteGen idgen.IGenerator = nil
	// todo: this is a temporary implementation
	localGen = idgen.NewLocalMUGen(1, true)
)

func getIDGen() idgen.IGenerator {
	if remoteGen != nil {
		return idgen.NewIDGen(remoteGen, localGen)
	}
	return localGen
}

func GenIDU64(ctx context.Context) (uint64, error) {
	id, err := getIDGen().Get(ctx)
	if err != nil {
		return 0, err
	}
	return uint64(id), err
}
