package share

import "sync-bot/types"

var dockerInfo map[int64][]types.DockerTask

func NewDockerTask() {
	if dockerInfo == nil {
		dockerInfo = make(map[int64][]types.DockerTask)
	}
}

func DockTask() map[int64][]types.DockerTask {
	return dockerInfo
}
