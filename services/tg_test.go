package services

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"reflect"
	"sync-bot/services/commands"
	"testing"
)

func TestTG_parserMessage(t1 *testing.T) {
	type fields struct {
		bot      *tgbot.BotAPI
		interval int
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantCmd commands.Command
		wantErr bool
	}{
		{
			name: "normal command test",
			fields: fields{
				bot:      nil,
				interval: 0,
			},
			args:    args{msg: "/docker-name vikings"},
			wantCmd: commands.CommonCmd{Cmd: "/docker-name", Value: "vikings"},
			wantErr: false,
		},
		{
			name: "more multiple spaces",
			fields: fields{
				bot:      nil,
				interval: 0,
			},
			args:    args{msg: " /docker-name   vikings "},
			wantCmd: commands.CommonCmd{Cmd: "/docker-name", Value: "vikings"},
			wantErr: false,
		},
		{
			name: "add image tag",
			fields: fields{
				bot:      nil,
				interval: 0,
			},
			args:    args{msg: " /docker-name   vikings/nginx:v1 "},
			wantCmd: commands.CommonCmd{Cmd: "/docker-name", Value: "vikings/nginx:v1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TG{
				bot:      tt.fields.bot,
				interval: tt.fields.interval,
			}
			gotCmd, err := t.parserCommand(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t1.Errorf("parserMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCmd, tt.wantCmd) {
				t1.Errorf("parserMessage() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}

func TestTG_parserCommand(t1 *testing.T) {
	type fields struct {
		bot      *tgbot.BotAPI
		interval int
		task     map[string]commands.Command
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantCmd commands.Command
		wantErr bool
	}{
		{
			name: "parser docker image",
			fields: fields{
				bot:      nil,
				interval: 0,
				task:     nil,
			},
			args: args{msg: " /docker-name  vikings/ingress"},
			wantCmd: commands.CommonCmd{
				Cmd:   "/docker-name",
				Value: "vikings/ingress",
			},
			wantErr: false,
		},
		{
			name: "parser docker build",
			fields: fields{
				bot:      nil,
				interval: 0,
				task:     nil,
			},
			args: args{msg: `/docker-build   FROM xxxx 
RUN xxxx
ENTRYPOINT ["ingress"] `},
			wantCmd: commands.CommonCmd{
				Cmd: "/docker-build",
				Value: `FROM xxxx 
RUN xxxx
ENTRYPOINT ["ingress"]`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TG{
				bot:      tt.fields.bot,
				interval: tt.fields.interval,
				task:     tt.fields.task,
			}
			gotCmd, err := t.parserCommand(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t1.Errorf("parserCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCmd, tt.wantCmd) {
				t1.Errorf("parserCommand() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}
