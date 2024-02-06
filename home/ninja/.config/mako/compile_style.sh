#! /bin/sh

set -e
script_dir=$(dirname $0)
theme_type=$1

printf "$GREEN "
gencolors $script_dir/../${theme_type}_colors.scss mako $script_dir/config.tmpl $script_dir/config_${theme_type}
printf $RESET
