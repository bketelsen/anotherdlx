#!/bin/bash
USER=`whoami`

set +e
lxc alias rm create
lxc alias rm flogin
lxc alias rm login
lxc alias rm zlogin
set -e
lxc alias add create 'launch dlxbase'
lxc alias add flogin "exec @ARGS@ --user 1000 --group 1000 --env HOME=/home/${USER} -- /usr/bin/fish --login"
lxc alias add login  "exec @ARGS@ --user 1000 --group 1000 --env HOME=/home/${USER} -- /bin/bash --login"
lxc alias add zlogin "exec @ARGS@ --user 1000 --group 1000 --env HOME=/home/${USER} -- /usr/bin/zsh --login"
