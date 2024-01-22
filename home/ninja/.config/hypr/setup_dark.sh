#! /bin/sh

set -ex
script_dir=$(dirname $0)

rm -f $script_dir/wallpapers/main.jpg
ln -sr $script_dir/wallpapers/dark_main.jpg $script_dir/wallpapers/main.jpg

rm -f $script_dir/wallpapers/side.jpg
ln -sr $script_dir/wallpapers/dark_side.jpg $script_dir/wallpapers/side.jpg
