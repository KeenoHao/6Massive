/**
 * @create by: Keeno
 * @description:
 * @create time: 2024/7/30 10:50
 *
 */
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"net"
	"os"
	"strings"
)

/**
 * @create by: Keeno
 * @description: 对[]*mat.Dense去重
 * @create time: 2024/7/30 10:51
 *
 */
func uniqueMatrices(matrices []*mat.Dense) []*mat.Dense {
	exists := make(map[string]bool)
	uniqueMatrices := []*mat.Dense{}

	for _, m := range matrices {
		// 将矩阵转换为字符串表示形式
		matrixString := matrixToString(m)

		if !exists[matrixString] {
			// 如果不存在，则添加到 uniqueMatrices 和 exists map 中
			uniqueMatrices = append(uniqueMatrices, m)
			exists[matrixString] = true
		}
	}

	return uniqueMatrices
}

// matrixToString 将 *mat.Dense 矩阵转换为字符串
func matrixToString(m *mat.Dense) string {
	// 获取矩阵的维度
	rows, cols := m.Dims()

	// 创建一个字符串缓冲区
	var b strings.Builder

	// 遍历矩阵并将元素添加到字符串缓冲区
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			b.WriteString(fmt.Sprintf("%f ", m.At(i, j)))
		}
		b.WriteString("\n")
	}

	return b.String()
}

/**
 * @create by: Keeno
 * @description: 将IPv6地址文件转化为矩阵
 * @create time: 2024/7/30 09:31
 *
 */
func parseIPv6File(filepath string) *mat.Dense {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var arrs [][]float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		ipv6 := net.ParseIP(ip)
		if ipv6 == nil {
			continue
		}
		ipv6 = ipv6.To16()
		row := make([]float64, 32)
		for i := 0; i < 16; i++ {
			row[i*2] = float64(ipv6[i] >> 4)
			row[i*2+1] = float64(ipv6[i] & 0x0F)
		}
		arrs = append(arrs, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	rows := len(arrs)
	if rows == 0 {
		return nil
	}
	cols := len(arrs[0])
	data := make([]float64, rows*cols)
	for i := 0; i < rows; i++ {
		copy(data[i*cols:], arrs[i])
	}

	return mat.NewDense(rows, cols, data)
}

/**
 * @create by: Keeno
 * @description: 将IPv6数组转为IPv6地址模式
 * @create time: 2024/7/30 09:32
 *
 */
func initSubspace(activeSeeds *mat.Dense) ([]string, int) {
	rows, cols := activeSeeds.Dims()

	if rows <= 1 {
		subspace := make([]string, cols)
		for i := range subspace {
			subspace[i] = "0"
		}
		return subspace, 0
	}

	Tars := mat.NewDense(cols, rows, nil)
	Tars.CloneFrom(activeSeeds.T())

	subspace := make([]string, cols)
	dimension := 0

	for i := 0; i < cols; i++ {
		uniqueValues := make(map[int]bool)
		var value int
		for j := 0; j < rows; j++ {
			value = int(Tars.At(i, j))
			uniqueValues[value] = true
		}

		if len(uniqueValues) > 1 {
			dimension++
			subspace[i] = "*"
		} else {
			subspace[i] = fmt.Sprintf("%x", value)
			//subspace[i] = strconv.Itoa(value)
		}
	}

	return subspace, dimension
}

/**
 * @create by: Keeno
 * @description: 根据IPv6地址模式生成所有情况地址
 * @create time: 2024/7/30 13:46
 *
 */
func expandWildcard(input string) []string {
	var results []string
	expandWildcardHelper("", input, &results)
	return results
}

func expandWildcardHelper(prefix, remaining string, results *[]string) {
	if len(remaining) == 0 {
		//resultHan := net.ParseIP(prefix).To16()
		//if !filter.Test(hashIPv6(resultHan)) {
		//	filter.Add(hashIPv6(resultHan))
		//	*results = append(*results, prefix)
		//}
		if containsInBloomFilter(prefix) {
			return
		}
		*results = append(*results, prefix)
		return
	}

	firstWildcard := strings.Index(remaining, "*")
	if firstWildcard == -1 {
		//resultHan := net.ParseIP(prefix + remaining).To16()
		//if !filter.Test(hashIPv6(resultHan)) {
		//	filter.Add(hashIPv6(resultHan))
		//	*results = append(*results, prefix+remaining)
		//}

		// Keys of Bloom filters
		if containsInBloomFilter(prefix + remaining) {
			return
		}
		*results = append(*results, prefix+remaining)
		return
	}

	// Find the part before the first '*'
	partBeforeWildcard := remaining[:firstWildcard]
	remaining = remaining[firstWildcard+1:]

	// Generate all combinations for '*'
	for i := 0; i <= 15; i++ {
		hexValue := fmt.Sprintf("%01x", i)
		expandWildcardHelper(prefix+partBeforeWildcard+hexValue, remaining, results)
	}
}

// ConvertStringArrayToIPv6 将字符串数组转换为 IPv6 地址格式
func ConvertStringArrayToIPv6(array []string) (string, int) {
	wildcardNum := 0
	var builder strings.Builder
	builder.Grow(len(array) + 7) // 预估大小

	for i, v := range array {
		if i > 0 && i%4 == 0 {
			builder.WriteByte(':')
		}
		if v == "*" {
			builder.WriteString("*")
			wildcardNum++
		} else {
			builder.WriteString(fmt.Sprintf("%01s", v))
		}
	}

	return builder.String(), wildcardNum
}

func murmur3(data []byte, seed uint32) uint32 {
	hash := seed
	for i := 0; i < len(data); i = i + 4 {
		k := binary.BigEndian.Uint32(data[i : i+4])
		k = k * 0xcc9e2d51
		k = (k << 15) | (k >> 17)
		k = k * 0x1b873593
		hash = hash ^ k
		hash = (hash << 13) | (hash >> 19)
		hash = hash*5 + 0xe6546b64
	}
	hash = hash ^ (hash >> 16)
	hash = hash * 0x85ebca6b
	hash = hash ^ (hash >> 13)
	hash = hash * 0xc2b2ae35
	hash = hash ^ (hash >> 16)
	return hash
}

func containsInBloomFilter(address string) bool {
	IPv6 := net.ParseIP(address).To16()
	i := murmur3(IPv6, 0x12345678)
	j := murmur3(IPv6, 0x87654321)
	// Check if the ip is in BitSet
	if BitSet[i/8]&(1<<(i%8)) != 0 && BitSet[j/8]&(1<<(j%8)) != 0 {
		return true
	}
	BitSet[i/8] |= (1 << (i % 8))
	BitSet[j/8] |= (1 << (j % 8))
	return false
}

/**
 * @create by: Keeno
 * @description: 将标准缩写格式IPv6地址扩展为完全格式IPv6地址
 * @create time: 2024/8/10 17:00
 *
 */
func expandIPv6Address(ip string) (string, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "", fmt.Errorf("invalid IPv6 address")
	}
	parsedIP = parsedIP.To16()
	if parsedIP == nil {
		return "", fmt.Errorf("not an IPv6 address")
	}
	// Format the IPv6 address in full
	fullFormat := fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		parsedIP[0], parsedIP[1], parsedIP[2], parsedIP[3],
		parsedIP[4], parsedIP[5], parsedIP[6], parsedIP[7],
		parsedIP[8], parsedIP[9], parsedIP[10], parsedIP[11],
		parsedIP[12], parsedIP[13], parsedIP[14], parsedIP[15])

	return fullFormat, nil
}

func convertIPv6toFull(originFileName, outputFileName string) {

	if originFileName == "" {
		fmt.Println("For activeAddress must be provided.")
		os.Exit(1)
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	file, err := os.Open(originFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		expandedIP, err := expandIPv6Address(ip)
		if err != nil {
			fmt.Println("IP:", ip, "   Error:", err)
			continue
		}
		fmt.Fprintf(writer, "%s\n", expandedIP)
	}
}

func normalizeIPv6(ip string) (isNotIPv6 bool, realIP string) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		isNotIPv6 = true
		return isNotIPv6, realIP
	}
	realIP = parsedIP.String()
	isNotIPv6 = false
	return isNotIPv6, realIP
}

func formateIP(originFile, targetFile string) {
	outputFile, err := os.Create(targetFile)
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()
	var count float64
	writer := bufio.NewWriter(outputFile)

	file, err := os.Open(originFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//isBig := false
	for scanner.Scan() {
		ip := scanner.Text()

		isNotIPv6, realIP := normalizeIPv6(ip)
		//println(realIP)
		if isNotIPv6 {
			println(len(ip), ":     ", ip)
			continue
		}
		if containsInBloomFilter(ip) {
			continue
		}
		fmt.Fprintf(writer, "%s\n", realIP)
		//if count > 1000000000 {
		//	isBig = true
		//}
		//if isBig {
		//	fmt.Fprintf(writer, "%s\n", realIP)
		//}
		count++

	}
	writer.Flush()
	fmt.Println("总数量:", count)
}
