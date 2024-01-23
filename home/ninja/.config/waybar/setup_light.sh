#! /bin/sh

set -e
script_dir=$(dirname $0)
abs_script_dir=$(readlink -f $script_dir)

ln -sfr $script_dir/light_style.css $script_dir/style.css
if [ $abs_script_dir != $HOME/.config/waybar ] && [ ! -v DONT_MODIFY_HOME ]; then
    ln -sfr $script_dir/style.css $HOME/.config/waybar/style.css
fi

echo -e $BLUE Restarting waybar ... $RESET
killall waybar
waybar &

echo -e $GREEN Sucessfully setup light style for waybar! $RESET

