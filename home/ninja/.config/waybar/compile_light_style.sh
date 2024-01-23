#! /bin/sh

set -e
script_dir=$(dirname $0)

sassc $script_dir/style.scss $script_dir/light_style.css

echo -e $GREEN Sucessfully compiled light stylesheets for waybar! $RESET

