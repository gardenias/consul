package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/command"
	"github.com/hashicorp/consul/command/agent"
	"github.com/hashicorp/consul/command/base"
	"github.com/hashicorp/consul/version"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &agent.Command{
				Revision:          version.GitCommit,
				Version:           version.Version,
				VersionPrerelease: version.VersionPrerelease,
				HumanVersion:      version.GetHumanVersion(),
				Ui:                ui,
				ShutdownCh:        make(chan struct{}),
			}, nil
		},

		"configtest": func() (cli.Command, error) {
			return &command.ConfigTestCommand{
				Command: base.Command{
					Flags: base.FlagSetNone,
					Ui:    ui,
				},
			}, nil
		},

		"event": func() (cli.Command, error) {
			return &command.EventCommand{
				Command: base.Command{
					Flags: base.FlagSetHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"exec": func() (cli.Command, error) {
			return &command.ExecCommand{
				ShutdownCh: makeShutdownCh(),
				Command: base.Command{
					Flags: base.FlagSetHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"force-leave": func() (cli.Command, error) {
			return &command.ForceLeaveCommand{
				Command: base.Command{
					Flags: base.FlagSetClientHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"info": func() (cli.Command, error) {
			return &command.InfoCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetClientHTTP,
				},
			}, nil
		},

		"join": func() (cli.Command, error) {
			return &command.JoinCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetClientHTTP,
				},
			}, nil
		},

		"keygen": func() (cli.Command, error) {
			return &command.KeygenCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetNone,
				},
			}, nil
		},

		"keyring": func() (cli.Command, error) {
			return &command.KeyringCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetClientHTTP,
				},
			}, nil
		},

		"kv": func() (cli.Command, error) {
			return &command.KVCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetNone,
				},
			}, nil
		},

		"kv delete": func() (cli.Command, error) {
			return &command.KVDeleteCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetHTTP,
				},
			}, nil
		},

		"kv get": func() (cli.Command, error) {
			return &command.KVGetCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetHTTP,
				},
			}, nil
		},

		"kv put": func() (cli.Command, error) {
			return &command.KVPutCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetHTTP,
				},
			}, nil
		},

		"kv export": func() (cli.Command, error) {
			return &command.KVExportCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetHTTP,
				},
			}, nil
		},

		"kv import": func() (cli.Command, error) {
			return &command.KVImportCommand{
				Command: base.Command{
					Ui:    ui,
					Flags: base.FlagSetHTTP,
				},
			}, nil
		},

		"leave": func() (cli.Command, error) {
			return &command.LeaveCommand{
				Ui: ui,
			}, nil
		},

		"lock": func() (cli.Command, error) {
			return &command.LockCommand{
				ShutdownCh: makeShutdownCh(),
				Command: base.Command{
					Flags: base.FlagSetHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"maint": func() (cli.Command, error) {
			return &command.MaintCommand{
				Ui: ui,
			}, nil
		},

		"members": func() (cli.Command, error) {
			return &command.MembersCommand{
				Ui: ui,
			}, nil
		},

		"monitor": func() (cli.Command, error) {
			return &command.MonitorCommand{
				ShutdownCh: makeShutdownCh(),
				Ui:         ui,
			}, nil
		},

		"operator": func() (cli.Command, error) {
			return &command.OperatorCommand{
				Ui: ui,
			}, nil
		},

		"reload": func() (cli.Command, error) {
			return &command.ReloadCommand{
				Ui: ui,
			}, nil
		},

		"rtt": func() (cli.Command, error) {
			return &command.RTTCommand{
				Ui: ui,
			}, nil
		},

		"snapshot": func() (cli.Command, error) {
			return &command.SnapshotCommand{
				Ui: ui,
			}, nil
		},

		"snapshot restore": func() (cli.Command, error) {
			return &command.SnapshotRestoreCommand{
				Ui: ui,
			}, nil
		},

		"snapshot save": func() (cli.Command, error) {
			return &command.SnapshotSaveCommand{
				Ui: ui,
			}, nil
		},

		"snapshot inspect": func() (cli.Command, error) {
			return &command.SnapshotInspectCommand{
				Ui: ui,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				HumanVersion: version.GetHumanVersion(),
				Ui:           ui,
			}, nil
		},

		"watch": func() (cli.Command, error) {
			return &command.WatchCommand{
				ShutdownCh: makeShutdownCh(),
				Ui:         ui,
			}, nil
		},
	}
}

// makeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// interrupt or SIGTERM received.
func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}
