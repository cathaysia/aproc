package lib

// https://blog.csdn.net/yiweiyi329/article/details/101030510

// 求并集
func Union(slice1, slice2 []int) []int {
	map1 := make(map[int]int)
	for _, v := range slice1 {
		map1[v]++
	}

	for _, v := range slice2 {
		times := map1[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}

	return slice1
}

// 求交集
func Intersect(slice1, slice2 []int) []int {
	map1 := make(map[int]int)
	result := make([]int, 0)

	for _, v := range slice1 {
		map1[v]++
	}

	for _, v := range slice2 {
		times := map1[v]
		if times == 1 {
			result = append(result, v)
		}
	}

	return result
}

// 求差集 slice1-并集
func Difference(slice1, slice2 []int) []int {
	map1 := make(map[int]int)
	result := make([]int, 0)
	inter := Intersect(slice1, slice2)

	for _, v := range inter {
		map1[v]++
	}

	for _, value := range slice1 {
		times := map1[value]
		if times == 0 {
			result = append(result, value)
		}
	}

	return result
}

// https://www.jianshu.com/p/80f5f5173fca

func StringSliceEqualBCE(str1, str2 []string) bool {
	if len(str1) != len(str2) {
		return false
	}

	if (str1 == nil) != (str2 == nil) {
		return false
	}

	str2 = str2[:len(str1)]
	for i, v := range str1 {
		if v != str2[i] {
			return false
		}
	}

	return true
}

func IntSliceEqualBCE(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	if (slice1 == nil) != (slice2 == nil) {
		return false
	}

	slice2 = slice2[:len(slice1)]
	for i, v := range slice1 {
		if v != slice2[i] {
			return false
		}
	}

	return true
}
