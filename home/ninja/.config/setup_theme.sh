#! /bin/sh

set -e
script_dir=$(dirname $0)
theme_type=$1

BLUE='\033[0;34m'
GREEN='\033[0;32m'
RESET='\033[0m'

echo -e $BLUE Setting up ${theme_type} theme ... $RESET

BLUE=$BLUE GREEN=$GREEN RESET=$RESET $script_dir/waybar/setup_theme.sh $theme_type
BLUE=$BLUE GREEN=$GREEN RESET=$RESET $script_dir/hypr/setup_theme.sh $theme_type
BLUE=$BLUE GREEN=$GREEN RESET=$RESET $script_dir/wofi/setup_theme.sh $theme_type
BLUE=$BLUE GREEN=$GREEN RESET=$RESET $script_dir/mako/setup_theme.sh $theme_type

echo -e $GREEN Successfully set up $theme_type theme ... $RESET

