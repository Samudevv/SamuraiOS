package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
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
	"playerctl",
	"mpv",
	"libmpeg2",
}

var yayPackages = []string{
	"xdg-desktop-portal-hyprland-git",
	"waybar-hyprland-no-systemd",
	"connman-gtk",
}

func main() {
	// Determine stage
	if len(os.Args) != 2 {
		logError("Invalid Arguments")
		os.Exit(1)
	}

	stage, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		logError("Failed to parse stage argument: ", err)
		os.Exit(1)
	}

	if stage == 1 {
		logInfo("Performing Stage 1 ...")

		// Update the system clock
		exe("dinitctl start ntpd")

		logInfo("Refreshing new mirrorlist ...")
		// Install rankmirrors and create better mirrorlist
		exe("pacman -S --noconfirm pacman-contrib")

		newMirrorlist := exeToString("rankmirrors -n 5 /etc/pacman.d/mirrorlist")
		{
			mirrorlist, err := os.Create("/etc/pacman.d/mirrorlist")
			if err != nil {
				logError(err)
				os.Exit(1)
			}

			mirrorlist.WriteString(newMirrorlist)
			mirrorlist.Close()
		}

		// Install base system
		logInfo("Installing base packages ...")
		prompt("Is UEFI yes|no (no)")
		isUEFIStr := inputWithDefault("no")
		isUEFI := strings.HasPrefix(strings.ToLower(isUEFIStr), "y")
		if isUEFI {
			packages = append(packages, "efibootmgr")
		}
		exe("basestrap /mnt " + strings.Join(packages, " "))

		// fstabgen
		exeAppendFile("fstabgen -U /mnt", "/mnt/etc/fstab")

		// chroot into system
		logInfo("Stage 1 Done")
		logInfo("Now clone the repo again and execute \"go run stage.go 2\"")

		exe("artix-chroot /mnt")
	} else if stage == 2 {
		logInfo("Performing Stage 2 ...")

		// set the time zone
		logInfo("Setting locale ...")
		prompt("Region (Europe)")
		region := inputWithDefault("Europe")
		prompt("City (Vienna)")
		city := inputWithDefault("Vienna")

		exe(fmt.Sprint("ln -sf /usr/share/zoneinfo/", region, "/", city, " /etc/localtime"))

		exe("hwclock --systohc")

		// Set locale
		prompt("Locale (comma seperated) (de_AT.UTF-8, en_GB.UTF-8)")
		localesStr := inputWithDefault("de_AT.UTF-8, en_GB.UTF-8")
		locales := strings.Split(localesStr, ",")
		for i := 0; i < len(locales); i++ {
			locales[i] = strings.TrimSpace(locales[i])
		}

		{
			localeGen, err := os.Open("/etc/locale.gen")
			if err != nil {
				logError(err)
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
				logError(err)
				os.Exit(1)
			}

			localeGen.WriteString(localeBuilder.String())
			localeGen.Close()
		}

		exe("locale-gen")
		exeAppendFile("echo export LANG=\""+locales[0]+"\"", "/etc/locale.conf")

		// Boot Loader
		logInfo("Installing boot loader (grub) ...")
		prompt("Is UEFI yes|no (no)")
		isUEFIStr := inputWithDefault("no")
		isUEFI := strings.HasPrefix(strings.ToLower(isUEFIStr), "y")

		if isUEFI {
			exe("grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=grub")
		} else {
			prompt("Device")
			var device string
			for device == "" {
				device = input()
			}

			exe("grub-install --recheck " + device)
		}

		exe("grub-mkconfig -o /boot/grub/grub.cfg")

		// Passwords, Username and Hostname
		logInfo("Enter root password")
		exe("passwd")

		prompt("Username (ninja)")
		userName := inputWithDefault("ninja")
		exe("useradd -m " + userName)

		logInfo("Enter password for ", userName)
		exe("passwd " + userName)

		prompt("Hostname (samurai)")
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

		exe("usermod -aG wheel " + userName)

		logInfo("Stage 2 Done")
		logInfo("Now reboot and boot into new drive and execute \"sudo dinitctl enable connmand\" to activate the network daemon. After that reconnect to the internet and execute \"go run stage.go 3\"")
	} else if stage == 3 {
		logInfo("Performing Stage 3 ...")

		homeDir, _ := os.UserHomeDir()
		curDir, _ := os.Getwd()

		// Installing yay
		logInfo("Installing yay ...")
		exe("mkdir -p " + filepath.Join(homeDir, "/repos/yay"))
		exe("git clone https://aur.archlinux.org/yay.git " + filepath.Join(homeDir, "/repos/yay"))
		os.Chdir(filepath.Join(homeDir, "/repos/yay"))
		exe("makepkg -si --noconfirm")

		// Install yay packages
		logInfo("Installing AUR packages ...")
		exe("yay -S --noconfirm " + strings.Join(yayPackages, " "))

		// Remove unneeded packages
		exeDontCare("sudo pacman -Rnsdd --noconfirm xdg-desktop-portal-gnome")
		exeDontCare("sudo pacman -Rnsdd --noconfirm xdg-desktop-portal-gtk")
		exeDontCare("sudo pacman -Rnsdd --noconfirm xdg-desktop-portal-kde")
		exeDontCare("sudo pacman -Rnsdd --noconfirm xdg-desktop-portal-wlr")

		// Install dinit-userservd
		logInfo("Installing dinit user service ...")
		exe("mkdir " + filepath.Join(homeDir, "/repos/dinit-userservd"))
		exe("git clone https://github.com/Xynonners/dinit-userservd.git " + filepath.Join(homeDir, "/repos/dinit-userservd"))
		os.Chdir(filepath.Join(homeDir, "/repos/dinit-userservd"))
		exe("makepkg -si --noconfirm")
		exe("sudo dinitctl enable dinit-userservd")
		os.Chdir(curDir)
		exeArgs("sudo", "go", "run", "append.go", "echo session optional pam_dinit_userservd.so", "/etc/pam.d/system-login")

		// Copy configuration files
		logInfo("Copying configuration files ...")
		copyConfig("etc/sddm.conf.d/default.conf")

		exe("mkdir -p " + filepath.Join(homeDir, ".config/dinit.d/boot.d"))
		copyConfig("home/samurai/.config/dinit.d/pipewire")
		copyConfig("home/samurai/.config/dinit.d/pipewire-pulse")

		copyConfig("home/samurai/.config/hypr/hyprland.conf")
		copyConfig("home/samurai/.config/micro/bindings.json")
		copyConfig("home/samurai/.config/micro/settings.json")
		copyConfig("home/samurai/.config/waybar/config")
		copyConfig("home/samurai/.config/waybar/style.css")

		logInfo("Stage 3 Done")
		logInfo("Now logout and login again and execute \"go run stage.go 4\"")
	} else if stage == 4 {
		logInfo("Performing Stage 4 ...")

		// Activate dinit user services
		logInfo("Activating dinit user services ...")
		exe("dinitctl enable pipewire")
		exe("dinitctl enable pipewire-pulse")

		currentUser, _ := user.Current()
		exe("sudo go run replace.go /etc/sddm.conf.d/default.conf samurai " + currentUser.Username)

		logInfo("Stage 4 Done")
		logInfo("The final step is to enable sddm which will launch you into hyprland: \"sudo dinitctl enable sddm\"")
	} else if stage == 255 {
		// Testing
		logInfo("Performing Tests ...")

		exe("mkdir -p /tmp/stage_test")
		exe("git clone https://github.com/PucklaJ/SamuraiOS.git /tmp/stage_test")
		logError("Task not failed successfully")
		exe("rm -rf /tmp/stage_test")

		prompt("Create which folder (/tmp/another/samurai)")
		folderName := inputWithDefault("/tmp/another/samurai")
		exe("mkdir -p " + folderName)

		exeAppendFile("echo Hello World", "/tmp/hello_samurai")
		exe("go run replace.go /tmp/hello_samurai Hello Greetings")

		logInfo("Tests Done")
	} else {
		logError("Invalid Stage ", stage)
		os.Exit(1)
	}
}

func input() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(text, "\n")
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
		logError("No Command")
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

	logScript(command)

	if err := cmd.Run(); err != nil {
		logError("\"", command, "\" failed: ", err)
		os.Exit(1)
	}
}

func exeDontCare(command string) {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		logError("No Command")
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

	logScript(command)

	cmd.Run()
}

func exeArgs(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	logScript(strings.Join(args, " "))

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\"%s\" failed: %s\n", strings.Join(args, " "), err)
		os.Exit(1)
	}
}

func exeAppendFile(command, filename string) {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		logError("No Command")
		os.Exit(1)
	}

	var args []string
	if len(words) > 1 {
		args = words[1:]
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		logError("Failed to open \"", filename, "\" for \"", command, "\": ", err)
		os.Exit(1)
	}
	defer file.Close()

	cmd := exec.Command(words[0], args...)
	cmd.Stdout = file
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	logScript(command, " >> ", filename)

	if err := cmd.Run(); err != nil {
		logError("\"", command, "\" failed: ", err)
		os.Exit(1)
	}
}

func exeToString(command string) string {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		fmt.Fprintln(os.Stderr, "No Command")
		logError("No Command")
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

	logScript(command)

	if err := cmd.Run(); err != nil {
		logError("\"", command, "\" failed: ", err)
		os.Exit(1)
	}

	return builder.String()
}

func copyConfig(src string) {
	var dst string

	homeDir, _ := os.UserHomeDir()

	if strings.HasPrefix(src, "home/samurai") {
		dst = strings.Replace(src, "home/samurai", homeDir, 1)
	} else {
		dst = "/" + src
	}

	dirname := filepath.Dir(dst)

	if strings.HasPrefix(dst, homeDir) {
		exe("mkdir -p " + dirname)
		exe("cp " + src + " " + dst)
	} else {
		exe("sudo mkdir -p " + dirname)
		exe("sudo cp " + src + " " + dst)
	}
}

func logInfo(msg ...any) {
	msgStr := fmt.Sprint(msg...)
	fmt.Printf("\n\n\033[30;46m[INFO]\033[0;33m %s\033[0m\n\n", msgStr)
}

func logError(msg ...any) {
	msgStr := fmt.Sprint(msg...)
	fmt.Fprintf(os.Stderr, "\n\n\033[30;41m[ERROR]\033[0;33m %s\033[0m\n\n", msgStr)
}

func logScript(msg ...any) {
	msgStr := fmt.Sprint(msg...)
	fmt.Printf("\033[30;47m[CALL]\033[0;33m %s\033[0m\n", msgStr)
}

func prompt(msg ...any) {
	msgStr := fmt.Sprint(msg...)
	fmt.Printf("\033[30;47m[PROMPT]\033[0;33m %s: \033[0m", msgStr)
}
