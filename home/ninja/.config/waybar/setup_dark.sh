#! /bin/sh

set -ex
script_dir=$(dirname $0)

rm -f $script_dir/style.css
ln -sr $script_dir/dark_style.css $script_dir/style.css

echo Sucessfully setup dark style for waybar!

