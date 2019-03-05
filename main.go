//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/proc"
	"github.com/portapps/portapps/pkg/win"
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

var cfg = config{
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

func init() {
	Papp.ID = "docker-toolbox-portable"
	Papp.Name = "Docker Toolbox"
	InitWithCfg(&cfg)

	win.SetConsoleTitle(fmt.Sprintf("%s Portable", Papp.Name))
}

func main() {
	Papp.AppPath = AppPathJoin("app")
	Papp.DataPath = CreateFolder(AppPathJoin("data"))
	Papp.Process = PathJoin(Papp.AppPath, "git", "bin", "bash.exe")
	Papp.Args = []string{
		"--login",
		"-i",
		PathJoin(Papp.AppPath, "start.sh"),
	}
	Papp.WorkingDir = Papp.AppPath

	sharedPath := CreateFolder(PathJoin(Papp.DataPath, "shared"))
	storagePath := CreateFolder(PathJoin(Papp.DataPath, "storage"))

	postInstallGit := PathJoin(Papp.AppPath, "git", "post-install.bat")
	if _, err := os.Stat(postInstallGit); err == nil {
		Log.Info("Initializing git...")
		if err = proc.QuickCmd("cmd", []string{"/k", postInstallGit}); err != nil {
			Log.Errorf("Cannot initialize git: %v", err)
		}
	}

	Log.Info("Setting machine environment...")
	OverrideEnv("MACHINE_NAME", cfg.Machine.Name)
	OverrideEnv("MACHINE_HOST_CIDR", cfg.Machine.HostCIDR)
	OverrideEnv("MACHINE_CPU", strconv.Itoa(cfg.Machine.CPU))
	OverrideEnv("MACHINE_RAM", strconv.Itoa(cfg.Machine.Ram))
	OverrideEnv("MACHINE_DISK", strconv.Itoa(cfg.Machine.Disk))
	OverrideEnv("MACHINE_STORAGE_PATH", FormatUnixPath(storagePath))
	OverrideEnv("MACHINE_SHARED_NAME", cfg.Machine.SharedName)
	OverrideEnv("MACHINE_SHARED_PATH", sharedPath)

	Log.Info("Adding app to PATH...")
	OverrideEnv("PATH", fmt.Sprintf("%s;%s", Papp.AppPath, os.Getenv("PATH")))

	Log.Info("Starting up the shell... ")
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   Papp.AppPath,
		Sys: &syscall.SysProcAttr{
			CmdLine: fmt.Sprintf(` --login -i "%s"`, PathJoin(Papp.AppPath, "start.sh")),
		},
	}
	process, err := os.StartProcess(Papp.Process, []string{}, &pa)
	if err != nil {
		Log.Fatal(err)
	}
	if _, err = process.Wait(); err != nil {
		Log.Fatal(err)
	}

	var exitArgs []string
	if cfg.Machine.OnExitRemove {
		exitArgs = []string{"rm", "-f", cfg.Machine.Name}
	} else if cfg.Machine.OnExitStop {
		exitArgs = []string{"stop", cfg.Machine.Name}
	}
	if len(exitArgs) > 0 {
		if err = proc.QuickCmd("docker-machine", exitArgs); err != nil {
			Log.Error(err)
		}
	}
}
