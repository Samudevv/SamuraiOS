#! /bin/sh

set -e
script_dir=$(dirname $0)
abs_script_dir=$(readlink -f $script_dir)
theme_type=$1

ln -sfr $script_dir/wallpapers/${theme_type}_main.jpg $script_dir/wallpapers/main.jpg
ln -sfr $script_dir/wallpapers/${theme_type}_side.jpg $script_dir/wallpapers/side.jpg

if [ $abs_script_dir != $HOME/.config/hypr ] && [ ! -v DONT_MODIFY_HOME ]; then
    ln -sfr $script_dir/wallpapers/main.jpg $HOME/.config/hypr/wallpapers/main.jpg
    ln -sfr $script_dir/wallpapers/side.jpg $HOME/.config/hypr/wallpapers/side.jpg
fi

if [ ! -v DONT_RESTART ]; then
    echo -e $BLUE Restarting hyprpaper ... $RESET
    killall hyprpaper
    hyprpaper &
fi

echo -e $GREEN Successfully setup ${theme_type} theme for hyprpaper $RESET

