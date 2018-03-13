<p align="center"><a href="https://github.com/portapps/docker-toolbox-portable" target="_blank"><img width="100" src="https://github.com/portapps/docker-toolbox-portable/blob/master/res/papp.png"></a></p>

<p align="center">
  <a href="https://github.com/portapps/docker-toolbox-portable/releases/latest"><img src="https://img.shields.io/github/release/portapps/docker-toolbox-portable.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/portapps/docker-toolbox-portable/releases/latest"><img src="https://img.shields.io/github/downloads/portapps/docker-toolbox-portable/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://ci.appveyor.com/project/portapps/docker-toolbox-portable"><img src="https://img.shields.io/appveyor/ci/portapps/docker-toolbox-portable.svg?style=flat-square" alt="AppVeyor"></a>
  <a href="https://goreportcard.com/report/github.com/portapps/docker-toolbox-portable"><img src="https://goreportcard.com/badge/github.com/portapps/docker-toolbox-portable?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/portapps/docker-toolbox-portable"><img src="https://img.shields.io/codacy/grade/439e341359d14857a0ee82f593a995e4.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://beerpay.io/portapps/portapps"><img src="https://img.shields.io/beerpay/portapps/portapps.svg?style=flat-square" alt="Beerpay"></a>
  <a href="https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=WQD7AQGPDEPSG"><img src="https://img.shields.io/badge/donate-paypal-7057ff.svg?style=flat-square" alt="Donate Paypal"></a>
</p>

## About

[Docker Toolbox](https://docs.docker.com/toolbox/overview/) portable app made with 🚀 [Portapps](https://github.com/portapps).<br />
Tested on Windows 7, Windows 8.1 and Windows 10.

## Requirements

* [VirtualBox](https://www.virtualbox.org/) needs to be installed

## Installation

There are different kinds of artifacts :

* `docker-toolbox-portable-win{32,64}-x.x.x-x-setup.exe` : Full portable release of Docker Toolbox as a setup. **Recommended way**!
* `docker-toolbox-portable-win{32,64}-x.x.x-x.7z` : Full portable release of Docker Toolbox as a 7z archive.
* `docker-toolbox-portable-win{32,64}.exe` : Only the portable binary (must be renamed `docker-toolbox-portable.exe`)
* `DockerToolbox-x.x.x.exe` : The original setup from the [official repository](https://github.com/docker/toolbox/releases).

### Fresh installation

Install `docker-toolbox-portable-win{32,64}-x.x.x-x-setup.exe` where you want then run `docker-toolbox-portable.exe`.

### App already installed

If you have already installed Docker Toolbox from the original setup :

* Stop the virtual machine in VirtualBox
* Move data located in `%USERPROFILE%\.docker\*` to `data\storage` folder

Run `docker-toolbox-portable.exe` and then you can [remove](https://support.microsoft.com/en-us/instantanswers/ce7ba88b-4e95-4354-b807-35732db36c4d/repair-or-remove-programs) Docker Toolbox from your computer.

### Upgrade

For an upgrade, simply download and install the [latest setup](https://github.com/portapps/docker-toolbox-portable/releases/latest).

## Configuration

A configuration file called `docker-toolbox-portable.json` is generated at first launch and can be customized :

```json
{
  "machine": {
    "name": "default",
    "host_cidr": "192.168.99.1/24",
    "cpu": 1,
    "ram": 1024,
    "disk": 20000,
    "share_name": "shared",
    "on_exit_stop": false,
    "on_exit_remove": false
  }
}
```

* `name` : Name of the virtual machine (default `default`)
* `host_cidr` : Specify the Host Only CIDR (default `192.168.99.1/24`)
* `cpu` : Number of CPUs for the machine (-1 to use the number of CPUs available ; default `1`)
* `ram` : Size of memory for host in MB (default `1024`)
* `disk` : Size of disk for host in MB (default `20000`)
* `share_name` : Name of the mounted directory (in `data\shared`) to use as volume (default `shared`)
* `on_exit_stop` : Stop the virtual machine on exit
* `on_exit_remove` : Remove the virtual machine on exit

### Mount a volume

The directory for volume persistence is located in `data\shared`.<br />
The share name can be customized in the configuration file and if you kept the default one (`shared`) you can mount a volume for persistence and fully portable : `-v /shared/matomo:/data`.

## How can i help ?

All kinds of contributions are welcomed :raised_hands:!<br />
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:<br />
But we're not gonna lie to each other, I'd rather you buy me a beer or two :beers:!

[![Beerpay](https://beerpay.io/portapps/portapps/badge.svg?style=beer-square)](https://beerpay.io/portapps/portapps)
or [![Paypal](https://cdn.rawgit.com/portapps/portapps/master/res/paypal.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=WQD7AQGPDEPSG)

## License

MIT. See `LICENSE` for more details.<br />
Rocket icon credit to [Squid Ink](http://thesquid.ink).
