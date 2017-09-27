package util

import (
	"os"
	"runtime"
	"runtime/pprof"
)

// Profiler is a structure to simplify generating CPU and memory profile files
// from a run of a program.
type Profiler struct {
	cpuFile    string
	memoryFile string
}

func createProfileFile(outfile string) (fp *os.File, err error) {
	fp, err = os.Create(outfile)
	if err != nil {
		err = E.Annotate(err, "Could not create profile file", outfile)
	}
	return
}

// SetupProfiler sets up an empty profiler with Options. The options should
// contain keys "profile-cpu-file" which is the file name for the CPU
// profile. Also the memory profile will be written to value of key
// "profile-mem-file".
func SetupProfiler(opts Options) (*Profiler, error) {
	var p *Profiler = &Profiler{
		cpuFile:    opts.Get("profile-cpu-file", ""),
		memoryFile: opts.Get("profile-mem-file", ""),
	}

	if p.cpuFile != "" {
		fp, err := createProfileFile(p.cpuFile)
		if err != nil {
			return nil, err
		}
		runtime.SetCPUProfileRate(1000)
		pprof.StartCPUProfile(fp)
	}

	return p, nil
}

// Close finishes the creation of the profile. Should be defer'd after running
// Setup.
func (p *Profiler) Close() error {
	if p.cpuFile != "" {
		pprof.StopCPUProfile()
	}
	if p.memoryFile != "" {
		var fp *os.File
		fp, err := createProfileFile(p.memoryFile)
		if err != nil {
			return err
		}
		defer fp.Close()
		pprof.WriteHeapProfile(fp)
	}
	return nil
}
