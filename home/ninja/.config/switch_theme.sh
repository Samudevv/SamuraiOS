#! /bin/sh

set -e
script_dir=$(dirname $0)

$script_dir/setup_other_theme.sh

if [ $(readlink $script_dir/setup_other_theme.sh) == 'setup_light_theme.sh' ]; then
    ln -sfr $script_dir/setup_dark_theme.sh $script_dir/setup_other_theme.sh
else
    ln -sfr $script_dir/setup_light_theme.sh $script_dir/setup_other_theme.sh
fi

