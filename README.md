# SamuraiOS

An Artix Linux configuration for Samurais 👹 and Ninjas 🥷

## Stage 1

1. Choose keyboard layout
```
loadkeys de
```

2. Partition the disk and create file system

3. Mount partitions

4. Connect to the internet

5. Install go and git on the host
```
pacman -Sy go git
```

1. Execute `go run stage.go 1`

## Stage 2

7. Execute `go run stage.go 2`

8. Make user able to use `sudo`
```
EDITOR=micro visudo # Uncomment the line with '%wheel ALL=(ALL:ALL) ALL'
usermod -aG wheel yourusername
```

9.  Reboot
```
exit
umount -R /mnt
reboot
```

## Stage 3

1.  After logging in execute `sudo dinitctl enable connmand`

2. Execute `go run stage.go 3`

3.  Logout

# Stage 4

12. After logging in again execute `go run stage.go 4`

13. Replace `User=samurai` in `/etc/sddm.conf.d/default.conf` with `User=yourusername`

14. Execute `sudo dinitctl enable sddm`

15. Done