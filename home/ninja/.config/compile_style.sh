#! /bin/sh

set -e
script_dir=$(dirname $0)

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RESET='\033[0m'

ln -sfr $script_dir/light_colors.scss $script_dir/colors.scss

echo -e $BLUE Compiling all light stylesheets ... $RESET
GREEN=$GREEN BLUE=$BLUE RESET=$RESET $script_dir/waybar/compile_style.sh light
GREEN=$GREEN BLUE=$BLUE RESET=$RESET $script_dir/wofi/compile_style.sh light

ln -sfr $script_dir/dark_colors.scss $script_dir/colors.scss

echo -e $BLUE Compiling all dark stylesheets ... $RESET
GREEN=$GREEN BLUE=$BLUE RESET=$RESET $script_dir/waybar/compile_style.sh dark
GREEN=$GREEN BLUE=$BLUE RESET=$RESET $script_dir/wofi/compile_style.sh dark

rm $script_dir/colors.scss

echo -e $GREEN Sucessfully compiled all stylesheets ... $RESET
