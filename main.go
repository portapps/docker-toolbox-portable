//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/portapps/portapps/v2"
	"github.com/portapps/portapps/v2/pkg/log"
	"github.com/portapps/portapps/v2/pkg/proc"
	"github.com/portapps/portapps/v2/pkg/utl"
)

type config struct {
	Machine machine `yaml:"machine" mapstructure:"machine"`
}

type machine struct {
	Name         string `yaml:"name" mapstructure:"name"`
	HostCIDR     string `yaml:"host_cidr" mapstructure:"host_cidr"`
	CPU          int    `yaml:"cpu" mapstructure:"cpu"`
	Ram          int    `yaml:"ram" mapstructure:"ram"`
	Disk         int    `yaml:"disk" mapstructure:"disk"`
	SharedName   string `yaml:"shared_name" mapstructure:"shared_name"`
	OnExitStop   bool   `yaml:"on_exit_stop" mapstructure:"on_exit_stop"`
	OnExitRemove bool   `yaml:"on_exit_remove" mapstructure:"on_exit_remove"`
}

var (
	app *portapps.App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Machine: machine{
			Name:         "default",
			HostCIDR:     "192.168.99.1/24",
			CPU:          1,
			Ram:          1024,
			Disk:         20000,
			SharedName:   "shared",
			OnExitStop:   false,
			OnExitRemove: false,
		},
	}

	// Init app
	if app, err = portapps.NewWithCfg("docker-toolbox-portable", "Docker Toolbox", cfg); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "git", "bin", "bash.exe")
	app.Args = []string{
		"--login",
		"-i",
		utl.PathJoin(app.AppPath, "start.sh"),
	}

	sharedPath := utl.CreateFolder(app.DataPath, "shared")
	storagePath := utl.CreateFolder(app.DataPath, "storage")

	postInstallGit := utl.PathJoin(app.AppPath, "git", "post-install.bat")
	if _, err := os.Stat(postInstallGit); err == nil {
		log.Info().Msg("Initializing git...")
		if err = proc.QuickCmd("cmd", []string{"/k", postInstallGit}); err != nil {
			log.Fatal().Err(err).Msg("Cannot initialize git")
		}
	}

	log.Info().Msg("Setting machine environment...")
	utl.OverrideEnv("MACHINE_NAME", cfg.Machine.Name)
	utl.OverrideEnv("MACHINE_HOST_CIDR", cfg.Machine.HostCIDR)
	utl.OverrideEnv("MACHINE_CPU", strconv.Itoa(cfg.Machine.CPU))
	utl.OverrideEnv("MACHINE_RAM", strconv.Itoa(cfg.Machine.Ram))
	utl.OverrideEnv("MACHINE_DISK", strconv.Itoa(cfg.Machine.Disk))
	utl.OverrideEnv("MACHINE_STORAGE_PATH", utl.FormatUnixPath(storagePath))
	utl.OverrideEnv("MACHINE_SHARED_NAME", cfg.Machine.SharedName)
	utl.OverrideEnv("MACHINE_SHARED_PATH", sharedPath)

	log.Info().Msg("Adding app to PATH...")
	utl.OverrideEnv("PATH", fmt.Sprintf("%s;%s", app.AppPath, os.Getenv("PATH")))

	log.Info().Msg("Starting up the shell... ")
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   app.AppPath,
		Sys: &syscall.SysProcAttr{
			CmdLine: fmt.Sprintf(` --login -i "%s"`, utl.PathJoin(app.AppPath, "start.sh")),
		},
	}

	defer func() {
		var exitArgs []string
		log.Info().Msg("Exiting...")

		if cfg.Machine.OnExitRemove {
			exitArgs = []string{"rm", "-f", cfg.Machine.Name}
		} else if cfg.Machine.OnExitStop {
			exitArgs = []string{"stop", cfg.Machine.Name}
		}

		if len(exitArgs) > 0 {
			if err := proc.QuickCmd("docker-machine", exitArgs); err != nil {
				log.Error().Err(err).Msg("docker-machine command error")
			}
		}
	}()

	process, err := os.StartProcess(app.Process, []string{}, &pa)
	if err != nil {
		log.Fatal().Err(err).Msg("Process failed")
	}
	if _, err := process.Wait(); err != nil {
		log.Error().Err(err).Msg("Process failed")
	}
}
