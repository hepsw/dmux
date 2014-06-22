dmux
====

docker + tmux = dmux !

## Synopsis

`dmux` create docker container environment with tmux pane. You can try some commands there and destroy it. You can also save your tried status as Docker image.

## Demo

![](http://deeeet.com/writing/images/post/dmux.gif)

## Requirement

- [Docker v0.12.0 or later](http://www.docker.com/)
- [Tmux](http://tmux.sourceforge.net/)

## Command

Before, you need to create new tmux session.

First, you need to create new container:

```bash
$ dmux init
```

After that new tmux pane is created and docker container environment is there. You can try some command or enjoy Docker environment.

If you want pause container, from anoher window:

```bash
$ dmux stop
```

This will pause container, it won't send `SIGTERM` and `SIGKILL`. It will freeze cotainer process with cgroup freezer. **You can restart process from the point where you exec stop**.

If you want to unpauze container and attach it again:

```bash
$ dmux start
```

If you delete container:

```bash
$ dmux delete
```

You can save your container as image:

```bash
$ dmux save [Image]
```

## Installation

Binary files are distributed.

[ ![Download](https://api.bintray.com/packages/tcnksm/dmux/dmux/images/download.png) ](http://dl.bintray.com/tcnksm/dmux/)

## Development

You can prepare it standard way:

```bash
$ go get install github.com/tcnksm/dmux
```

## VS.

- [jpetazzo/critmux](https://github.com/jpetazzo/critmux) - critmux uses [CRIU]() for stop container. `dmux` uses docker's `pause` and `unpause` command.

## Author

[tcnksm](https://github.com/tcnksm)
