# bobr-status  [![Build Status](https://drone.dxas90.xyz/api/badges/dxas90/bobr-status/status.svg)](https://drone.dxas90.xyz/dxas90/bobr-status)

## Table of Contents

- [bobr-status](#Readme.md)
  - [Table of Contents](#table-of-contents)
- [Description](#description)
  - [Requirements](#requirements)
  - [Install](#install)
  - [Usage](#usage)
  - [Monitoring](#monitoring)
  - [Troubleshooting](#troubleshooting)
  - [License](#license)

## Description

This repo is to test access to the K8s resources from a pod itself.

## Requirements

You will need on your computer

* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#install-kubectl-on-linux)
* [Go](https://golang.org/dl/)

## Install

Here you should document any install steps required to use this module. You should consider documenting any pre-requisites in this section too.

```console
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```

## Usage

To create the binary just run the `make` command.

```sh
make
```

## Troubleshooting

you can use the `--debug` flag to check the commands

## License

![Apache2](https://img.shields.io/github/license/dxas90/bobr-status)
