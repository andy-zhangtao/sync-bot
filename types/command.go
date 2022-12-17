package types

const (
	DockerName  = "/docker-name"
	DockerBuild = "/docker-build"
)

const (
	Ready = iota
	Start
	Finish
)

type DockerTask struct {
	Name  string
	Build string
	Stage int
}
