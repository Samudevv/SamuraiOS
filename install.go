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

// Packages that can be installed from the artix repos and the program basestrap
var basestrapPackages = []string{
	// Base packages to make the system work
	"base",
	"base-devel",
	"dinit",
	"elogind-dinit",
	"linux-lts",
	"linux-firmware",
	"micro",
	"grub",
	"os-prober",
	"dhcpcd",
	"wpa_supplicant",
	"connman-dinit",
	"pacman-contrib",
	"parallel",

	// Packages for working graphical system with audio
	"go",
	"fish",
	"pipewire",
	"pipewire-pulse",
	"pipewire-jack",
	"pipewire-alsa",
	"wireplumber",
	"libmpeg2",
	"sddm-dinit",
	"hyprland",
	"polkit-gnome",
	"man",
	"git",
	"dunst",
	"libnotify",
	"kitty",
	"noto-fonts",
	"noto-fonts-emoji",
	"noto-fonts-cjk",
	"nerd-fonts",
	"playerctl",
	"qt5ct",
	"gnome-keyring",
	"grim",
	"slurp",
	"wl-clipboard",
	"fcitx5",
	"fcitx5-qt",
	"fcitx5-gtk",
	"fcitx5-mozc",
	"fcitx5-configtool",
	"pavucontrol",

	// For eruption
	"rust",
	"protobuf-c",
	"gtksourceview4",
}

var archChaoticPackages = []string{
	// Packages for working graphical system with audio
	"swappy",
	"hyprpaper",
	"starship",
	"exa",
	"bat",
	"wofi",
	"xdg-desktop-portal-hyprland",
	"connman-gtk",
	"wlogout",
	"swaylock-effects",
	"wev",
	"dracula-icons-git",
	"dracula-cursors-git",
	"dracula-gtk-theme",
}

var aurPackages = []string{
	// Packages for working graphical system with audio
	"waybar-hyprland-no-systemd",
}

// Applications can be installed optionally (makes testing faster)
var applicationPackages = []string{
	"thunar",
	"mpv",
	"firefox",
	"evince",
	"eog",
	"godot",
	"glade",
	"texlive",
	"texlive-langgerman",
	"epiphany",
	"vscodium",
	"libreoffice-still",
	"libreoffice-still-de",
	"xmake",
	"biber",
	"pamac-aur",
	"mailspring",
	"teams",
	"anki",
}

func main() {
	// Parse Args
	var stage int = 1
	var allDefault bool
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-y" || arg == "--yes" {
				allDefault = true
			} else {
				stage = parseStage(os.Args[1])
			}
		}
	}

	if stage == 1 {
		logInfo("Performing Stage 1 ...")

		// Update the system clock
		exe("dinitctl start ntpd")

		logInfo("Refreshing new mirrorlist ...")
		// Install rankmirrors and create better mirrorlist
		exe("pacman -S --noconfirm --needed pacman-contrib parallel")
		rankmirrors("/etc/pacman.d/mirrorlist")

		// Install base system
		logInfo("Installing base packages ...")

		if isUEFI(allDefault) {
			basestrapPackages = append(basestrapPackages, "efibootmgr")
		}
		exe("basestrap /mnt " + strings.Join(basestrapPackages, " "))

		// fstabgen
		exeAppendFile("fstabgen -U /mnt", "/mnt/etc/fstab")

		// Copying repository into system
		exe("cp -r . /mnt/SamuraiOS")
		exe("chmod +x /mnt/SamuraiOS/scripts/install2.sh")

		// chroot into system
		logInfo("Stage 1 Done")
		logInfo("Now using chroot to go into /mnt ...")

		install2 := "artix-chroot /mnt /SamuraiOS/scripts/install2.sh"
		if allDefault {
			install2 += " -y"
		}

		exe(install2)
	} else if stage == 2 {
		logInfo("Performing Stage 2 ...")

		// set the time zone
		logInfo("Setting locale ...")
		region := promptWithDefault("Europe", allDefault, "Region")
		city := promptWithDefault("Vienna", allDefault, "City")

		exe(fmt.Sprint("ln -sf /usr/share/zoneinfo/", region, "/", city, " /etc/localtime"))

		exe("hwclock --systohc")

		// Set locale
		localesStr := promptWithDefault("de_AT.UTF-8, en_GB.UTF-8", allDefault, "Locale (comma seperated)")
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
		if isUEFI(allDefault) {
			exe("grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=grub")
		} else {
			exe("lsblk")
			prompt("What is your device? (e.g /dev/sda)")
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

		userName := promptWithDefault("ninja", allDefault, "Username")
		exe("useradd -m " + userName)

		logInfo("Enter password for ", userName)
		exe("passwd " + userName)

		hostName := promptWithDefault("samurai", allDefault, "Hostname")

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
		// Enable every user in the wheel group to use sudo
		exeArgs("go", "run", "scripts/replace.go", "/etc/sudoers", "# %wheel ALL=(ALL:ALL) ALL", "%wheel ALL=(ALL:ALL) ALL")
		// Show asteriks when typing sudo password
		exeArgs("go", "run", "scripts/replace.go", "/etc/sudoers", "# Defaults maxseq = 1000", "Defaults env_reset,pwfeedback")

		logInfo("Stage 2 Done")
		logInfo("Reboot into the new drive and execute \"sudo dinitctl enable connmand\" to activate the network daemon. After that reconnect to the internet and execute \"cd /SamuraiOS && go run install.go 3\"")
	} else if stage == 3 {
		logInfo("Performing Stage 3 ...")

		homeDir, _ := os.UserHomeDir()
		curDir, _ := os.Getwd()
		curUser, _ := user.Current()

		// Install arch repositories
		logInfo("Installing Arch repositories ...")
		sudoRankmirrors("etc/pacman.d/mirrorlist-arch")

		exe("sudo cp etc/pacman.d/mirrorlist-arch etc/pacman.d/mirrorlist-universe /etc/pacman.d/")
		exe("sudo cp etc/pacman.conf /etc/")
		// Install packages from arch repos and update repositories
		exe("sudo pacman-key --populate archlinux")

		// Install chaotic-aur
		logInfo("Installing chaotic-aur repository ...")
		exe("sudo pacman-key --recv-key 3056513887B78AEB --keyserver keyserver.ubuntu.com")
		exe("sudo pacman-key --lsign-key 3056513887B78AEB")
		exe("sudo pacman -U https://cdn-mirror.chaotic.cx/chaotic-aur/chaotic-keyring.pkg.tar.zst https://cdn-mirror.chaotic.cx/chaotic-aur/chaotic-mirrorlist.pkg.tar.zst")
		exeArgs("sudo", "go", "scripts/append.go", "echo [chaotic-aur]\nInclude = /etc/pacman.d/chaotic-mirrorlist", "/etc/pacman.conf")
		sudoRankmirrors("/etc/pacman.d/chaotic-mirrorlist")

		logInfo("Installing arch packages ...")
		exe("sudo pacman -Sy --noconfirm --needed artix-archlinux-support")
		exe("sudo pacman -S --noconfirm --needed " + strings.Join(archChaoticPackages, " "))

		// Install dinit-userservd
		if !dinitServiceExists("dinit-userservd") {
			logInfo("Installing dinit user service ...")
			exe("mkdir -p " + filepath.Join(homeDir, "/repos/dinit-userservd"))
			exe("git clone https://github.com/Xynonners/dinit-userservd.git " + filepath.Join(homeDir, "/repos/dinit-userservd"))
			os.Chdir(filepath.Join(homeDir, "/repos/dinit-userservd"))
			exe("makepkg -si --noconfirm")
			exe("sudo dinitctl enable dinit-userservd")
			os.Chdir(curDir)
			exeArgs("sudo", "go", "run", "scripts/append.go", "echo session optional pam_dinit_userservd.so", "/etc/pam.d/system-login")
			exe("rm -rf " + filepath.Join(homeDir, "/repos/dinit-userservd"))
		} else {
			logInfo("Skipping installation of dinit-userservd since it is already installed")
		}

		// Install eruptuion
		if !isInstalled("eruption") {
			logInfo("Installing eruption ...")
			exe("mkdir -p " + filepath.Join(homeDir, "repos/eruption"))
			exe("git clone --branch no-systemd https://github.com/PucklaJ/eruption.git " + filepath.Join(homeDir, "repos/eruption"))
			os.Chdir(filepath.Join(homeDir, "repos/eruption"))
			exe("make")
			exe("sudo make install")
			// Copying dinit services
			exe("mkdir -p " + filepath.Join(homeDir, ".config/dinit.d"))
			exe("cp support/dinit/eruption-audio-proxy " + filepath.Join(homeDir, ".config/dinit.d/"))
			exe("cp support/dinit/eruption-fx-proxy " + filepath.Join(homeDir, ".config/dinit.d/"))
			exe("cp support/dinit/eruption-process-monitor " + filepath.Join(homeDir, ".config/dinit.d/"))

			os.Chdir(curDir)
			exe("rm -rf " + filepath.Join(homeDir, "repos/eruption"))
			exe("sudo dinitctl enable eruption")
		} else {
			logInfo("Skipping Installation of eruption since it is already installed")
		}

		// Installing yay
		if !isInstalled("yay") {
			logInfo("Installing yay ...")
			exe("mkdir -p " + filepath.Join(homeDir, "/repos/yay"))
			exe("git clone https://aur.archlinux.org/yay.git " + filepath.Join(homeDir, "/repos/yay"))
			os.Chdir(filepath.Join(homeDir, "/repos/yay"))
			exe("makepkg -si --noconfirm")
			exe("rm -rf " + filepath.Join(homeDir, "/repos/yay"))
		} else {
			logInfo("Skipping installation of yay since it is already installed")
		}

		// Install yay packages
		logInfo("Installing AUR packages ...")
		exe("yay -S --noconfirm --needed " + strings.Join(aurPackages, " "))

		// Remove unneeded packages
		exeDontCare("sudo pacman -Rnsdd --noconfirm xdg-desktop-portal-gnome")
		exeDontCare("sudo pacman -Rnsdd --noconfirm xdg-desktop-portal-gtk")
		exeDontCare("sudo pacman -Rnsdd --noconfirm xdg-desktop-portal-kde")
		exeDontCare("sudo pacman -Rnsdd --noconfirm xdg-desktop-portal-wlr")

		logInfo("Clearing pacman cache ...")
		exe("sudo pacman -Scc --noconfirm")

		// Copy configuration files
		logInfo("Copying configuration files ...")

		// Copy non contents of repo
		repoEntries, err := os.ReadDir(curDir)
		if err != nil {
			logError(err)
			os.Exit(1)
		}

		var repoEntriesStr []string
		for _, e := range repoEntries {
			if e.IsDir() && !(e.Name() == "home" || e.Name() == "scripts") {
				repoEntriesStr = append(repoEntriesStr, e.Name())
			}
		}

		exe("sudo cp -r " + strings.Join(repoEntriesStr, " ") + " /")

		// Copy contents of home directory
		homeEntries, err := os.ReadDir("home/samurai")
		if err != nil {
			logError(err)
			os.Exit(1)
		}

		var homeEntriesStr []string
		for _, h := range homeEntries {
			homeEntriesStr = append(homeEntriesStr, filepath.Join("home/samurai", h.Name()))
		}

		exe("cp -r " + strings.Join(homeEntriesStr, " ") + " " + homeDir)
		exe("mkdir -p " + filepath.Join(homeDir, ".config/dinit.d/boot.d"))

		exe("go run scripts/replace.go " + filepath.Join(homeDir, "/.config/dinit.d/pipewire") + " samurai " + curUser.Username)
		exe("go run scripts/replace.go " + filepath.Join(homeDir, "/.config/dinit.d/pipewire-pulse") + " samurai " + curUser.Username)
		exe("go run scripts/replace.go " + filepath.Join(homeDir, "/.config/wlogout/style.css") + " samurai " + curUser.Username)
		exe("go run scripts/replace.go " + filepath.Join(homeDir, "/.config/qt5ct/qt5ct.conf") + " samurai " + curUser.Username)

		exe("sudo go run scripts/replace.go /etc/sddm.conf.d/default.conf samurai " + curUser.Username)

		exe("chmod +x " + filepath.Join(homeDir, ".config/gtk-3.0/import-gsettings"))

		// Copy wireplumber alsa configuration (Fix for broken headset audio)
		exe("sudo mkdir -p /etc/wireplumber/main.lua.d")
		exe("sudo cp /usr/share/wireplumber/main.lua.d/50-alsa-config.lua /etc/wireplumber/main.lua.d")
		exeArgs("sudo", "go", "run", "scripts/replace.go", "/etc/wireplumber/main.lua.d/50-alsa-config.lua", "--[\"api.alsa.headroom\"]      = 0", "[\"api.alsa.headroom\"]      = 1024")

		// Install go programs
		logInfo("Installing go programs ...")
		goDir := filepath.Join(homeDir, "go/src/github.com/PucklaJ")
		goPrograms, err := os.ReadDir(goDir)
		for _, gp := range goPrograms {
			os.Chdir(filepath.Join(goDir, gp.Name()))
			logInfo("Installing " + gp.Name() + " ...")
			exe("go install -buildvcs=false")
		}
		os.Chdir(curDir)

		// Change shell
		exe("chsh -s /usr/bin/fish")

		logInfo("Stage 3 Done")
		logInfo("Now logout and login again and execute \"cd /SamuraiOS && go run install.go 4\"")
	} else if stage == 4 {
		logInfo("Performing Stage 4 ...")

		homeDir, _ := os.UserHomeDir()

		// Activate dinit user services
		logInfo("Activating dinit user services ...")
		exe("mkdir -p " + filepath.Join(homeDir, "/.local/share/dinit"))
		exe("dinitctl enable pipewire")
		exe("dinitctl enable pipewire-pulse")
		exe("dinitctl enable eruption-audio-proxy")
		exe("dinitctl enable eruption-fx-proxy")

		logInfo("Installation Done")

		rmSamurai := promptWithDefaultYesNo(false, allDefault, "Remove /SamuraiOS")
		if rmSamurai {
			exe("sudo rm -rf /SamuraiOS")
		}

		launchHypr := promptWithDefaultYesNo(true, allDefault, "Execute \"sudo dinitctl enable sddm\" and launch into hyprland")
		if launchHypr {
			exe("sudo dinitctl enable sddm")
		}
	} else if stage == 5 {
		// Application Stage
		logInfo("Performing Stage 5 ...")

		exe("sudo pacman -S --noconfirm --needed " + strings.Join(applicationPackages, " "))
	} else if stage == 255 {
		// Testing
		logInfo("Performing Tests ...")

		exe("mkdir -p /tmp/stage_test")
		exe("git clone https://github.com/PucklaJ/SamuraiOS.git /tmp/stage_test")
		logError("Task not failed successfully")
		exe("rm -rf /tmp/stage_test")

		folderName := promptWithDefault("/tmp/another/samurai", allDefault, "Create which folder")
		exe("mkdir -p " + folderName)

		exeAppendFile("echo Hello World", "/tmp/hello_samurai")
		exe("go run scripts/replace.go /tmp/hello_samurai Hello Greetings")

		createTmpFile := promptWithDefaultYesNo(true, allDefault, "Create /tmp/install_test_os")
		if createTmpFile {
			exe("touch /tmp/install_test_os")
		}

		bak := backupName("/etc/pacman.d/mirrorlist")
		logInfo("Backup: ", bak)

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
		logError("\"", strings.Join(args, " "), "\" failed: ", err)
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

func promptWithDefault(defaultValue string, allDefault bool, msg ...any) string {
	msgStr := fmt.Sprint(msg...)
	msgStr = fmt.Sprint(msgStr, " (", defaultValue, ")")

	prompt(msgStr)
	if allDefault {
		fmt.Println(defaultValue)
		return defaultValue
	}

	value := inputWithDefault(defaultValue)
	return value
}

func promptWithDefaultYesNo(defaultValue, allDefault bool, msg ...any) bool {
	msgStr := fmt.Sprint(msg...)
	msgStr = fmt.Sprint(msgStr, " yes|no")

	var defaultValueStr string
	if defaultValue {
		defaultValueStr = "yes"
	} else {
		defaultValueStr = "no"
	}

	value := promptWithDefault(defaultValueStr, allDefault, msgStr)
	valueBool := strings.HasPrefix(strings.ToLower(value), "y")

	return valueBool
}

func isInstalled(program string) bool {
	_, err := exec.LookPath(program)
	return err == nil
}

func dinitServiceExists(service string) bool {
	cmd := exec.Command("sudo", "dinitctl", "status", service)
	err := cmd.Run()
	return err == nil
}

func parseStage(arg string) int {
	switch strings.ToLower(arg) {
	case "test":
		return 255
	case "apps":
		fallthrough
	case "application":
		fallthrough
	case "applications":
		return 5
	default:
		v, err := strconv.ParseUint(arg, 10, 64)
		if err != nil {
			logError("Failed to parse stage argument: ", err)
			os.Exit(1)
		}
		return int(v)
	}
}

func isUEFI(allDefault bool) bool {
	_, err := os.Stat("/sys/firmware/efi")
	if err != nil && !os.IsNotExist(err) {
		logInfo("Failed to check if the system is UEFI automatically: ", err, " Manual input required.")
		return promptWithDefaultYesNo(false, allDefault, "Is UEFI")
	}

	return err == nil
}

func backupName(filename string) string {
	for {
		filename = filename + ".bak"
		if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
			return filename
		}
	}
}

func sudoRankmirrors(mirrorlistPath string) {
	// Create back up
	mirrorlistBak := backupName(mirrorlistPath)
	exeArgs("sudo", "mv", mirrorlistPath, mirrorlistBak)
	// rank mirror list
	exeArgs("sudo", "go", "run", "scripts/append.go", "rankmirrors -n 5 -v -p "+mirrorlistBak, mirrorlistPath+".tmp")
	// Overwrite old mirrorlist
	exeArgs("sudo", "mv", mirrorlistPath+".tmp", mirrorlistPath)
}

func rankmirrors(mirrorlistPath string) {
	// Create back up
	mirrorlistBak := backupName(mirrorlistPath)
	exeArgs("mv", mirrorlistPath, mirrorlistBak)
	// rank mirror list
	exeAppendFile("rankmirrors -n 5 -v -p "+mirrorlistBak, mirrorlistPath+".tmp")
	// Overwrite old mirrorlist
	exeArgs("mv", mirrorlistPath+".tmp", mirrorlistPath)
}
