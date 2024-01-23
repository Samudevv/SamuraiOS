#! /bin/sh

set -e
script_dir=$(dirname $0)
abs_script_dir=$(readlink -f $script_dir)

ln -sfr $script_dir/wallpapers/dark_main.jpg $script_dir/wallpapers/main.jpg
ln -sfr $script_dir/wallpapers/dark_side.jpg $script_dir/wallpapers/side.jpg

if [ $abs_script_dir != $HOME/.config/hypr ] && [ ! -v DONT_MODIFY_HOME ]; then
    ln -sfr $script_dir/wallpapers/main.jpg $HOME/.config/hypr/wallpapers/main.jpg
    ln -sfr $script_dir/wallpapers/side.jpg $HOME/.config/hypr/wallpapers/side.jpg
fi

if [ ! -v DONT_RESTART ]; then
    echo -e $BLUE Restarting hyprpaper ... $RESET
    killall hyprpaper
    hyprpaper &
fi

echo -e $GREEN Successfully setup dark theme for hyprpaper $RESET

