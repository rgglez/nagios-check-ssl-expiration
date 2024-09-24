# check_ssl_expiration

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/nagios-check-ssl-expiration/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/nagios-check-ssl-expiration)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/nagios-check-ssl-expiration)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/nagios-check-ssl-expiration)](https://goreportcard.com/report/github.com/rgglez/nagios-check-ssl-expiration)
[![GitHub release](https://img.shields.io/github/release/rgglez/nagios-check-ssl-expiration.svg)](https://github.com/rgglez/gormcache/releases/)

**check_ssl_expiration** is a plugin for [Nagios](https://www.nagios.org) made in [Go](https://go.dev/). It retrieves the certificate from the given URL, and compares the [notValidAfter](https://clouddocs.f5.com/api/irules/X509__not_valid_after.html) field to the warn and crit parameters (if given) to see if the certificate is about to expire. 

## Usage

```bash
check_ssl_expiration --host=www.example.com --warn=10 --crit=5
```

This command checks the certificate for www.example.com (if any) and issues a normal warning if the certificate expires within 10 days, and a critical warning if it expires within 5 days.

### Command line parameters

* `--host` or `-H` specifies the URL to check. Example of valid values are: https://www.example.com, example.com or www.example.com/index.html.
* `--warn` or `-w` specifies the limit of days to issue a normal warning. Default value: 15 days.
* `--crit` or `-c` specifies the limit of days to issue a critical warning. Default value: 7 days.
* `--help` or `-h` shows the help.
* `--version` or `-v` shows the version of the program.

## Build and installation

### Build

To build the program, run:

```bash
make go
```

The executable will be created inside the ```build``` directory.

### Installation

Just copy the executable to your regular Nagios plugins directory.

## License

Copyright 2024 Rodolfo González González.

[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). Please read the [LICENSE](LICENSE.md) file.