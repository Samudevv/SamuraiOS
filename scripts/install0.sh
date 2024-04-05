#! /bin/sh

echo Installing SamuraiOS with arguments '"' $* '"' ...

set -ex

printf "#! /bin/sh\nset -ex\npacman-key --init\npacman-key --populate archlinuxarm\npacman -Sy --needed --noconfirm go git pacman-contrib parallel\ngit clone https://github.com/Samudevv/SamuraiOS.git -b aarch64 --depth 1\ncd SamuraiOS\ngo run install.go 1 $*\n" > install0to1.sh
chmod +x install0to1.sh
echo Install script contents
cat install0to1.sh

timedatectl

pacman -Sy --needed --noconfirm arch-install-scripts wget

arch_arm_tar=ArchLinuxARM-aarch64-latest.tar.gz

if [ -e "$arch_arm_tar" ]; then
  echo "Skipping download of $arch_arm_tar since it has already been downloaded"
else
  echo "Downloading $arch_arm_tar ..."
  wget http://os.archlinuxarm.org/os/$arch_arm_tar
  wget http://os.archlinuxarm.org/os/$arch_arm_tar.sig
  gpg --keyserver keyserver.ubuntu.com --recv-keys 68B3537F39A313B3E574D06777193F152BDBE6A6
  gpg --verify $arch_arm_tar.sig
fi

bsdtar -xpf $arch_arm_tar -C /mnt
cp install0to1.sh /mnt

genfstab -U /mnt >> /mnt/etc/fstab

arch-chroot /mnt /install0to1.sh
