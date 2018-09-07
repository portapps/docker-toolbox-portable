//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	. "github.com/portapps/portapps"
)

type config struct {
	Machine machine `json:"machine"`
}

type machine struct {
	Name         string `json:"name"`
	HostCIDR     string `json:"host_cidr"`
	CPU          int    `json:"cpu"`
	Ram          int    `json:"ram"`
	Disk         int    `json:"disk"`
	SharedName   string `json:"share_name"`
	OnExitStop   bool   `json:"on_exit_stop"`
	OnExitRemove bool   `json:"on_exit_remove"`
}

func init() {
	Papp.ID = "docker-toolbox-portable"
	Papp.Name = "Docker Toolbox"
	Init()

	SetConsoleTitle(fmt.Sprintf("%s Portable", Papp.Name))
}

func main() {
	var err error
	var cfg config
	var oldCfg config

	Papp.AppPath = AppPathJoin("app")
	Papp.DataPath = AppPathJoin("data")
	Papp.Process = PathJoin(Papp.AppPath, "git", "bin", "bash.exe")
	Papp.Args = []string{"--login", "-i", PathJoin(Papp.AppPath, "start.sh")}
	Papp.WorkingDir = Papp.AppPath

	sharedPath := CreateFolder(PathJoin(Papp.DataPath, "shared"))
	storagePath := CreateFolder(PathJoin(Papp.DataPath, "storage"))
	cfgPath := PathJoin(Papp.Path, fmt.Sprintf("%s.json", Papp.ID))
	cfgDefault := config{
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

	if err = createConfig(cfgPath, cfgDefault, &cfg, oldCfg); err != nil {
		Log.Fatal(err)
	}
	Log.Infof("Config: %v", cfg)

	postInstallGit := PathJoin(Papp.AppPath, "git", "post-install.bat")
	if _, err := os.Stat(postInstallGit); err == nil {
		Log.Info("Initializing git...")
		if err = QuickExecCmd("cmd", []string{"/k", postInstallGit}); err != nil {
			Log.Errorf("Cannot initialize git: %v", err)
		}
	}

	Log.Info("Setting environment...")
	os.Setenv("MACHINE_NAME", cfg.Machine.Name)
	os.Setenv("MACHINE_HOST_CIDR", cfg.Machine.HostCIDR)
	os.Setenv("MACHINE_CPU", strconv.Itoa(cfg.Machine.CPU))
	os.Setenv("MACHINE_RAM", strconv.Itoa(cfg.Machine.Ram))
	os.Setenv("MACHINE_DISK", strconv.Itoa(cfg.Machine.Disk))
	os.Setenv("MACHINE_STORAGE_PATH", FormatUnixPath(storagePath))
	os.Setenv("MACHINE_SHARED_NAME", cfg.Machine.SharedName)
	os.Setenv("MACHINE_SHARED_PATH", sharedPath)

	Log.Info("Adding app to PATH...")
	os.Setenv("PATH", fmt.Sprintf("%s;%s", Papp.AppPath, os.Getenv("PATH")))

	Log.Info("Starting up the shell... ")
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   Papp.AppPath,
		Sys: &syscall.SysProcAttr{
			CmdLine: fmt.Sprintf(` --login -i "%s"`, PathJoin(Papp.AppPath, "start.sh")),
		},
	}
	proc, err := os.StartProcess(Papp.Process, []string{}, &pa)
	if err != nil {
		Log.Fatal(err)
	}
	if _, err = proc.Wait(); err != nil {
		Log.Fatal(err)
	}

	var exitArgs []string
	if cfg.Machine.OnExitRemove {
		exitArgs = []string{"rm", "-f", cfg.Machine.Name}
	} else if cfg.Machine.OnExitStop {
		exitArgs = []string{"stop", cfg.Machine.Name}
	}
	if len(exitArgs) > 0 {
		if err = QuickExecCmd("docker-machine", exitArgs); err != nil {
			Log.Error(err)
		}
	}
}

func createConfig(confPath string, defaultConf config, conf *config, oldConf config) error {
	// Create config if not exists
	if _, err := os.Stat(confPath); err != nil {
		Log.Info("defaultJSON")
		defaultConfJSON, err := json.MarshalIndent(defaultConf, "", "  ")
		if err != nil {
			return err
		}
		Log.Info("Write")
		err = ioutil.WriteFile(confPath, defaultConfJSON, 0644)
		if err != nil {
			return err
		}
	}

	// Load current config
	raw, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, &oldConf)
	if err != nil {
		return err
	}

	// Merge config
	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return err
	}

	// Write config
	confJSON, _ := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(confPath, confJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
