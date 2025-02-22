/**
 * @create by: Keeno
 * @description:
 * @create time: 2024/7/30 10:50
 *
 */
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"log"
	"os"
	"strconv"
)

const (
	highDimPattern5       = 5
	highDimPattern6       = 6
	highDimPattern5Active = 1000
	highDimPattern6Active = 3000
)

/**
 * @create by: Keeno
 * @description: 最左分割
 * @create time: 2024/7/30 09:31
 *
 */
func leftmost(arrs *mat.Dense) [][]int {
	rows, cols := arrs.Dims()
	Tarrs := mat.NewDense(cols, rows, nil)
	Tarrs.CloneFrom(arrs.T())

	// 初始化变量
	var splitIndex int = -1
	var splitNibbles []int
	splitNibblesIndex := make([][]int, 16)

	// 遍历每一列
	for i := 0; i < 32; i++ {
		colArr := mat.Row(nil, i, Tarrs)
		coli := mat.NewDense(1, rows, colArr)
		if mat.Max(coli) == mat.Min(coli) {
			continue
		}
		splitIndex = i

		splits := make([]bool, 16)
		// 添加分割维度的不同取值
		for index, value := range colArr {
			if !splits[int(value)] {
				splits[int(value)] = true
				splitNibbles = append(splitNibbles, int(value))
				var curSplitNibblesIndex []int
				splitNibblesIndex[int(value)] = append(curSplitNibblesIndex, index)
			} else {
				splitNibblesIndex[int(value)] = append(splitNibblesIndex[int(value)], index)
			}
		}
		break

	}

	// 检查是否找到了合适的列
	if splitIndex == -1 {
		return nil
	}

	// 生成结果
	result := make([][]int, len(splitNibbles))
	for _, value := range splitNibbles {
		result = append(result, splitNibblesIndex[value])

	}

	return result
}

/**
 * @create by: Keeno
 * @description: 最大覆盖分割
 * @create time: 2024/7/30 09:31
 *
 */
func maxCovering(arrs *mat.Dense) [][]int {
	rows, cols := arrs.Dims()
	Tarrs := mat.NewDense(cols, rows, nil)
	Tarrs.CloneFrom(arrs.T())

	var Covering []float64
	leftmostIndex := -1
	leftmostCovering := -1.0

	for i := 0; i < 32; i++ {
		colArr := mat.Row(nil, i, Tarrs)
		coli := mat.NewDense(1, rows, colArr)
		if mat.Max(coli) == mat.Min(coli) {
			Covering = append(Covering, -1)
		} else {
			splits := make([]int, 16)
			//统计每个取值的地址个数
			for _, value := range colArr {
				splits[int(value)] += 1
			}
			//统计非离群地址个数
			var sum float64
			for _, value := range splits {
				if value != 1 {
					sum += float64(value)
				}

			}
			if leftmostIndex == -1 {
				leftmostIndex = i
				leftmostCovering = sum
			}
			Covering = append(Covering, sum)
		}
	}

	index := floats.MaxIdx(Covering)
	if floats.Max(Covering)-leftmostCovering <= float64(index-leftmostIndex) {
		index = leftmostIndex
	}

	colArr := mat.Row(nil, index, Tarrs)

	splits := make([][]int, 16)
	valueIsEmpytList := make([]bool, 16)
	//统计每个取值的地址个数
	for nibbleIndex, value := range colArr {
		valueIsEmpytList[int(value)] = true
		if len(splits[int(value)]) == 0 {
			var splitNibbles []int
			splits[int(value)] = append(splitNibbles, nibbleIndex)

		} else {
			splits[int(value)] = append(splits[int(value)], nibbleIndex)
		}
	}
	var result [][]int
	for valueIndex, _ := range valueIsEmpytList {
		if valueIsEmpytList[valueIndex] {
			result = append(result, splits[valueIndex])
		}
	}
	return result
}

/**
 * @create by: Keeno
 * @description: 最小熵分割
 * @create time: 2024/7/30 09:41
 *
 */
func minEntropy(arrs *mat.Dense) [][]int {
	rows, cols := arrs.Dims()
	Tarrs := mat.NewDense(cols, rows, nil)
	Tarrs.CloneFrom(arrs.T())

	minLeftEntropyIndex := -1
	minLeftEntropyValue := 17
	var splitNibbles []int

	for i := 0; i < cols && i < 32; i++ {
		splits := make([]int, 16)
		uniqueElements := make(map[int]bool)
		for j := 0; j < rows; j++ {
			value := int(Tarrs.At(i, j))
			splits[value]++
			uniqueElements[value] = true
		}

		uniqueCount := len(uniqueElements)
		if uniqueCount == 1 {
			continue
		}

		if uniqueCount == 2 {
			minLeftEntropyIndex = i
			for nibble, count := range splits {
				if count > 0 {
					splitNibbles = append(splitNibbles, nibble)
				}
			}
			break
		}

		if uniqueCount < minLeftEntropyValue {
			minLeftEntropyValue = uniqueCount
			splitNibbles = []int{}
			for nibble, count := range splits {
				if count > 0 {
					splitNibbles = append(splitNibbles, nibble)
				}
			}
			minLeftEntropyIndex = i
		}
	}

	if minLeftEntropyIndex == -1 {
		return nil
	}

	// 生成结果
	result := make([][]int, len(splitNibbles))
	for i, nibble := range splitNibbles {
		var indices []int
		for j := 0; j < rows; j++ {
			if int(Tarrs.At(minLeftEntropyIndex, j)) == nibble {
				indices = append(indices, j)
			}
		}
		result[i] = indices
	}

	return result
}

/**
 * @create by: Keeno
 * @description: 最右分割
 * @create time: 2024/7/30 10:45
 *
 */
func rightmost(arrs *mat.Dense) [][]int {
	rows, cols := arrs.Dims()
	Tarrs := mat.NewDense(cols, rows, nil)
	Tarrs.CloneFrom(arrs.T())

	// 初始化变量
	var splitIndex int = -1
	var splitNibbles []int
	splitNibblesIndex := make([][]int, 16)
	// 遍历每一列
	for i := 31; i < cols && i >= 0; i-- {
		colArr := mat.Row(nil, i, Tarrs)
		coli := mat.NewDense(1, rows, colArr)
		if mat.Max(coli) == mat.Min(coli) {
			continue
		}
		splitIndex = i

		splits := make([]bool, 16)
		// 添加分割维度的不同取值
		for index, value := range colArr {
			if !splits[int(value)] {
				splits[int(value)] = true
				splitNibbles = append(splitNibbles, int(value))
				var curSplitNibblesIndex []int
				splitNibblesIndex[int(value)] = append(curSplitNibblesIndex, index)
			} else {
				splitNibblesIndex[int(value)] = append(splitNibblesIndex[int(value)], index)
			}
		}
		break
	}

	// 检查是否找到了合适的列
	if splitIndex == -1 {
		return nil
	}

	// 生成结果
	result := make([][]int, len(splitNibbles))
	for _, value := range splitNibbles {
		result = append(result, splitNibblesIndex[value])

	}

	return result
}

/**
 * @create by: Keeno
 * @description: 聚类
 * @create time: 2024/7/30 09:31
 *
 */
func DHC(arrs *mat.Dense, dhcType string) []*mat.Dense {
	q := list.New()
	q.PushBack(arrs)
	var regionsArrs []*mat.Dense

	for q.Len() > 0 {
		element := q.Back()
		q.Remove(element)
		arr := element.Value.(*mat.Dense)

		rows, _ := arr.Dims()
		if rows < 16 {
			regionsArrs = append(regionsArrs, arr)
			continue
		}
		var splits [][]int
		//四种分裂规则
		switch dhcType {
		case "left":
			splits = leftmost(arr)
		case "right":
			splits = rightmost(arr)
		case "min":
			splits = minEntropy(arr)
		case "max":
			splits = maxCovering(arr)
		}

		if len(splits) > 0 {
			for _, s := range splits {
				if len(s) > 1 {
					subArr := mat.NewDense(len(s), arr.RawMatrix().Cols, nil)
					for i, rowIndex := range s {
						subArr.SetRow(i, arr.RawRowView(rowIndex))
					}
					q.PushBack(subArr)
				}
			}
		}

	}
	//return regionsArrs
	return RemoveOutliers(regionsArrs)
}

/**
 * @create by: Keeno
 * @description: 移除模式中的离群种子
 * @create time: 2024/8/2 18:24
 *
 */
func RemoveOutliers(regionsArrs []*mat.Dense) []*mat.Dense {
	var allRemoverArrs []*mat.Dense
	if len(regionsArrs) > 0 {
		for _, region := range regionsArrs {
			subRegions := OutlierDetect(region)
			for _, subRegion := range subRegions {
				allRemoverArrs = append(allRemoverArrs, subRegion)
			}
		}
	}
	return allRemoverArrs
}

/**
 * MDHC algorithm.
 * Clustering seed addresses to construct 4 IPv6 address space trees,
 * returning the low-dimensional patterns and high-dimensional patterns (pattern dimensions of 5 or 6)
 *in the space trees.
 */
func MDHC(seedsFile string) ([]string, []string, []string) {

	var lowDimPatterns []string
	var highDimPatterns_5 []string
	var highDimPatterns_6 []string
	dhcTypes := []string{"min", "left", "right", "max"}
	for _, dhcType := range dhcTypes {
		//将IPv6列表转化为二维数组
		arrs := parseIPv6File(seedsFile)
		if arrs == nil {
			fmt.Println("No valid IPv6 addresses found.")
			return nil, nil, nil
		}
		//DHC聚类
		regions := DHC(arrs, dhcType)
		//Deduplicate
		regions = uniqueMatrices(regions)
		for _, region := range regions {
			patternArray, wildcardNum := initSubspace(region)
			if wildcardNum == 0 {
				continue
			} else if wildcardNum <= 6 {
				pattern, _ := ConvertStringArrayToIPv6(patternArray)
				searchKeyword, _ := filter.Search(pattern)
				if len(searchKeyword) > 0 {
					continue
				}
				if wildcardNum <= 4 {
					lowDimPatterns = append(lowDimPatterns, pattern)
				} else if wildcardNum <= 5 {
					highDimPatterns_5 = append(highDimPatterns_5, pattern)
				} else {
					highDimPatterns_6 = append(highDimPatterns_6, pattern)
				}
			}
		}
	}

	fmt.Println("lowDimPatterns:", len(lowDimPatterns), ", highDimPatterns: ", len(highDimPatterns_5)+len(highDimPatterns_6))
	return lowDimPatterns, highDimPatterns_5, highDimPatterns_6

}

/**
 * Generate IPv6 target addresses in the low-dimensional pattern space
 */
func generateTargetAddress(patternFile []string, targetFile string) {
	outputFile, err := os.Create(targetFile)
	writer := bufio.NewWriter(outputFile)
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()

	countLowPattern := len(patternFile)
	var countAddress float64
	for _, pattern := range patternFile {
		println("Current remaining unfiltered prefixes: ", countLowPattern)
		countLowPattern--
		addresses := expandWildcard(pattern)
		for _, ip := range addresses {
			countAddress++
			fmt.Fprintf(writer, "%s\n", ip)
		}

	}
	writer.Flush()

}

/**
 * Save high-dimensional patterns
 */
func saveHighDimPatterns(highDimPatterns_5 []string, highDimPatterns_6 []string, highDimPatternFile string) {
	outputFile, err := os.Create(highDimPatternFile + "_5")
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	writer := bufio.NewWriter(outputFile)
	for _, pattern := range highDimPatterns_5 {
		fmt.Fprintf(writer, "%s\n", pattern)
	}
	writer.Flush()

	outputFile, err = os.Create(highDimPatternFile + "_6")
	writer = bufio.NewWriter(outputFile)
	for _, pattern := range highDimPatterns_6 {
		fmt.Fprintf(writer, "%s\n", pattern)
	}
	writer.Flush()
}

/**
 * Based on IPv6 active addresses in the low-dimensional pattern space as the filtering condition,
 * filter out active high-dimensional address patterns
 */
func feedback(patternFile string, patternDimension int, addressFile string, outputFileName string) {

	if addressFile == "" {
		fmt.Println("For activeAddress must be provided.")
		os.Exit(1)
	}
	if patternDimension == 0 {
		fmt.Println("For patternDimension must be provided. High-dimensional pattern's pattern dimension can be 5 or 6")
		os.Exit(1)
	}

	if patternFile == "" {
		fmt.Println("For highDimPatternFileName must be provided.")
		os.Exit(1)
	}

	if outputFileName == "" {
		outputFileName = "targetAddress_" + strconv.Itoa(patternDimension)
	}

	if patternDimension != 5 || patternDimension != 6 {
		fmt.Println("The pattern dimension of the high-dimensional address pattern should be 5 or 6.")
	}

	// Read high-dimensional patterns.
	file, _ := os.Open(patternFile)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		filter.Add([]string{scanner.Text()})
	}

	// Read the active address and match it with the high-dimensional patterns.
	file, err := os.Open(addressFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	prefixes := make(map[string]int)

	for scanner.Scan() {
		ip := scanner.Text()
		searchKeyword, _ := filter.Search(ip)
		if len(searchKeyword) > 0 {
			value, exists := prefixes[searchKeyword]
			if exists {
				containsInBloomFilter(ip)
				prefixes[searchKeyword] = value + 1
			} else {
				prefixes[searchKeyword] = 1
			}
		}
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)

	count := len(prefixes)
	var threshold int
	if patternDimension == highDimPattern5 {
		threshold = highDimPattern5Active
	} else {
		threshold = highDimPattern6Active
	}
	// Generate IPv6 target addresses in the active high-dimensional pattern space.
	for key, value := range prefixes {
		println("Current remaining unfiltered prefixes: ", count)
		count--

		if value >= threshold {
			addresses := expandWildcard(key)
			for _, ip := range addresses {
				fmt.Fprintf(writer, "%s\n", ip)
			}
		}
	}
	writer.Flush()
}
