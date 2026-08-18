//go:build !linux && !darwin

package usb

import "context"

func hotplug(context.Context) (<-chan HotplugEvent, error) { return nil, ErrNoHotplug }
