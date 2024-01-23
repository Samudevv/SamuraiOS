#! /bin/sh

set -e
script_dir=$(dirname $0)

BLUE='\033[0;34m'
GREEN='\033[0;32m'
RESET='\033[0m'

echo -e $BLUE Setting up dark theme ... $RESET

BLUE=$BLUE GREEN=$GREEN RESET=$RESET $script_dir/waybar/setup_dark.sh
BLUE=$BLUE GREEN=$GREEN RESET=$RESET $script_dir/hypr/setup_dark.sh

echo -e $GREEN Successfully set up dark theme ... $RESET

