#! /bin/sh

set -e
script_dir=$(dirname $0)
abs_script_dir=$(readlink -f $script_dir)
theme_type=$1

ln -sfr $script_dir/${theme_type}_style.css $script_dir/style.css

if [ $abs_script_dir != $HOME/.config/wofi ] && [ ! -v DONT_MODIFY_HOME ]; then
    ln -sfr $script_dir/style.css $HOME/.config/wofi/style.css
fi

wofi_name=$(printf 'wofi\n')
wofi_in_ps=$(ps -e | grep wofi | cut -f 9 -d ' ')

if [ ! -v DONT_RESTART ] && [[ $wofi_in_ps == *"$wofi_name"* ]]; then
    killall wofi
    wofi &
fi

