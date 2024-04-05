package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var basePackages = []string{
	// Base packages to make the system work
	"base-devel",
	"neovim",
	"dhcpcd",
	"networkmanager",
	"fzf",
	"whois",

	// Packages for working graphical system with audio
	"fish",
	"pipewire",
	"pipewire-pulse",
	"pipewire-jack",
	"pipewire-alsa",
	"wireplumber",
	"sddm",
	"hyprland",
	"waybar",
	"polkit-gnome",
	"man",
	"mako",
	"kitty",
	"foot",
	"noto-fonts",
	"noto-fonts-emoji",
	"noto-fonts-cjk",
	"nerd-fonts",
	"playerctl",
	"qt5ct",
	"gnome-keyring",
	"grim",
	"wl-clipboard",
	"fcitx5",
	"fcitx5-qt",
	"fcitx5-gtk",
	"fcitx5-mozc",
	"fcitx5-configtool",
	"qt5-wayland",
	"qt6-wayland",
	"bluez",
	"bluez-utils",
	"gifski",
	"wf-recorder",
	"swayidle",
	"htop",
	"btop",
	"weston",
	"xdg-desktop-portal-hyprland",
	"qt5-graphicaleffects",
	"qt5-quickcontrols2",
	"cifs-utils",
	"webkit2gtk",
	"clang",
	"man-pages-de",
	"keychain",
	"gvfs-mtp",
	"gvfs-gphoto2",
	"sassc",
	"hyfetch",
	"flatpak",
	"eza",
	"bat",
	"starship",
	"wofi",
	"nm-connection-editor",
	"blueman",
	"ttf-fantasque-sans-mono",
	"hyprpaper",
	"swappy",
	"wev",
}

var aurPackages = []string{
	"samurai-select",
	"aur/dracula-gtk-theme",
	"aur/dracula-cursors-git",
	"backlight_control",
	"poweralertd",
	"wlogout",
	"swaylock-effects",
	"dracula-icons-git",
}

// Applications can be installed optionally (makes testing faster)
var applicationPackages = []string{
	"thunar",
	"xmake",
	"biber",
	"speech-dispatcher",
	"thunar-archive-plugin",
	"file-roller",
	"openconnect",
	"openrgb",
}

var flatpaks = []string{
	"com.github.tchx84.Flatseal",
	"org.gnome.Evince",
	"org.kde.KStyle.Kvantum",
	"org.godotengine.Godot",
	"org.gnome.eog",
	"org.libreoffice.LibreOffice",
	"net.ankiweb.Anki",
	"com.github.IsmaelMartinez.teams_for_linux",
	"io.mpv.Mpv",
	"org.pulseaudio.pavucontrol",
	"org.inkscape.Inkscape",
}

var vscodeExtensions = []string{
	"jeff-hykin.better-cpp-syntax",
	"xaver.clang-format",
	"dracula-theme.theme-dracula",
	"MS-CEINTL.vscode-language-pack-de",
	"MS-CEINTL.vscode-language-pack-ja",
	"golang.go",
	"vscode-icons-team.vscode-icons",
	"streetsidesoftware.code-spell-checker",
	"streetsidesoftware.code-spell-checker-german",
	"ms-python.python",
	"llvm-vs-code-extensions.vscode-clangd",
	"vadimcn.vscode-lldb",
	"ms-vscode.hexeditor",
	"prince781.vala",
	"jeanp413.open-remote-ssh",
	"wmaurer.change-case",
	"danielgavin.ols",
	"yzhang.markdown-all-in-one",
	"asvetliakov.vscode-neovim",
	"ngtystr.ppm-pgm-viewer-for-vscode",
}

var virtualizationPackages = []string{
	"virt-install",
	"libvirt",
	"qemu-desktop",
	"virt-manager",
	"dnsmasq",
	"distrobox",
	"podman",
}

func main() {
	// Parse Args
	var stage int = 1
	var allDefault, userDefault, addUserDirectly bool
	userDefault = true
	var argUserName, argPassword, argHostname string
	if len(os.Args) > 1 {
		args := os.Args[1:]
		for i := 0; i < len(args); i++ {
			arg := args[i]
			if arg == "-y" || arg == "--yes" {
				allDefault = true
			} else if arg == "-u" || arg == "--user" {
				userDefault = false
			} else if arg == "--name" {
				i++
				if i == len(args) {
					logError("Invalid arguments: username required after \"--name\"")
					os.Exit(1)
				}
				argUserName = args[i]
			} else if arg == "--pass" {
				i++
				if i == len(args) {
					logError("Invalid arguments: password required after \"--pass\"")
					os.Exit(1)
				}
				argPassword = args[i]
			} else if arg == "--add-user" {
				addUserDirectly = true
			} else if arg == "--host" {
				i++
				if i == len(args) {
					logError("Invalid arguments: hostname required after \"--host\"")
				}
				argHostname = args[i]
			} else {
				stage = parseStage(os.Args[1])
			}
		}
	}

	if stage == 1 {
		logInfo("Performing Stage 1 ...")

		if userExists("alarm") {
			exe("userdel -r alarm")
		}

		// Enable ParallelDownloads
		exeArgs("go", "run", "scripts/replace.go", "/etc/pacman.conf", "#ParallelDownloads = 5", "ParallelDownloads = 5\nILoveCandy")

		rankmirrors("/etc/pacman.d/mirrorlist")
		exe("pacman -Syu --noconfirm --needed " + strings.Join(basePackages, " "))

		// set the time zone
		if !fileExists("/etc/localtime") {
			logInfo("Setting locale ...")
			region := promptWithDefault("Europe", allDefault, "Region")
			city := promptWithDefault("Vienna", allDefault, "City")

			exe(fmt.Sprint("ln -sf /usr/share/zoneinfo/", region, "/", city, " /etc/localtime"))
		}

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

		keymap := promptWithDefault("de", allDefault, "Keymap")
		exeAppendFile("echo KEYMAP="+keymap, "/etc/vconsole.conf")

		// Boot Loader
		if !fileExists("/boot/extlinux/extlinux.conf") {
			writeExtlinuxConf("/", "/boot/extlinux/extlinux.conf")
		}

		// Set root password to root
		passwordPrompt("root", "root", false)

		if addUserDirectly {
			if !userExists(argUserName) {
				addUser(argUserName, argPassword, allDefault, userDefault)
			}
		}

		if !fileExists("/etc/hostname") {
			var hostName string
			if len(argHostname) != 0 {
				hostName = argHostname
			} else {
				hostName = promptWithDefault(argHostname, allDefault, "Hostname")
			}

			{
				hostNameFile, err := os.Create("/etc/hostname")
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				hostNameFile.WriteString(hostName)
				hostNameFile.Close()
			}
		}

		if !fileExists("/etc/hosts") {
			hostName := readFileTrim("/etc/hostname")
			hosts, err := os.Create("/etc/hosts")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			hosts.WriteString("127.0.0.1\tlocalhost\n")
			hosts.WriteString("::1\tlocalhost\n")
			hosts.WriteString("127.0.1.1\t" + hostName + ".localdomain\t" + hostName + "\n")
		}

		// Enable every user in the wheel group to use sudo
		exeArgs("go", "run", "scripts/replace.go", "/etc/sudoers", "# %wheel ALL=(ALL:ALL) ALL", "%wheel ALL=(ALL:ALL) ALL")
		// Show asteriks when typing sudo password
		exeArgs("go", "run", "scripts/replace.go", "/etc/sudoers", "# Defaults maxseq = 1000", "Defaults env_reset,pwfeedback")

		if !userExists("installer") {
			logInfo("Add user installer")
			exe("useradd -M -G wheel installer")
			pwCrypt, err := exeToStringSilent("mkpasswd installer")
			if err == nil {
				pwCrypt = strings.TrimSpace(pwCrypt)
				exeArgs("usermod", "--password", pwCrypt, "installer")
			}
		}

		curDir, _ := os.Getwd()

		exeDontCare("systemctl enable NetworkManager.service")
		exeDontCare("systemctl enable bluetooth.service")

		// Installing yay
		if !isInstalled("yay") {
			installAURPackage("yay-bin")
		} else {
			logInfo("Skipping installation of yay since it is already installed")
		}

		// Copy configuration files
		logInfo("Copying system configuration files ...")

		// Copy non contents of repo
		repoEntries, err := os.ReadDir(curDir)
		if err != nil {
			logError(err)
			os.Exit(1)
		}

		var repoEntriesStr []string
		for _, e := range repoEntries {
			if e.IsDir() && !(e.Name() == "home" || e.Name() == "scripts" || strings.HasPrefix(e.Name(), ".")) {
				repoEntriesStr = append(repoEntriesStr, e.Name())
			}
		}

		exe("cp -r " + strings.Join(repoEntriesStr, " ") + " /")

		// Copy wireplumber alsa configuration (Fix for broken headset audio)
		exe("mkdir -p /etc/wireplumber/main.lua.d")
		exe("cp /usr/share/wireplumber/main.lua.d/50-alsa-config.lua /etc/wireplumber/main.lua.d")
		exeArgs("go", "run", "scripts/replace.go", "/etc/wireplumber/main.lua.d/50-alsa-config.lua", "--[\"api.alsa.headroom\"]      = 0", "[\"api.alsa.headroom\"]      = 1024")

		logInfo("Clearing pacman cache ...")
		exe("pacman -Scc --noconfirm")

		rmSamurai := promptWithDefaultYesNo(false, allDefault, "Remove /SamuraiOS")
		if rmSamurai {
			exe("rm -rf /SamuraiOS")
		}

		exe("userdel installer")

		logInfo("Stage 1 Done")
		logInfo("Now reboot into the system and do the following\n\t1. Login as root with pasword root\n\t2. Partion and mount the home partition\n\t3. Add users\n\t4. Execute `systemctl enable --now sddm.service`")
	} else if stage == 2 {
		// Application Stage
		logInfo("Performing Stage 2 ...")
		homeDir, _ := os.UserHomeDir()

		allPackages := append(aurPackages, applicationPackages...)
		exe("yay -S --noconfirm --needed " + strings.Join(allPackages, " "))

		installOdinfmt()
		installGoPrograms()

		// Disable eruption for now
		// exeDontCare("systemctl enable --user eruption-audio-proxy.service")
		// exeDontCare("systemctl enable --user eruption-fx-proxy.service")
		// exeDontCare("systemctl enable --user eruption-process-monitor.service")
		// exeDontCare("sudo systemctl enable --now eruption.service")

		for _, ext := range vscodeExtensions {
			exe("codium --install-extension " + ext)
		}

		// Install odin formatter for neoformat
		logInfo("Installing neovim formatter for odin ...")
		if ok := promptWithDefaultYesNo(true, allDefault, "You need to start neovim at least once before continuing. Did you do this?"); ok {
			exe("cp home/ninja/.local/share/nvim/site/pack/pckr/opt/neoformat/autoload/neoformat/formatters/odin.vim " + filepath.Join(homeDir, ".local/share/nvim/site/pack/pckr/opt/neoformat/autoload/neoformat/formatters/"))
		}

		logInfo("Installing flatpaks ...")
		exe("flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo")
		exe("flatpak install --assumeyes " + strings.Join(flatpaks, " "))

		logInfo("Stage 2 Done")
	} else if stage == 4 {
		// Virtualizazion Stage
		logInfo("Performing Virtualization Stage ...")

		exe("sudo pacman -S --noconfirm --needed " + strings.Join(virtualizationPackages, " "))

		exe("sudo systemctl enable --now libvirtd.service")
		exe("sudo systemctl enable --now docker.service")
		exe("sudo virsh net-start default")
		exe("sudo virsh net-autostart default")

		curUser, err := user.Current()
		if err != nil {
			logError("Failed to get user: ", err)
			os.Exit(1)
		}

		exe("sudo usermod -aG libvirt " + curUser.Username)
		exe("sudo usermod -aG libvirt-qemu " + curUser.Username)
		exe("sudo usermod -aG kvm " + curUser.Username)
		exe("sudo usermod -aG input " + curUser.Username)
		exe("sudo usermod -aG disk " + curUser.Username)

		logInfo("Stage 7 Done")
		logInfo("Now reboot and everything should be set up")
	} else if stage == 5 {
		logInfo("Performing User Stage ...")

		// User Stage to add another user
		addUser(argUserName, argPassword, allDefault, userDefault)

		logInfo("User Stage Done")
	} else if stage == 255 {
		// Testing
		logInfo("Performing Tests ...")

		scriptDir, _ := os.Getwd()
		compileStyles(filepath.Join(scriptDir, "/home/ninja"))

		writeExtlinuxConf("/", "/tmp/extlinux.conf")
		file, err := os.Open("/tmp/extlinux.conf")
		if err != nil {
			logError(err)
			os.Exit(1)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			logError(data)
			os.Exit(1)
		}
		fmt.Println(string(data))

		pkgName := searchPkgName("/tmp/yay-bin")
		fmt.Println("\nPackagename:", pkgName)

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

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0744)
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

func exeToStringSilent(command string) (string, error) {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		return "", errors.New("No Command")
	}

	var args []string
	if len(words) > 1 {
		args = words[1:]
	}

	var builder strings.Builder
	var stderr strings.Builder

	cmd := exec.Command(words[0], args...)
	cmd.Stdout = &builder
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("\"%s\" failed: %s", command, stderr.String())
	}

	return builder.String(), nil
}

func exeRetry(command string) {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		logError("No Command")
		os.Exit(1)
	}

	var args []string
	if len(words) > 1 {
		args = words[1:]
	}

	logScript(command)

	var tries int

	for {
		tries++

		cmd := exec.Command(words[0], args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			logError("failed: ", err, ". Trying again. ", 5-tries, " tries left ...")
			if tries == 5 {
				logError("failed 5 times quitting")
				os.Exit(1)
			}
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}
}

func exeEnv(command string, envNameValue ...string) {
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
	cmd.Env = envNameValue

	logScript(
		fmt.Sprint(
			strings.Join(envNameValue, " "),
			" ",
			command,
		),
	)

	if err := cmd.Run(); err != nil {
		logError("\"", command, "\" failed: ", err)
		os.Exit(1)
	}
}

func copyConfig(src string) {
	var dst string

	homeDir, _ := os.UserHomeDir()

	if strings.HasPrefix(src, "home/ninja") {
		dst = strings.Replace(src, "home/ninja", homeDir, 1)
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
	case "applications", "apps", "application":
		return 2
	case "gaming":
		return 3
	case "virt", "virtualization":
		return 4
	case "user":
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
	exeArgs("sudo", "go", "run", "scripts/append.go", "rankmirrors -n 0 -v -p "+mirrorlistBak, mirrorlistPath+".tmp")
	// Overwrite old mirrorlist
	exeArgs("sudo", "mv", mirrorlistPath+".tmp", mirrorlistPath)
}

func rankmirrors(mirrorlistPath string) {
	// Create back up
	mirrorlistBak := backupName(mirrorlistPath)
	exeArgs("mv", mirrorlistPath, mirrorlistBak)
	// rank mirror list
	exeAppendFile("rankmirrors -n 0 -v -p "+mirrorlistBak, mirrorlistPath+".tmp")
	// Overwrite old mirrorlist
	exeArgs("mv", mirrorlistPath+".tmp", mirrorlistPath)
}

func mountedDeviceName(path string, fullPartition bool) string {
	dfOut, err := exeToStringSilent("df")
	if err == nil {
		lines := strings.Split(dfOut, "\n")
		for _, line := range lines {
			words := strings.Split(line, " ")

			// Clear empty ones
			for i := 0; i < len(words); i++ {
				if words[i] == "" {
					words = append(words[:i], words[i+1:]...)
					i--
				}
			}

			if len(words) < 6 {
				continue
			}

			partition := words[0]
			directory := words[5]

			if directory == path {
				if fullPartition {
					return partition
				}
				return strings.Trim(partition, "0123456789")
			}
		}
	}

	logInfo("Failed to get mounted device name automatically. Manual input required.")
	exe("lsblk")
	for {
		prompt("Which device is currently mounted at " + path + " (e.g. /dev/sda)?")
		if device := input(); device != "" {
			return device
		}
	}
}

func passwordPrompt(username, argPassword string, allDefault bool) {
	var pw string
	if len(argPassword) != 0 {
		pw = argPassword
	} else {
		for {
			pw1 := promptWithDefault("s", allDefault, "Password")
			pw2 := promptWithDefault("s", allDefault, "Reenter Password")

			if pw1 != pw2 {
				logError("Passwords do not match. Please enter again!")
			} else {
				pw = pw1
				break
			}
		}
	}

	pwCrypt, err := exeToStringSilent("mkpasswd " + pw)
	if err == nil {
		pwCrypt = strings.TrimSpace(pwCrypt)
		// Set the root password if this is the first user
		entries, err := os.ReadDir("/home")
		if username != "root" && err == nil && len(entries) == 0 {
			exeArgs("usermod", "--password", pwCrypt, username)
		}
		exeArgs("usermod", "--password", pwCrypt, username)
	} else {
		logInfo("Failed to create password using mkpasswd. passwd is required.")
		logInfo("Root password")
		exe("passwd")
		logInfo("Password for user " + username)
		exe("passwd " + username)
	}
}

func chaoticInstalled() bool {
	cmd := exec.Command("pacman", "-Qk", "chaotic-mirrorlist")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func installOdinfmt() {
	if isInstalled("odinfmt") {
		logInfo("Skipping installation of odinfmt since it is already installed")
		return
	}

	logInfo("Installing odinfmt ...")
	curDir, _ := os.Getwd()
	tmpDir := os.TempDir()
	olsDir := filepath.Join(tmpDir, "ols")

	exeDontCare("rm -rf " + olsDir)
	exe("git clone https://github.com/DanielGavin/ols.git -b master --depth 1 " + olsDir)
	os.Chdir(olsDir)
	exe("./odinfmt.sh")
	exe("sudo cp odinfmt /usr/bin")
	os.Chdir(curDir)
}

func addUser(username, password string, allDefault, userDefault bool) {
	var userName string
	if len(username) != 0 {
		userName = username
	} else {
		userName = promptWithDefault("ninja", allDefault && userDefault, "Username")
	}
	exe("useradd -m " + userName)
	passwordPrompt(userName, password, allDefault && userDefault)

	exe("usermod -aG wheel " + userName)

	homeDir := "/home/" + userName

	// Copy contents of home directory
	homeEntries, err := os.ReadDir("home/ninja")
	if err != nil {
		logError(err)
		os.Exit(1)
	}

	var homeEntriesStr []string
	for _, h := range homeEntries {
		homeEntriesStr = append(homeEntriesStr, filepath.Join("home/ninja", h.Name()))
	}

	exe("cp -r " + strings.Join(homeEntriesStr, " ") + " " + homeDir)
	// Delete neoformat folder because pckr won't install neoformat because it believes that is already installed
	exe("rm -r " + filepath.Join(homeDir, ".local/share/nvim/site/pack/pckr/opt/neoformat"))

	exe("go run scripts/replace.go " + filepath.Join(homeDir, "/.config/qt5ct/qt5ct.conf") + " ninja " + userName)

	exe("chsh -s /usr/bin/fish " + userName)

	compileStyles(homeDir)

	exeArgs("chown", "-R", userName, homeDir)
}

func installGoPrograms() {
	homeDir, _ := os.UserHomeDir()
	curDir, _ := os.Getwd()
	// Install go programs
	logInfo("Installing go programs ...")
	goDir := filepath.Join(homeDir, "go/src/samurai")
	goPrograms, err := os.ReadDir(goDir)
	if err != nil {
		logError(err)
		os.Exit(1)
	}
	for _, gp := range goPrograms {
		if isInstalled(gp.Name()) {
			continue
		}

		os.Chdir(filepath.Join(goDir, gp.Name()))
		logInfo("Installing " + gp.Name() + " ...")
		exe("go install -buildvcs=false")
	}
	os.Chdir(curDir)
}

func searchPkgName(dirName string) string {
	entries, err := os.ReadDir(dirName)
	if err != nil {
		return ""
	}

	for _, e := range entries {
		if strings.Contains(e.Name(), ".pkg.tar") {
			return filepath.Join(dirName, e.Name())
		}
	}

	return ""
}

func installAURPackage(packageName string) {
	logInfo("Installing " + packageName + " ...")
	tmpDir := os.TempDir()
	curDir, _ := os.Getwd()
	yayDir := filepath.Join(tmpDir, packageName)
	exeDontCare("rm -rf " + yayDir)
	exe("mkdir -p " + yayDir)
	exe("git clone https://aur.archlinux.org/" + packageName + ".git --depth 1 " + yayDir)
	os.Chdir(yayDir)
	exe("chown -R installer " + yayDir)
	exe("sudo -u installer makepkg -s --noconfirm")
	pkgName := searchPkgName(yayDir)
	exe("pacman -U --needed --noconfirm " + pkgName)
	os.Chdir(curDir)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func userExists(userName string) bool {
	cmd := exec.Command("id", userName)
	err := cmd.Run()
	return err == nil
}

func readFileTrim(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		logError(err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logError(err)
		os.Exit(1)
	}

	return strings.TrimSpace(string(data))
}

func compileStyles(homeDir string) {
	configDir := filepath.Join(homeDir, ".config")
	compileScript := filepath.Join(configDir, "compile_style.sh")
	setupDarkScript := filepath.Join(configDir, "setup_dark_theme.sh")

	exe(compileScript)
	exeEnv(setupDarkScript, "DONT_RESTART=1", "DONT_MODIFY_HOME=1")
}

const extlinuxConfTemplate = `DEFAULT arch
MENU TITLE Boot Menu
PROMPT 0
TIMEOUT 50

LABEL arch
MENU LABEL Arch Linux ARM
LINUX /Image
INITRD /initramfs-linux.img
FDT /dtbs/rockchip/rk3399-pinebook-pro.dtb
APPEND root=UUID=%[1]s rw

LABEL arch-fallback
MENU LABEL Arch Linux ARM with fallback initramfs
LINUX /Image
INITRD /initramfs-linux-fallback.img
FDT /dtbs/rockchip/rk3399-pinebook-pro.dtb
APPEND root=UUID=%[1]s rw`

func writeExtlinuxConf(mountPath, filePath string) {
	logInfo("Generating extlinux config ...")
	rootDevice := mountedDeviceName(mountPath, true)

	blkid := exeToString("blkid " + rootDevice)
	uuidStart := strings.IndexRune(blkid, '"')
	uuidAndSuffix := blkid[uuidStart+1:]
	uuidEnd := strings.IndexRune(uuidAndSuffix, '"')
	uuid := uuidAndSuffix[:uuidEnd]

	dir := filepath.Dir(filePath)
	exe("mkdir -p " + dir)
	file, err := os.Create(filePath)
	if err != nil {
		logError(err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Fprintf(file, extlinuxConfTemplate, uuid)
}
