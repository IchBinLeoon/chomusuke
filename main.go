package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

const version = "1.0.0"

const (
	BOLDBLACK = "\033[1;30m"
	BOLDRED = "\033[1;31m"
	BOLDGREEN = "\033[1;32m"
	BOLDYELLOW = "\033[1;33m"
	BOLDBLUE = "\033[1;34m"
	BOLDMAGENTA = "\033[1;35m"
	BOLDCYAN = "\033[1;36m"
	BOLDWHITE = "\033[1;37m"
	BLACK = "\033[0;30m"
	RED = "\033[0;31m"
	GREEN = "\033[0;32m"
	YELLOW = "\033[0;33m"
	BLUE = "\033[0;34m"
	MAGENTA = "\033[0;35m"
	CYAN = "\033[0;36m"
	WHITE = "\033[0;37m"
	RESET = "\033[0;m"
)

const ASCII = `
 	  ,_     _		
 	  |\\_,-~/		
 	  / _  _ |    ,--.	
	 (  `+ YELLOW + `@  @` + RESET + ` )   / ,-'	
	  \  _` + RED + `T` + RESET + `_/-._( (		
 	 /         '. \		
	|         _  \ |	
	 \ \ ,  /      |	
 	  || |-_\__   /		
	 ((_/'(____,-'		
`

func main() {
	checkSystem()
	ascii := sliceASCII()
	host := fmt.Sprintf("\n\t\t\t\t%s%s%s@%s%s%s", BOLDGREEN, getUsername(), RESET, BOLDGREEN, getHostname(), RESET)
	fmt.Println(host)
	line := fmt.Sprintf("%s%s", ascii[1], strings.Repeat("-", len(fmt.Sprintf("%s@%s", getUsername(), getHostname()))))
	fmt.Println(line)
	fmt.Println(fmt.Sprintf("%s%sOS%s\t %s", ascii[2], BOLDCYAN, RESET, getOS()))
	fmt.Println(fmt.Sprintf("%s%sKernel%s\t %s", ascii[3], BOLDCYAN, RESET, getKernel()))
	fmt.Println(fmt.Sprintf("%s%sArch%s\t %s", ascii[4], BOLDCYAN, RESET, getArchitecture()))
	fmt.Println(fmt.Sprintf("%s%sShell%s\t %s", ascii[5], BOLDCYAN, RESET, getShell()))
	fmt.Println(fmt.Sprintf("%s%sDE%s\t %s", ascii[6], BOLDCYAN, RESET, getDE()))
	fmt.Println(fmt.Sprintf("%s%sUptime%s\t %s", ascii[7], BOLDCYAN, RESET, getUptime()))
	fmt.Println(fmt.Sprintf("%s%sCPU%s\t %s", ascii[8], BOLDCYAN, RESET, getCPU()))
	fmt.Println(fmt.Sprintf("%s%sMemory%s\t %s", ascii[9], BOLDCYAN, RESET, getMemory()))
	fmt.Println(ascii[10])
	palettes := getPalettes()
	fmt.Println(fmt.Sprintf("\t\t\t\t%s", palettes[0]))
	fmt.Println(fmt.Sprintf("\t\t\t\t%s\n", palettes[1]))
}

func checkSystem() {
	if runtime.GOOS != "linux" {
		fmt.Println("Unsupported Operating System")
		os.Exit(0)
	}
}

func sliceASCII() []string {
	ascii := strings.Split(ASCII, "\n")
	return ascii
}

func readFile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func getUsername() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
	}
	username := currentUser.Username
	return username
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

func getOS() string {
	distro := ""
	data := readFile("/etc/os-release")
	re := regexp.MustCompile(`PRETTY_NAME="(.*)"`)
	m := re.FindStringSubmatch(data)
	if m != nil {
		distro = m[1]
	}
	return distro
}

func getKernel() string {
	kernel := readFile("/proc/sys/kernel/osrelease")
	return kernel
}

func getArchitecture() string {
	utsname := syscall.Utsname{}
	if err := syscall.Uname(&utsname); err != nil {
		return ""
	}
	data := make([]byte, 0)
	for _, i := range utsname.Machine {
		if i == 0 {
			break
		}
		data = append(data, byte(i))
	}
	return string(data)
}

func getShell() string {
	shell := path.Base(os.Getenv("SHELL"))
	return shell
}

func getUptime() string {
	data := readFile("/proc/uptime")
	dataSlice := strings.Split(data, " ")
	t, err := strconv.ParseFloat(dataSlice[0], 64)
	if err != nil {
		return ""
	}
	h := math.Floor(t / 3600)
	m := math.Floor((t - h*3600) / 60)
	s := t - (h*3600 + m*60)
	uptime := fmt.Sprintf("%0.fh %0.fm %0.fs", h, m, s)
	return uptime
}

func getDE() string {
	if os.Getenv("XDG_CURRENT_DESKTOP") != "" {
		return os.Getenv("XDG_CURRENT_DESKTOP")
	} else if os.Getenv("DESKTOP_SESSION") != "" {
		return os.Getenv("DESKTOP_SESSION")
	} else {
		return ""
	}
}

func getCPU() string {
	cpu := ""
	data := readFile("/proc/cpuinfo")
	re := regexp.MustCompile("model name\t: (.*)")
	m := re.FindStringSubmatch(data)
	if m != nil {
		cpu = m[1]
	}
	return cpu
}

func getMemory() string {
	totalStr := ""
	availableStr := ""
	data := readFile("/proc/meminfo")
	reTotal := regexp.MustCompile("MemTotal:(.*) kB")
	mTotal := reTotal.FindStringSubmatch(data)
	if mTotal != nil {
		totalStr = strings.ReplaceAll(mTotal[1], " ", "")
	}
	reAvailable := regexp.MustCompile("MemAvailable:(.*) kB")
	mAvailable := reAvailable.FindStringSubmatch(data)
	if mAvailable != nil {
		availableStr = strings.ReplaceAll(mAvailable[1], " ", "")
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		total = 0
	}
	available, err := strconv.Atoi(availableStr)
	if err != nil {
		available = 0
	}
	memory := fmt.Sprintf("%dMB / %dMB", (total - available) / 1024, total / 1024)
	return memory
}

func getPalettes() [2]string {
	palettes := [2]string{}
	palettes[0] = fmt.Sprintf("%s● %s● %s● %s● %s● %s● %s● %s●", BLACK, RED, GREEN, YELLOW, BLUE, MAGENTA, CYAN, WHITE)
	palettes[1] = fmt.Sprintf("%s● %s● %s● %s● %s● %s● %s● %s●", BOLDBLACK, BOLDRED, BOLDGREEN, BOLDYELLOW, BOLDBLUE, BOLDMAGENTA, BOLDCYAN, BOLDWHITE)
	return palettes
}
