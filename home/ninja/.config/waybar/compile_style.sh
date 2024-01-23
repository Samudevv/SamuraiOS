#! /bin/sh

set -e
script_dir=$(dirname $0)
theme_type=$1

sassc $script_dir/style.scss $script_dir/${theme_type}_style.css

echo -e $GREEN Sucessfully compiled $theme_type stylesheets for waybar! $RESET

