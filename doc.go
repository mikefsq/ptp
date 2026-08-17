// Package ptp implements Picture Transfer Protocol to cameras over USB
//
// # Layout
//
//	ptp/       this package: container framing, sessions, objects, property
//	           descriptors, events. Pure Go — no cgo, so its tests run
//	           anywhere and it cross-compiles freely.
//	ptp/usb    the per-OS USB transports. The only cgo in the module.
//	ptp/fuji   Fujifilm X and GFX
//	ptp/sony   Sony Alpha
//	ptp/canon  not implemented
//	ptp/nikon  not implemented
//
// # A vendor package defines the items that are not standardized by PTP
// such as opcode, property tables, its capture and exposure semantics,
// and the names of its vendor-specific response codes.
package ptp
