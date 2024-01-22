#! /bin/sh

set -ex
script_dir=$(dirname $0)

rm -f $script_dir/style.css
ln -sr $script_dir/light_style.css $script_dir/style.css

echo Sucessfully setup light style for waybar!

