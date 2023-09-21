package main

import (
	"fmt"
	"log"
	"os"
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
	windowHeight    = 200
)

func main() {
	// Create a new app and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Lennart's System Information V1.0")
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
