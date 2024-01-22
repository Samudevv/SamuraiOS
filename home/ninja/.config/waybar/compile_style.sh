#! /bin/sh

set -ex
script_dir=$(dirname $0)

rm -f $script_dir/colors.scss
ln -sr $script_dir/dark_colors.scss $script_dir/colors.scss

sassc $script_dir/style.scss $script_dir/dark_style.css

rm $script_dir/colors.scss
ln -sr $script_dir/light_colors.scss $script_dir/colors.scss

sassc $script_dir/style.scss $script_dir/light_style.css

rm $script_dir/colors.scss

echo Sucessfully compiled dark and light stylesheets!

