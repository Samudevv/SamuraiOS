#! /bin/sh

set -e
script_dir=$(dirname $0)

generate_setup_other_theme() {
    echo Generating setup_other_theme.sh for $1
    printf '#! /bin/sh\nset -e\nscript_dir=$(dirname $0)\n$script_dir/setup_theme.sh '$1 > $script_dir/setup_other_theme.sh
    chmod +x $script_dir/setup_other_theme.sh
}

if [[ ! -e $script_dir/setup_other_theme.sh ]]; then
    generate_setup_other_theme light
fi

$script_dir/setup_other_theme.sh

if [[ $(cat $script_dir/setup_other_theme.sh) == *"light"* ]]; then
    generate_setup_other_theme dark
else
    generate_setup_other_theme light
fi

