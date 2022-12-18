package commands

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"strings"
	"sync-bot/services/github"
	"sync-bot/share"
	"sync-bot/utils"
)

type DockerBuild struct {
	Name  string
	Build string
	ctx   context.Context
}

type Travis struct {
	Language string   `yaml:"language"`
	Go       []string `yaml:"go"`
	Env      struct {
		Global []string `yaml:"global"`
	} `yaml:"env"`
	Services      []string `yaml:"services"`
	BeforeInstall []string `yaml:"before_install"`
	Install       []string `yaml:"install"`
	Script        []string `yaml:"script"`
	AfterSuccess  []string `yaml:"after_success"`
	AfterFailure  []string `yaml:"after_failure"`
	Notifications struct {
		Slack struct {
			Rooms    []string `yaml:"rooms"`
			Template []string `yaml:"template"`
		} `yaml:"slack"`
		Email struct {
			Recipients []string `yaml:"recipients"`
			OnSuccess  string   `yaml:"on_success"`
			OnFailure  string   `yaml:"on_failure"`
		} `yaml:"email"`
	} `yaml:"notifications"`
}

func (dn *DockerBuild) Kind() string {
	return "docker-build"
}

func (dn *DockerBuild) Content() string {
	return dn.Build
}

func (dn *DockerBuild) SetContext(ctx context.Context) {
	dn.ctx = ctx
}

func (dn *DockerBuild) Context() context.Context {
	return dn.ctx
}

func (dn *DockerBuild) SetReply(reply string) {
	dn.Name = utils.GrabCommand(reply)
}

func (dn *DockerBuild) SetSource(chatId int64) {
	for _, task := range share.DockTask()[chatId] {
		if task.Name != "" && task.Build == "" {
			task.Build = dn.Build
			dn.Name = task.Name
			break
		}
	}
}

func (dn *DockerBuild) Answer() (*string, bool) {
	result := fmt.Sprintf("[%s] has start", dn.Name)
	return &result, true
}

func (dn *DockerBuild) Run() (result string, err error) {

	log.Infof("name: %s  build: %s", dn.Name, dn.Build)
	// get last commit sha
	commit, err := share.GithubHelper().GetLastCommit()
	if err != nil {
		return result, err
	}

	// get repo tree with sha
	tree, err := share.GithubHelper().GetTreeWithSha(commit.Object.Sha)
	if err != nil {
		return result, err
	}

	// modify file content
	travis, err := dn.travisConfig(dn.extractImageAndTag(dn.Name))
	if err != nil {
		return result, err
	}

	trees := []github.RepoTree{
		{
			Path:    "Dockerfile",
			Mode:    "100644",
			Type:    "blob",
			Content: dn.Build,
		},
		{
			Path:    ".travis.yml",
			Mode:    "100644",
			Type:    "blob",
			Content: travis,
		},
	}
	resp, err := share.GithubHelper().NewTreeWithContent(tree.Sha, trees)
	if err != nil {
		return result, err
	}

	// commit prev modify
	commitSha, err := share.GithubHelper().CreateNewCommit(commit.Object.Sha, resp.Sha, dn.ctx.Value("auth").(string), dn.ctx.Value("email").(string))
	if err != nil {
		return result, err
	}

	// update branch
	err = share.GithubHelper().UpdateBranch(commitSha)
	if err != nil {
		return result, err
	}

	return fmt.Sprintf("[%s] Has in build queue.", dn.Name), nil
}

func (dn *DockerBuild) Inspect(update tgbotapi.Update) {

}

func (dn *DockerBuild) extractImageAndTag(name string) (image, tag string) {
	_name := strings.Split(name, ":")
	if len(_name) == 1 {
		return name, "latest"
	}

	return _name[0], _name[1]
}

func (dn *DockerBuild) travisConfig(image, tag string) (result string, err error) {
	t := Travis{
		Language: "go",
		Go:       []string{"1.15.5"},
		Env: struct {
			Global []string `yaml:"global"`
		}{
			Global: []string{
				"AREA=hongkong",
				fmt.Sprintf("IMAGE=%s", image),
				fmt.Sprintf("TAG=%s", tag),
			},
		},
		Services: []string{
			"docker",
		},
		BeforeInstall: []string{
			"echo \"$DOCKERPASSWD\" | docker login -u \"$DOCKERUSER\" --password-stdin registry.cn-$AREA.aliyuncs.com",
			"echo \"$OFFDOCKERPASSWD\" | docker login -u \"$OFFDOCKERUSER\" --password-stdin",
		},
		Install: []string{
			"docker build -t registry.cn-$AREA.aliyuncs.com/$IMAGE:$TAG .",
			"docker tag registry.cn-$AREA.aliyuncs.com/$IMAGE:$TAG $IMAGE:$TAG",
		},
		Script: []string{
			"echo \"ignore go test command \"",
		},
		AfterSuccess: []string{
			"docker push registry.cn-$AREA.aliyuncs.com/$IMAGE:$TAG",
			"docker push $IMAGE:$TAG",
			"curl \"https://api.telegram.org/bot$TGBOOT/sendMessage?chat_id=-1001203454731&parse_mode=Markdown&text=*registry.cn-$AREA.aliyuncs.com/$IMAGE:$TAG SUCCESS*\"",
		},
		AfterFailure: []string{
			"curl \"https://api.telegram.org/bot$TGBOOT/sendMessage?chat_id=-1001203454731&parse_mode=Markdown&text=*registry.cn-$AREA.aliyuncs.com/$IMAGE:$TAG FAILED*\"",
		},
		Notifications: struct {
			Slack struct {
				Rooms    []string `yaml:"rooms"`
				Template []string `yaml:"template"`
			} `yaml:"slack"`
			Email struct {
				Recipients []string `yaml:"recipients"`
				OnSuccess  string   `yaml:"on_success"`
				OnFailure  string   `yaml:"on_failure"`
			} `yaml:"email"`
		}{
			Slack: struct {
				Rooms    []string `yaml:"rooms"`
				Template []string `yaml:"template"`
			}(struct {
				Rooms    []string
				Template []string
			}{Rooms: []string{
				"sync-from-docker-hub:5Q3mjXPHAuJ50LeUcUrfIWmx",
			}, Template: []string{
				"\"Build <%{build_url}|#%{build_number}> (<%{compare_url}|%{commit}>) of %{repository_slug}@%{branch} by %{author} %{result} in %{duration}\"",
				"\"%{commit_message}\"",
				"\"Result: %{result}\"",
			}}),
			Email: struct {
				Recipients []string `yaml:"recipients"`
				OnSuccess  string   `yaml:"on_success"`
				OnFailure  string   `yaml:"on_failure"`
			}(struct {
				Recipients []string
				OnSuccess  string
				OnFailure  string
			}{Recipients: []string{
				"ztao8607@gmail.com",
			}, OnSuccess: "always", OnFailure: "always"}),
		},
	}

	data, err := yaml.Marshal(t)
	if err != nil {
		return result, err
	}

	return string(data), nil
}
