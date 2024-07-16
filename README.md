# Hardware Monitor

This is a Go project that monitors and displays hardware status for a laptop. The project is OS agnostic and updates the values every 500ms. It highlights changes in hardware parameters by altering the background color of the updated values.

## Features

- Monitors CPU usage, memory usage, disk usage, network I/O, system uptime, and GPU status.
- Updates displayed values every 500ms.
- Highlights updated values with a contrasting background color for easy identification.
- OS agnostic, works on Windows, macOS, and Linux.

## Prerequisites

- Go 1.16 or higher
- NVIDIA drivers installed
- `nvidia-smi` available in the system PATH

## Installation

## Clone the repository:
```sh
   git clone https://github.com/rajasatyajit/hw-monitor.git
   cd hw-monitor
```
 
## Initialize the Go module
```sh
    go mod init hardware-monitor
```

## Install dependencies
```sh
    go get -u github.com/shirou/gopsutil/v4
```

## Usage

    Ensure nvidia-smi is available in your system PATH. nvidia-smi is typically installed with the NVIDIA drivers.

    Save the code in main.go.

    Run the project:

```sh
    go run main.go
```

## How It Works

    The program fetches hardware status using the gopsutil library and nvidia-smi for GPU metrics.
    It displays CPU usage, memory usage, disk usage, network I/O, system uptime, and GPU status.
    The display is updated every 500ms.
    Changes in hardware parameters are highlighted with a red background and white text.

## Customization

    You can customize the colors by modifying the ANSI escape codes in the const section.
    You can add more hardware parameters by extending the HardwareStatus struct and modifying the getHardwareStatus and display functions.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
## Acknowledgements

    This project uses the gopsutil library for hardware monitoring.
    This project uses the nvidia-smi tool for GPU monitoring.

This approach uses the `nvidia-smi` command-line tool to get GPU metrics on Windows. It parses the output and integrates it with the rest of the hardware monitoring information. This should avoid the issues related to Unix-specific headers and provide a working solution for GPU monitoring on Windows.