package stats

import (
	"os/exec"
	"strings"
	"strconv"
)
//wmic cpu get loadpercentage

func GetStats() (cpu, mem float32){
	out, err := exec.Command("wmic","cpu", "get", "loadpercentage").Output()
	if err != nil {
		return -1
	}
	strOut := string(out)
	line := strOut[strings.IndexAny(strOut,"1234567890"):strings.LastIndexAny(strOut, "1234567890")+1]
	i, err := strconv.Atoi(line)
	if err != nil{
		return -2
	}
	return float32(i), -1
}