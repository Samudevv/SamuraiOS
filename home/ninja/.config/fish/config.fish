# Set values
# Hide welcome message
set fish_greeting
set VIRTUAL_ENV_DISABLE_PROMPT 1
set -x MANPAGER "bat -l man -p"
set -x LANG de_AT.UTF-8

## Environment setup
# Apply .profile: use this to put fish compatible .profile stuff in
if test -f ~/.fish_profile
    source ~/.fish_profile
end

# Add ~/.local/bin to PATH
if test -d ~/.local/bin
    if not contains -- ~/.local/bin $PATH
        set -p PATH ~/.local/bin
    end
end

## Starship prompt
if status --is-interactive
    source ("/usr/bin/starship" init fish --print-full-init | psub)
end

## Functions
# Functions needed for !! and !$ https://github.com/oh-my-fish/plugin-bang-bang
function __history_previous_command
    switch (commandline -t)
        case "!"
            commandline -t $history[1]
            commandline -f repaint
        case "*"
            commandline -i !
    end
end

function __history_previous_command_arguments
    switch (commandline -t)
        case "!"
            commandline -t ""
            commandline -f history-token-search-backward
        case "*"
            commandline -i '$'
    end
end

if [ "$fish_key_bindings" = fish_vi_key_bindings ]
    bind -Minsert ! __history_previous_command
    bind -Minsert '$' __history_previous_command_arguments
else
    bind ! __history_previous_command
    bind '$' __history_previous_command_arguments
end

# Fish command history
function history
    builtin history --show-time='%F %T '
end

## Useful aliases
# Replace ls with exa
alias ls='exa -al --color=always --group-directories-first --icons' # preferred listing
alias ll='exa -l --color=always --group-directories-first --icons' # long format

# Replace some more things with better alternatives
alias cat='bat --style header --style rule --style snip --style changes --style header'

# Common use
alias untar='tar -zxvf '
alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'
alias .....='cd ../../../..'
alias ......='cd ../../../../..'
alias grep='grep --color=auto'
alias fgrep='grep -F --color=auto'
alias egrep='grep -E --color=auto'
alias pacman='pacman --color=auto'

# User configs
# Start ssh-agent through keychain to use ssh keys withoud passphrase
# eval (keychain --quiet --eval --agents ssh id_ed25519)

# Improve performance of command not found
function __fish_command_not_found_handler --on-event fish_command_not_found
    echo "fish: Unknown command '$argv'"
end

set -x --path GOPATH ~/go
if test -d ~/repos/SamuraiOS/home/ninja/go
    set -x --path GOPATH $GOPATH ~/repos/SamuraiOS/home/ninja/go
end

set -x CGO_CFLAGS '-g -O2 -Wdeprecated-declarations'
set -x CGO_CXXFLAGS '-g -O2 -Wdeprecated-declarations'

set -x --path PATH \
    $GOPATH/bin \
    $HOME/repos/Odin \
    $PATH

set -x BROWSER flatpak run net.waterfox.waterfox
set -x EDITOR nvim
set -x _JAVA_OPTIONS "-Dawt.useSystemAAFontSettings=on -Dswing.aatext=true"
# set -x _JAVA_OPTIONS "-Dawt.useSystemAAFontSettings=on -Dswing.aatext=true -Dswing.defaultlaf=com.sun.java.swing.plaf.gtk.GTKLookAndFeel"
set -x FZF_DEFAULT_COMMAND 'find . -type f ! -path "*/.*" -readable'
set -x MICRO_TRUECOLOR 1

alias fishconf="$EDITOR ~/.config/fish/config.fish"
alias sourcefish="source ~/.config/fish/config.fish"
alias hiddenconf="$EDITOR ~/.config/fish/hidden.fish"
alias nvimconf="$EDITOR ~/.config/nvim/init.lua"
alias qtileconf="$EDITOR ~/.config/qtile/config.py"
alias hyprconf="$EDITOR ~/.config/hypr/general.conf"
alias swayconf="$EDITOR ~/.config/sway/general"
alias gitcheckconf="git config user.name && git config user.email"
alias gitac="git add -A && git commit -m"
alias gits="git status"
alias gitd="git diff"
alias gitp="git push"
alias cmd="wine ~/.wine/drive_c/windows/system32/cmd.exe"
alias go-windows-amd64="GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go"
alias pacmanhist="cat /var/log/pacman.log | grep upgraded | less"
alias cpmingwdll="cp /usr/x86_64-w64-mingw32/bin/libwinpthread-1.dll /usr/x86_64-w64-mingw32/bin/libstdc++-6.dll /usr/x86_64-w64-mingw32/bin/libgcc_s_seh-1.dll ."
alias gitlsmod="git ls-files -om --exclude-standard"
alias microconf="$EDITOR $HOME/.config/micro/settings.json"
alias microkeys="$EDITOR $HOME/.config/micro/bindings.json"
alias microf="micro (fzf)"
alias e="$EDITOR"
alias ef="$EDITOR (fzf)"
alias code="gtk-launch codium-wayland"
alias clipboard="wl-copy --trim-newline"
alias icat="kitty +kitten icat"
alias odindemo="$EDITOR $HOME/repos/Odin/examples/demo/demo.odin"
alias neofetch="hyfetch"
alias update-mirrors="echo Updating mirrorlist ... && sudo cp /etc/pacman.d/mirrorlist /etc/pacman.d/mirrorlist.bak sudo reflector --latest 5 --sort rate --save /etc/pacman.d/mirrorlist && echo 'Update done!'"

if [ "$TERM" = xterm-kitty ]
    alias ssh="kitty +kitten ssh"
end

# Source hidden.fish if it exists
if test -f (status dirname)/hidden.fish
    source (status dirname)/hidden.fish
end

if test -n "$KEYCHAIN_SSH_KEYS"
    if status --is-interactive
        eval (keychain --eval --quiet --agents ssh $KEYCHAIN_SSH_KEYS)
    end
end
