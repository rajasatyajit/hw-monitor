package nvml

/*
#cgo LDFLAGS: -lnvidia-ml
#include <nvml.h>
#include <windows.h>

nvmlReturn_t nvmlInitWrapper() {
    return nvmlInit();
}

nvmlReturn_t nvmlShutdownWrapper() {
    return nvmlShutdown();
}

nvmlReturn_t nvmlDeviceGetCountWrapper(unsigned int *deviceCount) {
    return nvmlDeviceGetCount(deviceCount);
}

nvmlReturn_t nvmlDeviceGetHandleByIndexWrapper(unsigned int index, nvmlDevice_t *device) {
    return nvmlDeviceGetHandleByIndex(index, device);
}

nvmlReturn_t nvmlDeviceGetNameWrapper(nvmlDevice_t device, char *name, unsigned int length) {
    return nvmlDeviceGetName(device, name, length);
}

nvmlReturn_t nvmlDeviceGetUtilizationRatesWrapper(nvmlDevice_t device, nvmlUtilization_t *utilization) {
    return nvmlDeviceGetUtilizationRates(device, utilization);
}

nvmlReturn_t nvmlDeviceGetMemoryInfoWrapper(nvmlDevice_t device, nvmlMemory_t *memory) {
    return nvmlDeviceGetMemoryInfo(device, memory);
}
*/
import "C"
import (
	"fmt"
)

type Device struct {
	handle C.nvmlDevice_t
}

func Init() error {
	if res := C.nvmlInitWrapper(); res != C.NVML_SUCCESS {
		return fmt.Errorf("failed to initialize NVML: %d", res)
	}
	return nil
}

func Shutdown() error {
	if res := C.nvmlShutdownWrapper(); res != C.NVML_SUCCESS {
		return fmt.Errorf("failed to shutdown NVML: %d", res)
	}
	return nil
}

func DeviceCount() (int, error) {
	var count C.uint
	if res := C.nvmlDeviceGetCountWrapper(&count); res != C.NVML_SUCCESS {
		return 0, fmt.Errorf("failed to get device count: %d", res)
	}
	return int(count), nil
}

func DeviceHandleByIndex(index int) (*Device, error) {
	var device C.nvmlDevice_t
	if res := C.nvmlDeviceGetHandleByIndexWrapper(C.uint(index), &device); res != C.NVML_SUCCESS {
		return nil, fmt.Errorf("failed to get device handle by index: %d", res)
	}
	return &Device{handle: device}, nil
}

func (d *Device) Name() (string, error) {
	var name [C.NVML_DEVICE_NAME_BUFFER_SIZE]C.char
	if res := C.nvmlDeviceGetNameWrapper(d.handle, &name[0], C.NVML_DEVICE_NAME_BUFFER_SIZE); res != C.NVML_SUCCESS {
		return "", fmt.Errorf("failed to get device name: %d", res)
	}
	return C.GoString(&name[0]), nil
}

func (d *Device) UtilizationRates() (int, error) {
	var utilization C.nvmlUtilization_t
	if res := C.nvmlDeviceGetUtilizationRatesWrapper(d.handle, &utilization); res != C.NVML_SUCCESS {
		return 0, fmt.Errorf("failed to get utilization rates: %d", res)
	}
	return int(utilization.gpu), nil
}

func (d *Device) MemoryInfo() (uint64, error) {
	var memory C.nvmlMemory_t
	if res := C.nvmlDeviceGetMemoryInfoWrapper(d.handle, &memory); res != C.NVML_SUCCESS {
		return 0, fmt.Errorf("failed to get memory info: %d", res)
	}
	return uint64(memory.used), nil
}
