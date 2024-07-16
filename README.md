# Hardware Monitor

## Overview
The Hardware Monitor project is a Go application that monitors and displays hardware information from various devices in an OS-agnostic manner. It provides real-time updates on CPU usage, memory usage, disk usage, network I/O, and system uptime.

## Features
- Monitor CPU usage for each core
- Display memory usage statistics
- Track disk usage for different partitions
- Monitor network I/O traffic
- Display system uptime

## Dependencies
- [gopsutil](https://github.com/shirou/gopsutil/v4): Go package for hardware and OS statistics

## Usage
1. Ensure you have Go installed on your machine.
2. Install the required dependencies using `go get`.
3. Run the application using `go run main.go`.

## Contributing
Contributions are welcome! If you have any ideas for improvements or new features, feel free to open an issue or submit a pull request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.