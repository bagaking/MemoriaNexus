package utils

import (
	"context"
	"fmt"

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

func GenIDU64(ctx context.Context) (UInt64, error) {
	id, err := getIDGen().Get(ctx)
	if err != nil {
		return 0, err
	}
	return UInt64(id), err
}

func MGenIDU64(ctx context.Context, count int) ([]UInt64, error) {
	ids, err := getIDGen().MGet(ctx, int64(count))
	if err != nil {
		return nil, err
	}
	ret := make([]UInt64, 0, len(ids))
	for i := range ids {
		if ids[i] <= 0 {
			return nil, fmt.Errorf("idgen: invalid id: %d", ids[i])
		}
		ret = append(ret, UInt64(ids[i]))
	}

	return ret, err
}
