#! /bin/sh

set -x

timedatectl

mv /etc/pacman.d/mirrorlist /etc/pacman.d/mirrorlist.bak
reflector --latest 5 --sort rate --save /etc/pacman.d/mirrorlist.tmp
mv /etc/pacman.d/mirrorlist.tmp /etc/pacman.d/mirrorlist

pacstrap -K /mnt \
base \
base-devel \
linux-lts \
linux-firmware \
micro \
neovim \
grub \
os-prober \
dhcpcd \
wpa_supplicant \
networkmanager \
reflector \
fzf \
whois \
go \
git

genfstab -U /mnt >> /mnt/etc/fstab
