#! /bin/sh

set -e
script_dir=$(dirname $0)
abs_script_dir=$(readlink -f $script_dir)
theme_type=$1

ln -sfr $script_dir/config_${theme_type} $script_dir/config

if [ $abs_script_dir != $HOME/.config/mako ] && [ ! -v DONT_MODIFY_HOME ]; then
    ln -sfr $script_dir/config $HOME/.config/mako/config
fi

if [ ! -v DONT_RESTART ]; then
    echo -e $BLUE Restarting mako ... $RESET
    killall mako
    mako &
fi

echo -e $GREEN Successfully setup ${theme_type} theme for mako $RESET

