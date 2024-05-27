package services

import "online-judge/dao/mysql"

func convertTestCases(testCases []*TestCase) []*mysql.TestCase {
	// 提前转换类型
	var convertedTestCases []*mysql.TestCase
	for _, tc := range testCases {
		// 进行类型转换
		convertedTestCases = append(convertedTestCases, &mysql.TestCase{
			TID:      tc.TID,
			PID:      tc.PID,
			Input:    tc.Input,
			Expected: tc.Expected,
		})
	}
	return convertedTestCases
}
