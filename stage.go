package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var packages = []string{
	// Base packages to make the system work
	"base",
	"base-devel",
	"dinit",
	"elogind-dinit",
	"linux-zen",
	"linux-firmware",
	"micro",
	"grub",
	"os-prober",
	"dhcpcd",
	"wpa_supplicant",
	"connman-dinit",
	"connman-gtk",
	"go",

	// Packages for good working system
	"pipewire",
	"wireplumber",
	"pipewire-jack",
	"pipewire-pulse",
	"noto-fonts",
	"hyprland",
	"kitty",
	"dunst",
	"git",
	"polkit-kde-agent",
	"libnotify",
	"pavucontrol",
	"pcmanfm",
	"sddm-dinit",
}

var yayPackages = []string{
	"xdg-desktop-hyprland-git",
	"waybar-hyprland-no-systemd",
}

var dinitServices = []string{
	"connmand",
}

func input() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return text
}

func inputWithDefault(defaultValue string) string {
	text := input()
	if text == "" {
		return defaultValue
	}
	return text
}

func exe(command string) {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		fmt.Fprintln(os.Stderr, "No Command")
		os.Exit(1)
	}

	var args []string
	if len(words) > 1 {
		args = words[1:]
	}

	cmd := exec.Command(words[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println(command)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\"%s\" failed: %s\n", command, err)
		os.Exit(1)
	}
}

func exeAppendFile(command, filename string) {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		fmt.Fprintln(os.Stderr, "No Command")
		os.Exit(1)
	}

	var args []string
	if len(words) > 1 {
		args = words[1:]
	}

	file, err := os.OpenFile(filename, os.O_APPEND, 0770)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open \"%s\" for \"%s\": %s\n", filename, command, err)
		os.Exit(1)
	}
	defer file.Close()

	cmd := exec.Command(words[0], args...)
	cmd.Stdout = file
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println(fmt.Sprint(command, " >> ", filename))

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\"%s\" failed: %s\n", command, err)
		os.Exit(1)
	}
}

func exeToString(command string) string {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		fmt.Fprintln(os.Stderr, "No Command")
		os.Exit(1)
	}

	var args []string
	if len(words) > 1 {
		args = words[1:]
	}

	var builder strings.Builder

	cmd := exec.Command(words[0], args...)
	cmd.Stdout = &builder
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println(command)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\"%s\" failed: %s\n", command, err)
		os.Exit(1)
	}

	return builder.String()
}

func copyConfig(src string) {
	var dst string

	if strings.HasPrefix(src, "home/samurai") {
		dst = strings.Replace(src, "home/samurai", "$HOME", 1)
	} else {
		dst = "/" + src
	}

	dirname := filepath.Dir(dst)

	exe("mkdir -p " + dirname)
	if strings.HasPrefix(dst, "$HOME") {
		exe("cp " + src + " " + dst)
	} else {
		exe("sudo cp " + src + " " + dst)
	}
}

func main() {
	// Determine stage
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid Arguments")
		os.Exit(1)
	}

	stage, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse stage argument: ", err)
		os.Exit(1)
	}

	if stage == 1 {
		fmt.Println("Performing Stage 1 ...")

		// Update the system clock
		exe("dinitctl start ntpd")

		// Install rankmirrors and create better mirrorlist
		exe("pacman -S pacman-contrib")

		newMirrorlist := exeToString("rankmirrors -n 5 /etc/pacman.d/mirrorlist")
		{
			mirrorlist, err := os.Create("/etc/pacman.d/mirrorlist")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			mirrorlist.WriteString(newMirrorlist)
			mirrorlist.Close()
		}

		// Install base system
		fmt.Println("Is UEFI yes/no (no): ")
		isUEFIStr := inputWithDefault("no")
		isUEFI := strings.HasPrefix(strings.ToLower(isUEFIStr), "y")
		if isUEFI {
			packages = append(packages, "efibootmgr")
		}
		exe("basestrap /mnt " + strings.Join(packages, " "))

		// fstabgen
		exeAppendFile("fstabgen -U /mnt", "/mnt/etc/fstab")

		// chroot into system
		exe("artix-chroot /mnt")

		fmt.Println("Stage 1 Done")
	} else if stage == 2 {
		fmt.Println("Performing Stage 2 ...")

		// set the time zone
		fmt.Println("Region (Europe): ")
		region := inputWithDefault("Europe")
		fmt.Println("City (Vienna): ")
		city := inputWithDefault("Vienna")

		exe(fmt.Sprint("ln -sf /usr/share/zoneinfo/", region, "/", city, " /etc/localtime"))

		exe("hwclock --systohc")

		// Set locale
		fmt.Println("Locale (comma seperated) (de_AT.UTF-8, en_GB.UTF-8): ")
		localesStr := inputWithDefault("de_AT.UTF-8, en_GB.UTF-8")
		locales := strings.Split(localesStr, ",")
		for i := 0; i < len(locales); i++ {
			locales[i] = strings.TrimSpace(locales[i])
		}

		{
			localeGen, err := os.Open("/etc/locale.gen")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			var localeBuilder strings.Builder

			localeScanner := bufio.NewScanner(localeGen)

			for localeScanner.Scan() {
				line := localeScanner.Text()
				var isWanted bool
				for _, l := range locales {
					if strings.HasPrefix(line, "#"+l) {
						isWanted = true
						break
					}
				}

				if isWanted {
					localeBuilder.WriteString(strings.TrimPrefix(line, "#") + "\n")
				} else {
					localeBuilder.WriteString(line + "\n")
				}
			}

			localeGen.Close()

			localeGen, err = os.Create("/etc/locale.gen")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			localeGen.WriteString(localeBuilder.String())
			localeGen.Close()
		}

		exe("locale-gen")
		exeAppendFile("echo export LANG=\""+locales[0]+"\"", "/etc/locale.conf")

		// Boot Loader
		fmt.Println("Is UEFI yes/no (no): ")
		isUEFIStr := inputWithDefault("no")
		isUEFI := strings.HasPrefix(strings.ToLower(isUEFIStr), "y")

		if isUEFI {
			exe("grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=grub")
		} else {
			fmt.Println("Device: ")
			device := input()

			exe("grub-install --recheck " + device)
		}

		exe("grub-mkconfig -o /boot/grub/grub.cfg")

		// Passwords, Username and Hostname
		fmt.Println("Enter root password")
		exe("passwd")

		fmt.Println("Username (ninja): ")
		userName := inputWithDefault("ninja")
		exe("useradd -m " + userName)

		fmt.Println("Enter password for " + userName)
		exe("passwd " + userName)

		fmt.Println("Hostname (samurai): ")
		hostName := inputWithDefault("samurai")

		{
			hostNameFile, err := os.Create("/etc/hostname")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			hostNameFile.WriteString(hostName)
			hostNameFile.Close()
		}

		{
			hosts, err := os.Create("/etc/hosts")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			hosts.WriteString("127.0.0.1\tlocalhost\n")
			hosts.WriteString("::1\tlocalhost\n")
			hosts.WriteString("127.0.1.1\t" + hostName + ".localdomain\t" + hostName + "\n")
		}

		fmt.Println("Stage 2 Done")
	} else if stage == 3 {
		fmt.Println("Performing Stage 3 ...")

		// Activate connman
		exe("sudo dinitctl enable connmand")

		// Installing yay
		exe("mkdir -p $HOME/repos/yay")
		exe("git clone https://aur.archlinux.org/yay.git $HOME/repos/yay")
		exe("makepkg -sip $HOME/repos/yay/PKGBUILD")

		// Install yay packages
		exe("yay -S " + strings.Join(yayPackages, " "))

		// Remove unneeded packages
		exe("yay -Rnsdd xdg-desktop-portal-gnome xdg-desktop-portal-gtk xdg-desktop-portal-kde xdg-desktop-portal-wlr")

		// Install dinit-userservd
		exe("mkdir $HOME/repos/dinit-userservd")
		exe("git clone https://github.com/Xynonners/dinit-userservd.git $HOME/repos/dinit-userservd")
		exe("makepkg -sip $HOME/repos/dinit-userservd/PKGBUILD")
		exe("sudo dinitctl enable dinit-userservd")

		exeAppendFile("sudo echo session optional pam_dinit_userservd.so", "/etc/pam.d/system-login")

		// Copy configuration files
		copyConfig("etc/sddm.conf.d/default.conf")

		copyConfig("home/samurai/.config/dinit.d/dunst")
		copyConfig("home/samurai/.config/dinit.d/pipewire")
		copyConfig("home/samurai/.config/dinit.d/pipewire-pulse")
		copyConfig("home/samurai/.config/dinit.d/wireplumber")

		copyConfig("home/samurai/.config/hypr/hyprland.conf")
		copyConfig("home/samurai/.config/micro/bindings.json")
		copyConfig("home/samurai/.config/micro/settings.json")
		copyConfig("home/samurai/.config/waybar/config")
		copyConfig("home/samurai/.config/waybar/style.css")

		fmt.Println("Stage 3 Done")
	} else if stage == 4 {
		fmt.Println("Performing Stage 4")

		// Activate dinit user services
		exe("dinitctl enable dunst")
		exe("dinitctl enable pipewire")
		exe("dinitctl enable pipewire-pulse")
		exe("dinitctl enable wireplumber")

		fmt.Println("Stage 4 Done")
	} else {
		fmt.Fprintln(os.Stderr, "Invalid Stage", stage)
		os.Exit(1)
	}
}
