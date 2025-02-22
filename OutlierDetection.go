// /**
// * @create by: Keeno
// * @description:
// * @create time: 2024/8/2 08:00
// *
// */
package main

import (
	"container/list"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func IoslatedForest(Weights []float64, Tarr_i []float64) {
	counts := make(map[float64][]int, 16)
	OutlierNum := 16
	if len(Tarr_i) < OutlierNum {
		OutlierNum = len(Tarr_i)
	}
	valueSingle := make(map[float64]bool, 16)
	for index, split := range Tarr_i {
		if counts[split] == nil {
			var indexList []int
			indexList = append(indexList, index)
			counts[split] = indexList
			valueSingle[split] = true
		} else {
			OutlierNum = OutlierNum - 1
			valueSingle[split] = false
		}
	}

	for key, value := range valueSingle {
		if value {
			index := counts[key][0]
			Weights[index] += 1.0 / float64(OutlierNum)
		}
	}
}

func FourD(Weights []float64) []int {
	//if len(Weights) <= 2 {
	//	return []int{}
	//}
	outLierIndex := floats.MaxIdx(Weights)

	outRemovedWeights := append([]float64{}, Weights...)
	outRemovedWeights = append(outRemovedWeights[:outLierIndex], outRemovedWeights[outLierIndex+1:]...)
	//标准差
	outRemovedD := stat.StdDev(outRemovedWeights, nil)
	//均值
	outRemovedAvg := stat.Mean(outRemovedWeights, nil)

	if Weights[outLierIndex]-outRemovedAvg > 3*outRemovedD {
		return append([]int{outLierIndex}, FourD(outRemovedWeights)...)
	} else {
		return []int{}
	}
}

func iterDevide(arrs *mat.Dense) []*mat.Dense {
	q := list.New()
	q.PushBack(arrs)
	var regionsArrs []*mat.Dense
	for q.Len() > 0 {
		element := q.Back()
		q.Remove(element)
		arr := element.Value.(*mat.Dense)
		splits := maxCovering(arr)

		found := false
		for _, s := range splits {
			if len(s) == 1 {
				regionsArrs = append(regionsArrs, arr)
				found = true
				break
			}
		}
		if !found {
			for _, s := range splits {
				subArr := mat.NewDense(len(s), arr.RawMatrix().Cols, nil)
				for i, rowIndex := range s {
					subArr.SetRow(i, arr.RawRowView(rowIndex))
				}
				q.PushBack(subArr)
			}
		}
	}
	return regionsArrs
}

func removeRow(yuanlai *mat.Dense, rowIndex int) *mat.Dense {
	rows, cols := yuanlai.Dims()
	// 创建一个新的 Dense 矩阵，尺寸比原来少一行
	newMat := mat.NewDense(rows-1, cols, nil)

	// 拷贝原始矩阵中除了第 rowIndex 行之外的所有行到新矩阵中
	for r := 0; r < rows; r++ {
		if r < rowIndex {
			// 复制 rowIndex 之前的行
			newMat.SetRow(r, mat.Row(nil, r, yuanlai))
		} else if r > rowIndex {
			// 复制 rowIndex 之后的行，同时调整索引以适应新矩阵
			newMat.SetRow(r-1, mat.Row(nil, r, yuanlai))
		}
	}
	return newMat
}

func OutlierDetect(arrs *mat.Dense) []*mat.Dense {
	var regionsArrs []*mat.Dense
	numRows, numCols := arrs.Dims()
	if numRows == 1 {
		return regionsArrs
	}
	if numRows == 2 {
		_, dimension := initSubspace(arrs)
		if dimension > 5 {
			return regionsArrs
		} else {
			regionsArrs = append(regionsArrs, arrs)
			return regionsArrs
		}

	}
	Tarrs := mat.NewDense(numCols, numRows, nil)
	Tarrs.CloneFrom(arrs.T())

	freeDimensionNum := 0
	Weights := make([]float64, numRows)

	for i := 0; i < 32; i++ {
		colArr := mat.Row(nil, i, Tarrs)
		coli := mat.NewDense(1, numRows, colArr)
		if mat.Max(coli) == mat.Min(coli) {
			continue
		}
		freeDimensionNum++
		IoslatedForest(Weights, colArr)
	}

	for _, oW := range FourD(Weights) {
		arrs = removeRow(arrs, oW)
	}

	patterns := iterDevide(arrs)

	return patterns
}
