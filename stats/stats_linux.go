package stats

import(
	"io/ioutil"
	"strings"
	"strconv"
	"time"
)
/*
You just need to sum the 2nd (user), 3rd (nice) and 4th (system) row, and divide it by all the available time (5th row is the idle time). Pseudo-code example: 

tmp = $2 + $3 + $4 
usage = $tmp / ($tmp + $5) 

user: normal processes executing in user mode
nice: niced processes executing in user mode
system: processes executing in kernel mode
idle: twiddling thumbs
iowait: waiting for I/O to complete
irq: servicing interrupts
softirq: servicing softirqs

*/
func getMem() (mem float32){
	content, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		//handle
	}
	var totalM float32
	var freeM float32 = 0
	lines := strings.Split(string(content), "\n")
	for _, line := range(lines){
		fields := strings.Fields(line)
		if len(fields) > 0 {
			if fields[0] == "MemTotal:" {
				val, err := strconv.Atoi(fields[1])
				if err != nil {
					//handle
				}
				totalM = float32(val)
			} else if fields[0] == "MemFree:" {
				val, err := strconv.Atoi(fields[1])
				if err != nil {
					//handle
				}
				freeM += float32(val)
			} else if fields[0] == "Cached:" {
				val, err := strconv.Atoi(fields[1])
				if err != nil {
					//handle
				}
				freeM += float32(val)
				} else if fields[0] == "Buffers:" {
				val, err := strconv.Atoi(fields[1])
				if err != nil {
					//handle
				}
				freeM += float32(val)
			}
		}
	}
	return (totalM-freeM)/totalM*100
}

func getCpuRaw() (total, idle float32){
	content, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		//handle
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			nfields := len(fields)
			for i := 1; i < nfields; i++ {
				val, err := strconv.Atoi(fields[i])
				if err != nil {
					//handle
				}
				total += float32(val)
				if i == 4 {
					idle = float32(val)
				}
			}
			return
		}
	}
	return
}

func GetStats() (cpu, mem float32){
	total0, idle0 := getCpuRaw()
	time.Sleep(time.Millisecond * 200)
	total1, idle1 := getCpuRaw()
	diffTotal := total1-total0
	cpu = (diffTotal-(idle1-idle0))/diffTotal*100
	if cpu < 0 {
		cpu = 0
	}
	return cpu, getMem()
}
