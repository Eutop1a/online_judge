package utils

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

func GetInputAndExpectedFromFile(inputDir, expectedDir string) ([]string, []string, error) {

	// 读取所有输入文件
	inputFiles, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, nil, err
	}

	// 创建两个字符串切片用于存储文件内容
	var Inputs []string
	var Expected []string

	// 确保文件名排序，保证匹配顺序
	sort.Slice(inputFiles, func(i, j int) bool {
		return inputFiles[i].Name() < inputFiles[j].Name()
	})

	// 遍历输入文件
	for _, file := range inputFiles {
		// 构建对应的输出文件名
		expectedFileName := strings.Replace(file.Name(), "input", "expected", 1) // 更正替换逻辑
		expectedFileName = strings.Replace(expectedFileName, ".in", ".out", 1)   // 确保扩展名正确替换

		// 读取输入文件内容
		inputContent, err := ioutil.ReadFile(filepath.Join(inputDir, file.Name()))
		if err != nil {
			return nil, nil, err
		}
		Inputs = append(Inputs, string(inputContent))

		// 读取输出文件内容
		expectedContent, err := ioutil.ReadFile(filepath.Join(expectedDir, expectedFileName))
		if err != nil {
			return nil, nil, err
		}
		Expected = append(Expected, string(expectedContent))
	}
	return Inputs, Expected, nil
}
