#! /bin/sh

echo Installing SamuraiOS with arguments '"' $* '"' ...

set -ex

timedatectl

pacman -Sy --needed --noconfirm reflector arch-install-scripts

mv /etc/pacman.d/mirrorlist /etc/pacman.d/mirrorlist.bak
reflector --latest 5 --sort rate --save /etc/pacman.d/mirrorlist.tmp
mv /etc/pacman.d/mirrorlist.tmp /etc/pacman.d/mirrorlist

pacstrap -K /mnt \
base \
linux-lts \
linux-firmware \
go \
git

genfstab -U /mnt >> /mnt/etc/fstab

printf "#! /bin/sh\nset -ex\ngit clone https://github.com/PucklaJ/SamuraiOS.git -b aarch64 --depth 1\ncd SamuraiOS\ngo run install.go 1 $*" > install0to1.sh
chmod +x install0to1.sh

cp install0to1.sh /mnt

arch-chroot /mnt /install0to1.sh
