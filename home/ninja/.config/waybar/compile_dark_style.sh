#! /bin/sh

set -e
script_dir=$(dirname $0)

sassc $script_dir/style.scss $script_dir/dark_style.css

echo -e $GREEN Sucessfully compiled dark stylesheets for waybar! $RESET

