package internal

import "strings"

func arrcmp(src []string, dest []string) ([]string, []string) {
	msrc := make(map[string]byte) //按源数组建索引
	mall := make(map[string]byte) //源+目所有元素建索引
	var set []string              //交集
	//1.源数组建立map
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	//2.目数组中，存不进去，即重复元素，所有存不进去的集合就是并集
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) { //长度变化，即可以存
			l = len(mall)
			_ = l
		} else { //存不了，进并集
			set = append(set, v)
		}
	}
	//3.遍历交集，在并集中找，找到就从并集中删，删完后就是补集（即并-交=所有变化的元素）
	for _, v := range set {
		delete(mall, v)
	}
	//4.此时，mall是补集，所有元素去源中找，找到就是删除的，找不到的必定能在目数组中找到，即新加的
	var added, deleted []string
	for v, m := range mall {
		_ = m
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}

	return added, deleted
}

// 删除上个request append到args的值
func delHttpapiConfArgs(httpapiArgsLen int) {
	apiHubHttpConf.Args = append(apiHubHttpConf.Args[:0], apiHubHttpConf.Args[httpapiArgsLen:]...)
}

func delHttpapiPrivates(httpapiPrivatesLen int) {
	apiHubHttpPrivates.Privates = append(apiHubHttpPrivates.Privates[:0], apiHubHttpPrivates.Privates[httpapiPrivatesLen:]...)
}

func getStringBetweenDoubleQuotationMarks(inputStrings string) (outputString string, outputIndex int) {
	return getStringBetweenSpecifySymbols(inputStrings, "\"", "\"")
}
func getStringBetweenDoubleBrackets(inputStrings string) (outputString string, outputIndex int) {
	return getStringBetweenSpecifySymbols(inputStrings, "{{", "}}")
}

// 获取指定字符中间的字符串，并返回字符串最右的索引值，backIndex = -1表示错误
func getStringBetweenSpecifySymbols(inputStrings string, specifySymbolBefore string, specifySymbolAfter string) (outputString string, outputIndex int) {
	currentIndex := strings.Index(inputStrings, specifySymbolBefore)
	if currentIndex != -1 {
		nextIndex := strings.Index(inputStrings[currentIndex+len(specifySymbolBefore):], specifySymbolAfter)
		if nextIndex != -1 {
			outputString = inputStrings[currentIndex+len(specifySymbolBefore) : currentIndex+len(specifySymbolBefore)+nextIndex]
			outputIndex = nextIndex + len(specifySymbolAfter) + currentIndex
		} else {
			outputString = ""
			outputIndex = -1
		}
	} else {
		outputString = ""
		outputIndex = -1
	}
	return outputString, outputIndex
}
