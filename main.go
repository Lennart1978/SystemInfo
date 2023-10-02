package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var trashDialog *dialog.CustomDialog

const (
	// Refresh interval
	refreshInterval = time.Second
	windowWidth     = 350
	windowHeight    = 600
)

func getLinuxVersion() (string, error) {
	cmd := exec.Command("lsb_release", "-d", "-s")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func getLinuxKernelVersion() (string, error) {
	cmd := exec.Command("uname", "-r")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func getDesktopEnvironment() (string, error) {
	desktopEnv := os.Getenv("DESKTOP_SESSION")
	if desktopEnv != "" {
		return desktopEnv, nil
	}

	return "", fmt.Errorf("desktop-Umgebung nicht erkannt")
}

func getCPUModel() (string, error) {
	cpuinfo, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(cpuinfo), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "model name") {
			fields := strings.Fields(line)
			model := strings.Join(fields[3:], " ")
			return model, nil
		}
	}

	return "", fmt.Errorf("cpu-Modell nicht gefunden")
}

func getCPUSpeed() (string, error) {
	cpuinfo, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(cpuinfo), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu MHz") {
			fields := strings.Fields(line)
			speed := fields[len(fields)-1]
			return speed + " MHz", nil
		}
	}

	return "", fmt.Errorf("cpu-Geschwindigkeit nicht gefunden")
}

func getTemperature() (string, error) {
	// Pfad zum Temperatursensor (kann je nach System variieren)
	tempFile := "/sys/class/thermal/thermal_zone0/temp"

	temperature, err := os.ReadFile(tempFile)
	if err != nil {
		return "", err
	}

	// Konvertieren Sie die Temperatur in Grad Celsius
	tempValue := strings.TrimSpace(string(temperature))
	tempInt := tempValue[:len(tempValue)-3]
	tempFloat := tempInt[:len(tempInt)-2] + tempInt[len(tempInt)-2:]

	return tempFloat + " °C", nil
}

func getGPUModel() (string, error) {
	cmd := exec.Command("lspci", "-v")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	var gpuModel string

	for _, line := range lines {
		if strings.Contains(line, "VGA compatible controller") {
			gpuModel = strings.TrimSpace(line)
			break
		}
	}
	_, gpuModel, _ = strings.Cut(gpuModel, "VGA compatible controller: ")
	return gpuModel, nil
}

func getRAMSize() (string, error) {
	cmd := exec.Command("free", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("ungültige Ausgabe")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return "", fmt.Errorf("ungültige Ausgabe")
	}

	totalRAM := fields[1]
	return totalRAM + " MB", nil
}

func main() {
	// Create a new app and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Lennart's System Information V1.5")
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))
	myWindow.SetFixedSize(true)

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("error retrieving memory information: %v", err)
	}

	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Fatalf("error retrieving disk information: %v", err)
	}

	// Create progress bars
	memProgressBar := widget.NewProgressBar()
	memProgressBar.Max = 100
	memProgressBar.SetValue(100 - memInfo.UsedPercent)

	diskProgressBar := widget.NewProgressBar()
	diskProgressBar.Max = 100
	diskProgressBar.SetValue(diskInfo.UsedPercent)

	// Create labels
	memLabel := widget.NewLabel(fmt.Sprintf("Memory Usage: %.2f%%", 100-memInfo.UsedPercent))
	memLabel.Alignment = fyne.TextAlignCenter
	diskLabel := widget.NewLabel(fmt.Sprintf("Disk Space Usage: %.2f%%", diskInfo.UsedPercent))
	diskLabel.Alignment = fyne.TextAlignCenter

	// added after V1.0
	// Username:
	userName := widget.NewLabel(fmt.Sprintf("User: %s", os.Getenv("USER")))

	// Hostname:
	hName, err := os.Hostname()
	if err != nil {
		log.Fatalf("error retrieving hostname: %v", err)
		hName = "unknown"
	}
	hostname := widget.NewLabel(fmt.Sprintf("Hostname: %s", hName))

	// OS name:
	osName := widget.NewLabel(fmt.Sprintf("OS: %s", runtime.GOOS))

	// Arch:
	osArch := widget.NewLabel(fmt.Sprintf("Arch: %s", runtime.GOARCH))

	// OS Version:
	osV, err := getLinuxVersion()
	osV = strings.TrimSpace(osV)
	if err != nil {
		log.Fatalf("error retrieving linux version: %v", err)
		osV = "unknown"
	}
	osVersion := widget.NewLabel(fmt.Sprintf("Version: %s", osV))

	// Kernel Version:
	kernelV, err := getLinuxKernelVersion()
	kernelV = strings.TrimSpace(kernelV)
	if err != nil {
		log.Fatalf("error retrieving linux kernel version: %v", err)
		kernelV = "unknown"
	}
	kernelVersion := widget.NewLabel(fmt.Sprintf("Kernel Version: %s", kernelV))

	// Desktop Environment:
	desktopEnv, err := getDesktopEnvironment()
	if err != nil {
		log.Fatalf("error retrieving desktop environment: %v", err)
		desktopEnv = "unknown"
	}
	desktopEnvironment := widget.NewLabel(fmt.Sprintf("Desktop Environment: %s", desktopEnv))

	// CPU Cores:
	cpuCores := widget.NewLabel(fmt.Sprintf("CPU Cores: %d", runtime.NumCPU()))

	// CPU Model:
	cpuM, err := getCPUModel()
	if err != nil {
		log.Fatalf("error retrieving cpu model: %v", err)
		cpuM = "unknown"
	}
	cpuModel := widget.NewLabel(fmt.Sprintf("CPU Model: %s", cpuM))

	// CPU Speed:
	cpuS, err := getCPUSpeed()
	if err != nil {
		log.Fatalf("error retrieving cpu speed: %v", err)
		cpuS = "unknown"
	}
	cpuSpeed := widget.NewLabel(fmt.Sprintf("CPU Speed: %s", cpuS))

	// CPU Temperature:
	cpuTemp, err := getTemperature()
	if err != nil {
		log.Fatalf("error retrieving cpu temperature: %v", err)
		cpuTemp = "unknown"
	}
	cpuTemperature := widget.NewLabel(fmt.Sprintf("CPU Temperature: %s", cpuTemp))

	// GPU Name:
	gpuN, err := getGPUModel()
	if err != nil {
		log.Fatalf("error retrieving gpu model: %v", err)
		gpuN = "unknown"
	}
	gpuName := widget.NewLabel(fmt.Sprintf("GPU: %s", gpuN))

	// Memory:
	m, err := getRAMSize()
	if err != nil {
		log.Fatalf("error retrieving memory size: %v", err)
		m = "unknown"
	}
	memory := widget.NewLabel(fmt.Sprintf("Memory: %s", m))

	// Start monitoring the progress bars in a seperate Goroutine
	go startMonitoring(myWindow, memLabel, memProgressBar, diskLabel, diskProgressBar)

	// Create a button to empty the trash
	emptyTrashButton := widget.NewButton("Empty Trash", func() {
		content := container.NewVBox(
			widget.NewLabel("Are you sure you want to empty the trash?"),
			container.NewHBox(
				widget.NewButton("Yes", func() {
					err := emptyTrash()
					if err != nil {
						dialog.NewError(err, myWindow)
					} else {
						// Schließe das Dialogfeld nach erfolgreicher Löschung
						trashDialog.Hide()
					}
				}),
			),
		)
		trashDialog = dialog.NewCustom("Confirm Emptying Trash", "Cancel", content, myWindow)
		trashDialog.Show()
	})

	// Create a container for the widgets
	content := container.NewVBox(
		memLabel,
		memProgressBar,
		diskLabel,
		diskProgressBar,
		emptyTrashButton,
		userName,
		hostname,
		osName,
		osArch,
		osVersion,
		kernelVersion,
		desktopEnvironment,
		cpuCores,
		cpuModel,
		cpuSpeed,
		cpuTemperature,
		gpuName,
		memory,
	)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func emptyTrash() error {
	// Channel zur Fehlerkommunikation erstellen
	errCh := make(chan error, 1)
	// Empty the trash in einem separaten Goroutine ausführen
	fmt.Println("Emptying trash...")
	go func() {
		defer close(errCh)
		err := run()
		if err != nil {
			errCh <- err // Fehler in den Channel senden
		}
	}()
	// Warten auf das Ende der Goroutine
	<-errCh

	fmt.Println("Trash emptied !")
	return nil
}

func run() error {
	trashPath := os.ExpandEnv("$HOME/.local/share/Trash/files/")
	err := os.RemoveAll(trashPath)
	if err != nil {
		return fmt.Errorf("error emptying trash: %v", err)
	}
	return nil
}

func startMonitoring(window fyne.Window, memLabel *widget.Label, memProgressBar *widget.ProgressBar, diskLabel *widget.Label, diskProgressBar *widget.ProgressBar) {
	for {
		// Aktualisiere die Systeminformationen
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Fehler beim Abrufen von Memory-Informationen: %v", err)
		}

		diskInfo, err := disk.Usage("/")
		if err != nil {
			log.Printf("Fehler beim Abrufen von Disk-Informationen: %v", err)
		}

		// Aktualisiere die Anzeige mit den neuen Informationen
		memLabel.SetText(fmt.Sprintf("Memory Usage: %.2f%%", 100-memInfo.UsedPercent))
		memProgressBar.SetValue(100 - memInfo.UsedPercent)
		diskLabel.SetText(fmt.Sprintf("Disk Space Usage: %.2f%%", diskInfo.UsedPercent))
		diskProgressBar.SetValue(diskInfo.UsedPercent)

		// Warte für das nächste Aktualisierungsintervall
		time.Sleep(refreshInterval)
	}
}
