package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	BitSet = make([]byte, 1<<30)
	filter = New()
)

func main() {

	// operations, MDHC, expand, feedback, extremeExpansion, convert
	var operate string
	// the IPv6 seed address file
	var seedSetFile string
	// the target address file
	var targetAddressFile string
	// the high-dimensional pattern file
	var highDimPatternFile string
	// the active address file
	var activeAddressFile string
	// the pattern dimension of the high-dimensional patterns that need to be filtered
	var patternDimension int
	// the Seed address source file
	var seedSourceFileName string
	// the number of seed address sets constructed.
	//Randomly select a certain number of active addresses from
	//the seed address source as the seed address set by downsampling and other means
	var seedSetNum int
	// Single seed address set size
	var seedSetSize int
	// Seed address set filename prefix
	var seedSetPrefix string

	flag.StringVar(&operate, "o", "MDHC", "The operation to perform (MDHC, convert, feedback, expand)")
	flag.StringVar(&seedSetFile, "s", "random100K1", "The name of the seed file")
	flag.StringVar(&targetAddressFile, "t", "targetAddress", "The name of the target address file")
	flag.StringVar(&highDimPatternFile, "h", "highDimPattern", "The name of the high-dimensional pattern file")
	flag.StringVar(&activeAddressFile, "a", "", "The name of the active address file")
	flag.IntVar(&patternDimension, "p", 0, "The pattern dimension (positive integer)")
	flag.StringVar(&seedSourceFileName, "S", "hitlist_2024_07_20", "The name of the active address file")
	flag.IntVar(&seedSetNum, "num", 10, "The number of seed address sets(positive integer)")
	flag.IntVar(&seedSetSize, "size", 1000000, "The number of active IPv6 addresses in the seed address set(positive integer)")
	flag.StringVar(&seedSetPrefix, "prefix", "random100K", "The seed address set file name prefix")
	flag.Parse()

	if operate == "" {
		fmt.Println("For operate must be provided. Operate can be MDHC, format, feedback, or expand.")
		os.Exit(1)
	}

	switch operate {
	case "expand":
		// Sampling-based IPv6 seed address expansion strategy
		// command: ./6Massive -o expand -S hitlist_2024_07_20 -num 10 -size 100000 -prefix random100K
		// Input.
		//			-S:      Seed address source
		//			-num:    The number of seed address set
		//			-size:   The size of single seed address set
		//          -prefix: Name prefix of seed address set
		// Output.
		//			All seed address sets of the same size
		//
		//

		// Using the Hitlist of July 24, 2024, as the seed address source, 6Massive constructs 10 sets of seed addresses,
		//with each set having a scale of 100,000 being more effective.
		//The parameters can be adjusted according to actual needs.
		// Execute the shell script
		cmd := exec.Command("bash", "expand.sh", seedSourceFileName, strconv.Itoa(seedSetNum), strconv.Itoa(seedSetSize), seedSetPrefix)
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	case "MDHC":
		// MDHC (merged divisive hierarchical clustering) strategy
		// command: ./6Massive -o MDHC -s random100K1 -t targetAddress -h highDimPattern
		// Input.
		//			-s:      Seed address set
		// Output.
		//			-t:      Target address set file
		//		    -h:		 high-dimensional address pattern file
		//
		//
		// Construct 4 IPv6 address space trees for seed addresses through MDHC strategy,
		//generate low-dimensional patterns and high-dimensional patterns, and generate
		//IPv6 target addresses in the low-dimensional address pattern space for probing

		//By constructing four IPv6 address space trees through the MDHC strategy,
		//merging low-dimensional patterns and high-dimensional patterns.
		lowDimPatterns, highDimPatterns_5, highDimPatterns_6 := MDHC(seedSetFile)

		//Generate IPv6 target addresses in the low-dimensional pattern space
		generateTargetAddress(lowDimPatterns, targetAddressFile)
		//Save high-dimensional patterns for generating more IPv6 target addresses in the feedback stage.
		saveHighDimPatterns(highDimPatterns_5, highDimPatterns_6, highDimPatternFile)

		//Subsequent scanning of IPv6 target addresses to obtain active IPv6 addresses can be performed
		//using the Zmapv6 tool.
		//Command:
		//sudo zmap --probe-module=icmp6_echoscan --ipv6-target-file=targetAddress  --output-file=activeAddress --ipv6-source-ip=2001::1(Host IP)   --bandwidth=30M --cooldown-time=4
	case "feedback":
		// Feedback strategy based on pattern space intersection
		// command: ./6Massive -o feedback -a activeAddress -h highDimPattern5 -p 5 -t targetAddress
		// Input.
		//			-h:      high-dimensional address pattern file
		//          -p:      Pattern dimensions of high-dimensional address patterns (5 or 6)
		//          -a       Active addresses in low-dimensional schema space
		// Output.
		//			-t:      Target addresses generated in the active high-dimensional schema space
		//
		//
		//Based on the experimental results, the pattern dimension (patternDimension) is better when it equals 5 or 6.
		//Execute feedback strategy, filter active high-dimensional patterns based on IPv6 active addresses
		//probed in the low-dimensional mode space, and generate IPv6 target addresses in its space.
		feedback(highDimPatternFile, patternDimension, activeAddressFile, targetAddressFile+strconv.Itoa(patternDimension))
	case "convert":
		// Before executing the feedback strategy, format the IPv6 active addresses scanned
		// command: ./6Massive -o convert -a activeAddress -t targetAddress
		// Input.
		//          -a       Active addresses in low-dimensional schema space (eg. 2001:4860:4860::8888)
		// Output.
		//			-t:      Active addresses after format (eg. 2001:4860:4860:0000:0000:0000:0000:8888)
		//
		convertIPv6toFull(activeAddressFile, targetAddressFile)
	case "extremeExpansion":
		// The -S parameter represents all target address sets that need to be generated for IPv6 target address generation, separated by _
		// command: ./6Massive -o extremeExpansion -S random1M1_random1M2_random1M3_random1M4_random1M5_random1M6_random1M7_random1M8_random1M9_random1M10 -t targetAddress
		// Input.
		//          -S       All seed address set files of the same size
		// Output.
		//			-t:      Target address
		//

		var lowDimPatterns []string
		if strings.Contains(seedSetFile, "_") {
			subStrings := strings.Split(seedSetFile, "_")
			for _, str := range subStrings {
				curLowDimPatterns, _, _ := MDHC(str)
				lowDimPatterns = append(lowDimPatterns, curLowDimPatterns...)
			}
		} else {
			lowDimPatterns, _, _ = MDHC(seedSetFile)
		}
		filter = New()
		lowDimPatterns = removeDuplicates(lowDimPatterns)
		//Generate IPv6 target addresses in the low-dimensional pattern space
		generateTargetAddress(lowDimPatterns, targetAddressFile)
	default:
		// If the operation type is invalid
		fmt.Println("Invalid operation. Please use 'MDHC', 'convert', 'feedback', or 'expand'.")
		os.Exit(1)
	}

}
