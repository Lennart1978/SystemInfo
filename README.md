# Lennart's System Information v1.5

This Go application provides information about your system, including memory usage, disk space usage, user, hostname, operating system details, CPU information, GPU information, and more. It also allows you to empty the trash on your system.

## Features

- Display system information in an easy-to-read format.
- Monitor memory and disk usage with progress bars.
- Empty the trash with the click of a button.
- Supports Linux systems.

## Requirements

- Go 1.16 or later
- Fyne library for GUI

## Installation

1. Clone this repository to your local machine.

```bash
git clone https://github.com/Lennart1978/SystemInfo

- cd SystemInfo

- go build main.go

## Usage

Launch the application, and it will display system information in a graphical user interface (GUI).
The "Memory Usage" and "Disk Space Usage" sections display the current memory and disk usage, respectively, with progress bars.
Click the "Empty Trash" button to empty the trash on your system.

## Known Issues

This application is designed for Linux systems and may not work on other operating systems.
