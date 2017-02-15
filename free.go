package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var unitB *bool // output in bytes
var unitK *bool // output in kilobytes
var unitM *bool // output in megabytes
var unitG *bool // output in gigabytes
var human *bool // human-readable output: show unit abbrev, 2 decimal places,
//                        dynamic units unless specified other

var units = [6]string{"B", "K", "M", "G", "T", "P"} // up to petabytes

var base = 1024

// check binary is available, execute it and capture output
func execOutput(cmd string, args ...string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal(cmd + ": not found")
	}

	var output []byte
	output, err = exec.Command(path, args...).Output()
	if err != nil {
		log.Fatal(cmd + ": failed")
	}

	return string(output)
}

// extract Ith group from the string
func reGrp(r string, s string, i int) string {
	result := regexp.MustCompile(r).FindStringSubmatch(s)
	if len(result) < i {
		return ""
	}

	return result[i]
}

// convert string to int
func toInt(s string, id string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(id + " = " + s + ": not a number")
	}
	return result
}

func toFloat(s string, id string) float64 {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(id + " = " + s + ": not a number")
	}
	return result
}

func scaleSize(v, p int, format string) string {
	if p > -1 {
		u := math.Pow(float64(base), float64(p))
		var r string
		if *human {
			r = fmt.Sprintf(format+"%s", float64(v)/u, units[p])
		} else {
			r = fmt.Sprintf(format, float64(v)/u)
		}
		return r
	}

	// dynamic sizing
	exp := 0
	u := 1
	r := v / u
	for r > base {
		exp++
		u = u * 1024
		r = v / u
	}

	return fmt.Sprintf(format+"%s", float64(v)/float64(u), units[exp])
}

func main() {

	unitB = flag.Bool("b", false, "show output in bytes")
	unitK = flag.Bool("k", false, "show output in kilobytes")
	unitM = flag.Bool("m", false, "show output in megabytes")
	unitG = flag.Bool("g", false, "show output in gigabytes")
	human = flag.Bool("h", false, "show human-readable output")
	flag.Parse()

	var unitPower int    // dynamic uniting
	var pofmt = " %11.f" // padded format for nubmers
	var ofmt = " %.f"    // right-aligned non-padded format for numbers

	switch {
	case *unitB:
		unitPower = 0
		if *human {
			pofmt = " %10.f"
			ofmt = " %.f"
		}
	case *unitK:
		unitPower = 1
		if *human {
			pofmt = " %10.f"
			ofmt = " %.f"
		}
	case *unitM:
		unitPower = 2
		if *human {
			pofmt = " %10.2f"
			ofmt = " %.2f"
		}
	case *unitG:
		unitPower = 3
		if *human {
			pofmt = " %10.2f"
			ofmt = " %.2f"
		}
	default: // dynamic uniting
		unitPower = -1
		pofmt = " %10.2f"
		ofmt = " %.2f"
	}

	// Getting data
	vmStat := execOutput("vm_stat", "")
	pageSize := toInt(reGrp("page size of (\\d+) bytes", vmStat, 1), "pageSize")
	pagesFree := toInt(reGrp("Pages free: \\s*(\\d+)", vmStat, 1), "pagesFree")
	anonPages := toInt(reGrp("Anonymous pages: \\s*(\\d+)", vmStat, 1), "anonPages")
	pagesWired := toInt(reGrp("Pages wired down: \\s*(\\d+)", vmStat, 1), "pagesWired")
	pagesPurgeable := toInt(reGrp("Pages purgeable: \\s*(\\d+)", vmStat, 1), "pagesPurgeable")
	pagesCompSrc := toInt(reGrp("Pages stored in compressor: \\s*(\\d+)", vmStat, 1), "pagesCompSrc")
	pagesCompRes := toInt(reGrp("Pages occupied by compressor: \\s*(\\d+)", vmStat, 1), "pagesCompRes")
	filePages := toInt(reGrp("File-backed pages: \\s*(\\d+)", vmStat, 1), "filePages")

	compRatio := (pagesCompSrc - pagesCompRes) * 100 / pagesCompSrc

	hwMemSize := execOutput("sysctl", "-n", "hw.memsize")
	memSize := toInt(strings.Trim(hwMemSize, " \n"), "hw.memsize")

	// swapusage is always in megabytes, internally will store in bytes
	swapUsage := execOutput("sysctl", "-noh", "vm.swapusage")
	swapTotal := int(toFloat(reGrp("total = (\\S+)M", swapUsage, 1), "swapTotal") * 1024 * 1024)
	swapUsed := int(toFloat(reGrp("used = (\\S+)M", swapUsage, 1), "swapUsed") * 1024 * 1024)
	swapFree := int(toFloat(reGrp("free = (\\S+)M", swapUsage, 1), "swapFree") * 1024 * 1024)

	// output
	fmt.Printf("              total        used        free      appmem       wired   compressed (ratio)\n")
	// 1st line
	fmt.Printf("%-7s", "Mem:")
	fmt.Print(scaleSize(memSize, unitPower, pofmt))
	fmt.Print(scaleSize((memSize - pagesFree*pageSize), unitPower, pofmt))
	fmt.Print(scaleSize(pagesFree*pageSize, unitPower, pofmt))
	fmt.Print(scaleSize((anonPages-pagesPurgeable)*pageSize, unitPower, pofmt))
	fmt.Print(scaleSize(pagesWired*pageSize, unitPower, pofmt))
	fmt.Print("  ")
	fmt.Print(scaleSize(pagesCompSrc*pageSize, unitPower, ofmt))
	fmt.Print(" ->")
	fmt.Print(scaleSize(pagesCompRes*pageSize, unitPower, ofmt))
	fmt.Printf(" (%d%%)", compRatio)
	fmt.Println()
	// 2nd line
	fmt.Printf("%-19s", "+/- Cache:")
	fmt.Print(scaleSize(memSize-(pagesFree+filePages+pagesPurgeable)*pageSize, unitPower, pofmt))
	fmt.Print(scaleSize((pagesFree+filePages+pagesPurgeable)*pageSize, unitPower, pofmt))
	//fmt.Printf("  (%.2f%s fcache + %.2f%s purgeable)", filePages*pageSize/unit, unitName, pagesPurgeable*pageSize/unit, unitName)
	fmt.Println()
	// 3rd line
	fmt.Printf("%-7s", "Swap:")
	fmt.Print(scaleSize(swapTotal, unitPower, pofmt))
	fmt.Print(scaleSize(swapUsed, unitPower, pofmt))
	fmt.Print(scaleSize(swapFree, unitPower, pofmt))
	fmt.Println()

}
