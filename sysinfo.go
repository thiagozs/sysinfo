package sysinfo

import (
	"fmt"
	"sync"
	"syscall"
	"time"
)

// Sysinfo system info model
// Work on Unix only ***
type Sysinfo struct {
	Uptime       time.Duration // time since boot
	Loads        [3]float64    // 1, 5, and 15 minute load averages
	Procs        uint64        // number of current processes
	TotalRAM     uint64        // total usable main memory size [kB]
	FreeRAM      uint64        // available memory size [kB]
	SharedRAM    uint64        // amount of shared memory [kB]
	BufferRAM    uint64        // memory used by buffers [kB]
	TotalSwap    uint64        // total swap space size [kB]
	FreeSwap     uint64        // swap space still available [kB]
	TotalHighRAM uint64        // total high memory size [kB]
	FreeHighRAM  uint64        // available high memory size [kB]
	mu           sync.Mutex    // protects fields
}

// Get the sysinfo
func Get() (*Sysinfo, error) {

	si := &syscall.Sysinfo_t{}
	sis := &Sysinfo{}

	err := syscall.Sysinfo(si)
	if err != nil {
		return nil, err
	}

	scale := 65536.0               // magic
	unit := uint64(si.Unit) * 1024 // kB

	sis.Loads[0] = float64(si.Loads[0]) / scale
	sis.Loads[1] = float64(si.Loads[1]) / scale
	sis.Loads[2] = float64(si.Loads[2]) / scale

	sis.Uptime = time.Duration(si.Uptime) * time.Second
	sis.Procs = uint64(si.Procs)

	sis.TotalRAM = uint64(si.Totalram) / unit
	sis.FreeRAM = uint64(si.Freeram) / unit
	sis.BufferRAM = uint64(si.Bufferram) / unit
	sis.TotalSwap = uint64(si.Totalswap) / unit
	sis.FreeSwap = uint64(si.Freeswap) / unit
	sis.TotalHighRAM = uint64(si.Totalhigh) / unit
	sis.FreeHighRAM = uint64(si.Freehigh) / unit

	sis.mu = sync.Mutex{}

	return sis, nil
}

// String return a formated results of sysinfo without locking
func (si *Sysinfo) String() string {
	return fmt.Sprintf("uptime\t\t%v\nload\t\t%2.2f %2.2f %2.2f\nprocs\t\t%d\n"+
		"ram  total\t%d kB\nram  free\t%d kB\nram  buffer\t%d kB\n"+
		"swap total\t%d kB\nswap free\t%d kB",
		si.Uptime, si.Loads[0], si.Loads[1], si.Loads[2], si.Procs,
		si.TotalRAM, si.FreeRAM, si.BufferRAM,
		si.TotalSwap, si.FreeSwap,
	)
}

// ToString return a formated result about sysinfo with lock
func (si *Sysinfo) ToString() string {
	defer si.mu.Unlock()
	si.mu.Lock()
	return si.String()
}
