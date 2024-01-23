#! /bin/sh

set -e
script_dir=$(dirname $0)
abs_script_dir=$(readlink -f $script_dir)
theme_type=$1

ln -sfr $script_dir/${theme_type}_style.css $script_dir/style.css
if [ $abs_script_dir != $HOME/.config/waybar ] && [ ! -v DONT_MODIFY_HOME ]; then
    ln -sfr $script_dir/style.css $HOME/.config/waybar/style.css
fi

if [ ! -v DONT_RESTART ]; then
    echo -e $BLUE Restarting waybar ... $RESET
    killall waybar
    waybar &
fi

echo -e $GREEN Sucessfully setup $theme_type style for waybar! $RESET

