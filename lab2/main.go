package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"regexp"
	"sort"
)

type probs struct {
	FP float64
	FN float64
}

type texts struct {
	text10    []string
	text100   []string
	text1000  []string
	text10000 []string
}

type textt struct {
	vigenere1  []string
	vigenere5  []string
	vigenere10 []string
	affine     []string
	randomm    []string
	reccurent  []string
}

func main() {
	alphabetStr := "абвгдеєжзиіїйклмнопрстуфхцчшщьюя"
	alphabetMap := make(map[string]bool)
	for _, v := range alphabetStr {
		alphabetMap[string(v)] = false
	}

	alphabetstat := make(map[string]float64)
	alphabetFrequencies := make(map[string]float64)
	data := string(readData(alphabetMap))
	for _, letter := range data {
		alphabetstat[string(letter)] += 1
	}

	datalen := float64(len(data))
	for letter, sum := range alphabetstat {
		alphabetFrequencies[letter] = sum / datalen
	}
	monogram := make(map[string]int)
	i := 0
	for _, letter := range alphabetStr {
		monogram[string(letter)] = i
		i++
	}
	monogramStats := make(map[string]float64)
	for _, letter := range data {
		monogramStats[string(letter)]++

	}

	var allBigrams []string
	bigramsStatistics := make(map[string]float64)
	bigramsFrequencies := make(map[string]float64)
	bigramToNumber := make(map[string]int)
	for _, v1 := range alphabetStr {
		for _, v2 := range alphabetStr {
			bigram := string(v1) + string(v2)
			allBigrams = append(allBigrams, bigram)
			bigramsStatistics[bigram] = 0
			bigramsFrequencies[bigram] = 0
		}
	}

	for i := 0; i < int(datalen)-2; i += 4 {
		bigramsStatistics[string(data[i:i+2])+string(data[i+2:i+4])] += 1
	}

	for key, v := range bigramsStatistics {
		bigramsFrequencies[key] = v / (datalen - 1)
	}

	for i, bigram := range allBigrams {
		bigramToNumber[bigram] = i
	}

	// bigramLen := len(allBigrams)

	var monogranEntropy float64
	for _, v := range alphabetStr {
		monogranEntropy += alphabetFrequencies[string(v)] * math.Log2(alphabetFrequencies[string(v)])
	}
	monogranEntropy *= -1

	var bigramEntropy float64
	for _, v := range allBigrams {
		if bigramsFrequencies[v] != 0 {
			bigramEntropy += bigramsFrequencies[v] * math.Log2(bigramsFrequencies[v])
		}
	}

	var coincidenceIndex float64
	for _, v := range alphabetStr {
		coincidenceIndex += alphabetstat[string(v)] * (alphabetstat[string(v)] - 1)
	}
	coincidenceIndex /= datalen
	coincidenceIndex /= (datalen - 1)

	textCount := 10000
	T := texts{
		textGeneretor(10, textCount, 10, data),
		textGeneretor(100, textCount, 100, data),
		textGeneretor(1000, textCount, 100, data),
		textGeneretor(10000, textCount, 1000, data),
	}

	vkey1 := generateKey(len(alphabetStr), 1)
	vkey5 := generateKey(len(alphabetStr), 5)
	vkey10 := generateKey(len(alphabetStr), 10)
	affinekey := generateKey(len(alphabetStr), 2)
	for affinekey[0]%2 != 1 {
		affinekey = generateKey(len(alphabetStr), 2)
	}
	testmono_10 := textt{}
	testmono_100 := textt{}
	testmono_1000 := textt{}
	testmono_10000 := textt{}

	for _, inputText := range T.text10 {
		testmono_10.vigenere1 = append(testmono_10.vigenere1, vigenere(inputText, vkey1, alphabetStr, monogram))
		testmono_10.vigenere5 = append(testmono_10.vigenere5, vigenere(inputText, vkey5, alphabetStr, monogram))
		testmono_10.vigenere10 = append(testmono_10.vigenere10, vigenere(inputText, vkey10, alphabetStr, monogram))
		testmono_10.affine = append(testmono_10.affine, affine(inputText, affinekey, monogram, alphabetStr))
		testmono_10.randomm = append(testmono_10.randomm, uniform(10, alphabetStr))
		testmono_10.reccurent = append(testmono_10.reccurent, reccurent(10, alphabetStr, monogram))
	}
	fmt.Println("Test 10 finished")

	for _, inputText := range T.text100 {
		testmono_100.vigenere1 = append(testmono_100.vigenere1, vigenere(inputText, vkey1, alphabetStr, monogram))
		testmono_100.vigenere5 = append(testmono_100.vigenere5, vigenere(inputText, vkey5, alphabetStr, monogram))
		testmono_100.vigenere10 = append(testmono_100.vigenere10, vigenere(inputText, vkey10, alphabetStr, monogram))
		testmono_100.affine = append(testmono_100.affine, affine(inputText, affinekey, monogram, alphabetStr))
		testmono_100.randomm = append(testmono_100.randomm, uniform(100, alphabetStr))
		testmono_100.reccurent = append(testmono_100.reccurent, reccurent(100, alphabetStr, monogram))
	}
	fmt.Println("Test 100 finished")

	for _, inputText := range T.text1000 {
		testmono_1000.vigenere1 = append(testmono_1000.vigenere1, vigenere(inputText, vkey1, alphabetStr, monogram))
		testmono_1000.vigenere5 = append(testmono_1000.vigenere5, vigenere(inputText, vkey5, alphabetStr, monogram))
		testmono_1000.vigenere10 = append(testmono_1000.vigenere10, vigenere(inputText, vkey10, alphabetStr, monogram))
		testmono_1000.affine = append(testmono_1000.affine, affine(inputText, affinekey, monogram, alphabetStr))
		testmono_1000.randomm = append(testmono_1000.randomm, uniform(1000, alphabetStr))
		testmono_1000.reccurent = append(testmono_1000.reccurent, reccurent(1000, alphabetStr, monogram))
	}
	fmt.Println("Test 1000 finished")

	// for i, inputText := range T.text10000 {
	// 	testmono_10000.vigenere1 = append(testmono_10000.vigenere1, vigenere(inputText, vkey1, alphabetStr, monogram))
	// 	testmono_10000.vigenere5 = append(testmono_10000.vigenere5, vigenere(inputText, vkey5, alphabetStr, monogram))
	// 	testmono_10000.vigenere10 = append(testmono_10000.vigenere10, vigenere(inputText, vkey10, alphabetStr, monogram))
	// 	testmono_10000.affine = append(testmono_10000.affine, affine(inputText, affinekey, monogram, alphabetStr))
	// 	testmono_10000.randomm = append(testmono_10000.randomm, uniform(10000, alphabetStr))
	// 	testmono_10000.reccurent = append(testmono_10000.reccurent, reccurent(10000, alphabetStr, monogram))
	// 	if i%1000 == 0 {
	// 		fmt.Println(i)
	// 	}
	// }
	fmt.Println("Test 10000 finished")

	/////////////////////////BI
	bivkey1 := generateKey(len(allBigrams), 1)
	bivkey5 := generateKey(len(allBigrams), 5)
	bivkey10 := generateKey(len(allBigrams), 10)
	biaffinekey := generateKey(len(allBigrams), 2)
	for biaffinekey[0]%2 != 1 {
		affinekey = generateKey(len(allBigrams), 2)
	}
	testbi_10 := textt{}
	testbi_100 := textt{}
	testbi_1000 := textt{}
	testbi_10000 := textt{}

	for _, inputText := range T.text10 {
		testbi_10.vigenere1 = append(testbi_10.vigenere1, bivigenere(inputText, bivkey1, allBigrams, bigramToNumber))
		testbi_10.vigenere5 = append(testbi_10.vigenere5, bivigenere(inputText, bivkey5, allBigrams, bigramToNumber))
		testbi_10.vigenere10 = append(testbi_10.vigenere10, bivigenere(inputText, bivkey10, allBigrams, bigramToNumber))
		testbi_10.affine = append(testbi_10.affine, biaffine(inputText, affinekey, bigramToNumber, allBigrams))
		testbi_10.randomm = append(testbi_10.randomm, biuniform(10, allBigrams))
		testbi_10.reccurent = append(testbi_10.reccurent, bireccurent(10, allBigrams, bigramToNumber))
	}
	fmt.Println("Test 10 finished")

	for _, inputText := range T.text100 {
		testbi_100.vigenere1 = append(testbi_100.vigenere1, bivigenere(inputText, bivkey1, allBigrams, bigramToNumber))
		testbi_100.vigenere5 = append(testbi_100.vigenere5, bivigenere(inputText, bivkey5, allBigrams, bigramToNumber))
		testbi_100.vigenere10 = append(testbi_100.vigenere10, bivigenere(inputText, bivkey10, allBigrams, bigramToNumber))
		testbi_100.affine = append(testbi_100.affine, biaffine(inputText, affinekey, bigramToNumber, allBigrams))
		testbi_100.randomm = append(testbi_100.randomm, biuniform(10, allBigrams))
		testbi_100.reccurent = append(testbi_100.reccurent, bireccurent(10, allBigrams, bigramToNumber))
	}
	fmt.Println("Test 100 finished")

	for _, inputText := range T.text1000 {
		testbi_1000.vigenere1 = append(testbi_1000.vigenere1, bivigenere(inputText, bivkey1, allBigrams, bigramToNumber))
		testbi_1000.vigenere5 = append(testbi_1000.vigenere5, bivigenere(inputText, bivkey5, allBigrams, bigramToNumber))
		testbi_1000.vigenere10 = append(testbi_1000.vigenere10, bivigenere(inputText, bivkey10, allBigrams, bigramToNumber))
		testbi_1000.affine = append(testbi_1000.affine, biaffine(inputText, affinekey, bigramToNumber, allBigrams))
		testbi_1000.randomm = append(testbi_1000.randomm, biuniform(10, allBigrams))
		testbi_1000.reccurent = append(testbi_1000.reccurent, bireccurent(10, allBigrams, bigramToNumber))
	}
	fmt.Println("Test 1000 finishedd")

	// for i, inputText := range T.text10000 {
	// 	testbi_10000.vigenere1 = append(testbi_10000.vigenere1, bivigenere(inputText, bivkey1, allBigrams, bigramToNumber))
	// 	testbi_10000.vigenere5 = append(testbi_10000.vigenere5, bivigenere(inputText, bivkey5, allBigrams, bigramToNumber))
	// 	testbi_10000.vigenere10 = append(testbi_10000.vigenere10, bivigenere(inputText, bivkey10, allBigrams, bigramToNumber))
	// 	testbi_10000.affine = append(testbi_10000.affine, biaffine(inputText, affinekey, bigramToNumber, allBigrams))
	// 	testbi_10000.randomm = append(testbi_10000.randomm, biuniform(10, allBigrams))
	// 	testbi_10000.reccurent = append(testbi_10000.reccurent, bireccurent(10, allBigrams, bigramToNumber))
	// 	if i%1000 == 0 {
	// 		fmt.Println(i)
	// 	}
	// }
	fmt.Println("Test 10000 finished")

	/////////////CRITER 1.0
	gram := canceledgrams(monogramStats)
	bigram := canceledgrams(bigramsStatistics)
	{
		fmt.Println("CRITER 1.0")
		fmt.Println("mono real")
		L10 := crit1_0(T.text10, gram)
		L100 := crit1_0(T.text100, gram)
		L1000 := crit1_0(T.text1000, gram)
		L10000 := crit1_0(T.text10000, gram)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", L10, L100, L1000, L10000)
		fmt.Println("mono generated")
		Lm10 := crit1_0struct(testmono_10, gram)
		Lm100 := crit1_0struct(testmono_100, gram)
		Lm1000 := crit1_0struct(testmono_1000, gram)
		Lm10000 := crit1_0struct(testmono_10000, gram)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10[0], Lm10[1], Lm10[2], Lm10[3], Lm10[4], Lm10[5])
		fmt.Printf("L100:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm100[0], Lm100[1], Lm100[2], Lm100[3], Lm100[4], Lm100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm1000[0], Lm1000[1], Lm1000[2], Lm1000[3], Lm1000[4], Lm1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10000[0], Lm10000[1], Lm10000[2], Lm10000[3], Lm10000[4], Lm10000[5])
		fmt.Println("bi real")
		Lbi10 := crit1_0bi(T.text10, bigram)
		Lbi100 := crit1_0bi(T.text100, bigram)
		Lbi1000 := crit1_0bi(T.text1000, bigram)
		Lbi10000 := crit1_0bi(T.text10000, bigram)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", Lbi10, Lbi100, Lbi1000, Lbi10000)
		fmt.Println("bi generated")
		LBI10 := bicrit1_0struct(testbi_10, bigram)
		LBI100 := bicrit1_0struct(testbi_100, bigram)
		LBI1000 := bicrit1_0struct(testbi_1000, bigram)
		LBI10000 := bicrit1_0struct(testbi_10000, bigram)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10[0], LBI10[1], LBI10[2], LBI10[3], LBI10[4], LBI10[5])
		fmt.Printf("L100:\n\viginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI100[0], LBI100[1], LBI100[2], LBI100[3], LBI100[4], LBI100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI1000[0], LBI1000[1], LBI1000[2], LBI1000[3], LBI1000[4], LBI1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10000[0], LBI10000[1], LBI10000[2], LBI10000[3], LBI10000[4], LBI10000[5])
	}

	///////////CRITER 1.1
	pg := prohibitedGrams(alphabetstat, 0.9)
	pgbi := prohibitedGrams(bigramsStatistics, 0.75)
	{
		fmt.Println("CRITER 1.1")
		fmt.Println("mono real")
		L10 := crit1_1(T.text10, 1, pg)
		L100 := crit1_1(T.text100, 2, pg)
		L1000 := crit1_1(T.text1000, 3, pg)
		L10000 := crit1_1(T.text10000, 4, pg)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", L10, L100, L1000, L10000)
		fmt.Println("mono generated")
		Lm10 := crit1_1struct(testmono_10, 1, pg)
		Lm100 := crit1_1struct(testmono_100, 2, pg)
		Lm1000 := crit1_1struct(testmono_1000, 3, pg)
		Lm10000 := crit1_1struct(testmono_10000, 4, pg)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10[0], Lm10[1], Lm10[2], Lm10[3], Lm10[4], Lm10[5])
		fmt.Printf("L100:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm100[0], Lm100[1], Lm100[2], Lm100[3], Lm100[4], Lm100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm1000[0], Lm1000[1], Lm1000[2], Lm1000[3], Lm1000[4], Lm1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10000[0], Lm10000[1], Lm10000[2], Lm10000[3], Lm10000[4], Lm10000[5])
		fmt.Println("bi real")
		Lbi10 := crit1_1bi(T.text10, 1, pgbi)
		Lbi100 := crit1_1bi(T.text100, 10, pgbi)
		Lbi1000 := crit1_1bi(T.text1000, 50, pgbi)
		Lbi10000 := crit1_1bi(T.text10000, 100, pgbi)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", Lbi10, Lbi100, Lbi1000, Lbi10000)
		fmt.Println("bi generated")
		LBI10 := bicrit1_1struct(testbi_10, 1, pgbi)
		LBI100 := bicrit1_1struct(testbi_100, 10, pgbi)
		LBI1000 := bicrit1_1struct(testbi_1000, 50, pgbi)
		LBI10000 := bicrit1_1struct(testbi_10000, 100, pgbi)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10[0], LBI10[1], LBI10[2], LBI10[3], LBI10[4], LBI10[5])
		fmt.Printf("L100:\n\viginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI100[0], LBI100[1], LBI100[2], LBI100[3], LBI100[4], LBI100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI1000[0], LBI1000[1], LBI1000[2], LBI1000[3], LBI1000[4], LBI1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10000[0], LBI10000[1], LBI10000[2], LBI10000[3], LBI10000[4], LBI10000[5])
	}

	/////////////CRITER 1.2
	{
		fmt.Println("CRITER 1.2")
		fmt.Println("mono real")
		L10 := crit1_2(T.text10, 0.1, monogram)
		L100 := crit1_2(T.text100, 0.025, monogram)
		L1000 := crit1_2(T.text1000, 0.015, monogram)
		L10000 := crit1_2(T.text10000, 0.008, monogram)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", L10, L100, L1000, L10000)
		fmt.Println("mono generated")
		Lm10 := crit1_2struct(testmono_10, 0.1, monogram)
		Lm100 := crit1_2struct(testmono_100, 0.025, monogram)
		Lm1000 := crit1_2struct(testmono_1000, 0.015, monogram)
		Lm10000 := crit1_2struct(testmono_10000, 0.008, monogram)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10[0], Lm10[1], Lm10[2], Lm10[3], Lm10[4], Lm10[5])
		fmt.Printf("L100:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm100[0], Lm100[1], Lm100[2], Lm100[3], Lm100[4], Lm100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm1000[0], Lm1000[1], Lm1000[2], Lm1000[3], Lm1000[4], Lm1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10000[0], Lm10000[1], Lm10000[2], Lm10000[3], Lm10000[4], Lm10000[5])
		fmt.Println("bi real")
		Lbi10 := bicrit1_2(T.text10, 0.1)
		Lbi100 := bicrit1_2(T.text100, 0.0125)
		Lbi1000 := bicrit1_2(T.text1000, 0.0035)
		Lbi10000 := bicrit1_2(T.text10000, 0.001)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", Lbi10, Lbi100, Lbi1000, Lbi10000)
		fmt.Println("bi generated")
		LBI10 := bicrit1_2struct(testbi_10, 0.1)
		LBI100 := bicrit1_2struct(testbi_100, 0.0125)
		LBI1000 := bicrit1_2struct(testbi_1000, 0.0035)
		LBI10000 := bicrit1_2struct(testbi_10000, 0.001)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10[0], LBI10[1], LBI10[2], LBI10[3], LBI10[4], LBI10[5])
		fmt.Printf("L100:\n\viginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI100[0], LBI100[1], LBI100[2], LBI100[3], LBI100[4], LBI100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI1000[0], LBI1000[1], LBI1000[2], LBI1000[3], LBI1000[4], LBI1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10000[0], LBI10000[1], LBI10000[2], LBI10000[3], LBI10000[4], LBI10000[5])
	}

	/////////////CRITER 1.3
	{
		fmt.Println("CRITER 1.3")
		fmt.Println("mono real")
		L10 := crit1_3(T.text10, 0.05, gram)
		L100 := crit1_3(T.text100, 0.025, gram)
		L1000 := crit1_3(T.text1000, 0.0125, gram)
		L10000 := crit1_3(T.text10000, 0.01, gram)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", L10, L100, L1000, L10000)
		fmt.Println("mono generated")
		Lm10 := crit1_3struct(testmono_10, 0.05, gram)
		Lm100 := crit1_3struct(testmono_100, 0.025, gram)
		Lm1000 := crit1_3struct(testmono_1000, 0.0125, gram)
		Lm10000 := crit1_3struct(testmono_10000, 0.01, gram)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10[0], Lm10[1], Lm10[2], Lm10[3], Lm10[4], Lm10[5])
		fmt.Printf("L100:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm100[0], Lm100[1], Lm100[2], Lm100[3], Lm100[4], Lm100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm1000[0], Lm1000[1], Lm1000[2], Lm1000[3], Lm1000[4], Lm1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10000[0], Lm10000[1], Lm10000[2], Lm10000[3], Lm10000[4], Lm10000[5])
		fmt.Println("bi real")
		Lbi10 := bicrit1_3(T.text10, 0.05)
		Lbi100 := bicrit1_3(T.text100, 0.0125)
		Lbi1000 := bicrit1_3(T.text1000, 0.01)
		Lbi10000 := bicrit1_3(T.text10000, 0.004)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", Lbi10, Lbi100, Lbi1000, Lbi10000)
		fmt.Println("bi generated")
		LBI10 := bicrit1_3struct(testbi_10, 0.03)
		LBI100 := bicrit1_3struct(testbi_100, 0.0125)
		LBI1000 := bicrit1_3struct(testbi_1000, 0.008)
		LBI10000 := bicrit1_3struct(testbi_10000, 0.004)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10[0], LBI10[1], LBI10[2], LBI10[3], LBI10[4], LBI10[5])
		fmt.Printf("L100:\n\viginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI100[0], LBI100[1], LBI100[2], LBI100[3], LBI100[4], LBI100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI1000[0], LBI1000[1], LBI1000[2], LBI1000[3], LBI1000[4], LBI1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10000[0], LBI10000[1], LBI10000[2], LBI10000[3], LBI10000[4], LBI10000[5])
	}

	/////////////CRITER 3.0
	{
		fmt.Println("CRITER 3.0")
		fmt.Println("mono real")
		L10 := crit3_0(T.text10, 2.95, monogranEntropy)
		L100 := crit3_0(T.text100, 1, monogranEntropy)
		L1000 := crit3_0(T.text1000, 0.16, monogranEntropy)
		L10000 := crit3_0(T.text10000, 0.05, monogranEntropy)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", L10, L100, L1000, L10000)
		fmt.Println("mono generated")
		Lm10 := crit3_0struct(testmono_10, 0.1, monogranEntropy)
		Lm100 := crit3_0struct(testmono_100, 0.5, monogranEntropy)
		Lm1000 := crit3_0struct(testmono_1000, 0.08, monogranEntropy)
		Lm10000 := crit3_0struct(testmono_10000, 0.025, monogranEntropy)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10[0], Lm10[1], Lm10[2], Lm10[3], Lm10[4], Lm10[5])
		fmt.Printf("L100:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm100[0], Lm100[1], Lm100[2], Lm100[3], Lm100[4], Lm100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm1000[0], Lm1000[1], Lm1000[2], Lm1000[3], Lm1000[4], Lm1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10000[0], Lm10000[1], Lm10000[2], Lm10000[3], Lm10000[4], Lm10000[5])
		fmt.Println("bi real")
		Lbi10 := bicrit3_0(T.text10, 2, bigramEntropy)
		Lbi100 := bicrit3_0(T.text100, 2, bigramEntropy)
		Lbi1000 := bicrit3_0(T.text1000, 0.275, bigramEntropy)
		Lbi10000 := bicrit3_0(T.text10000, 0.075, bigramEntropy)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", Lbi10, Lbi100, Lbi1000, Lbi10000)
		fmt.Println("bi generated")
		LBI10 := bicrit3_0struct(testbi_10, 2, bigramEntropy)
		LBI100 := bicrit3_0struct(testbi_100, 2, bigramEntropy)
		LBI1000 := bicrit3_0struct(testbi_1000, 0.275, bigramEntropy)
		LBI10000 := bicrit3_0struct(testbi_10000, 0.075, bigramEntropy)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10[0], LBI10[1], LBI10[2], LBI10[3], LBI10[4], LBI10[5])
		fmt.Printf("L100:\n\viginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI100[0], LBI100[1], LBI100[2], LBI100[3], LBI100[4], LBI100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI1000[0], LBI1000[1], LBI1000[2], LBI1000[3], LBI1000[4], LBI1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10000[0], LBI10000[1], LBI10000[2], LBI10000[3], LBI10000[4], LBI10000[5])
	}
	/////////////CRITER 5.1
	{
		fmt.Println("CRITER 5.1")
		mono := getCommonGrams(alphabetstat, 10)
		bi50 := getCommonGrams(alphabetstat, 50)
		bi100 := getCommonGrams(alphabetstat, 100)
		bi200 := getCommonGrams(alphabetstat, 200)
		fmt.Println("mono real")
		L10 := crit5_1(T.text10, 7, mono)
		L100 := crit5_1(T.text100, 0, mono)
		L1000 := crit5_1(T.text1000, 0, mono)
		L10000 := crit5_1(T.text10000, 0, mono)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", L10, L100, L1000, L10000)
		fmt.Println("mono generated")
		Lm10 := crit5_1struct(testmono_10, 7, mono)
		Lm100 := crit5_1struct(testmono_100, 0, mono)
		Lm1000 := crit5_1struct(testmono_1000, 0, mono)
		Lm10000 := crit5_1struct(testmono_10000, 0, mono)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10[0], Lm10[1], Lm10[2], Lm10[3], Lm10[4], Lm10[5])
		fmt.Printf("L100:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm100[0], Lm100[1], Lm100[2], Lm100[3], Lm100[4], Lm100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm1000[0], Lm1000[1], Lm1000[2], Lm1000[3], Lm1000[4], Lm1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10000[0], Lm10000[1], Lm10000[2], Lm10000[3], Lm10000[4], Lm10000[5])
		fmt.Println("bi real")
		Lbi10 := bicrit5_1(T.text10, 47, bi50)
		Lbi100 := bicrit5_1(T.text100, 30, bi50)
		Lbi1000 := bicrit5_1(T.text1000, 10, bi100)
		Lbi10000 := bicrit5_1(T.text10000, 0, bi200)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", Lbi10, Lbi100, Lbi1000, Lbi10000)
		fmt.Println("bi generated")
		LBI10 := bicrit5_1struct(testbi_10, 47, bi50)
		LBI100 := bicrit5_1struct(testbi_100, 30, bi50)
		LBI1000 := bicrit5_1struct(testbi_1000, 10, bi100)
		LBI10000 := bicrit5_1struct(testbi_10000, 0, bi200)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10[0], LBI10[1], LBI10[2], LBI10[3], LBI10[4], LBI10[5])
		fmt.Printf("L100:\n\viginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI100[0], LBI100[1], LBI100[2], LBI100[3], LBI100[4], LBI100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI1000[0], LBI1000[1], LBI1000[2], LBI1000[3], LBI1000[4], LBI1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10000[0], LBI10000[1], LBI10000[2], LBI10000[3], LBI10000[4], LBI10000[5])
	}

	/////////////struct CRITER
	{
		fmt.Println("CRITER struct")
		fmt.Println("mono real")
		L10 := criteria(T.text10, 0.001, alphabetStr)
		L100 := criteria(T.text100, 0.025, alphabetStr)
		L1000 := criteria(T.text1000, 0.15, alphabetStr)
		L10000 := criteria(T.text10000, 0.25, alphabetStr)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", L10, L100, L1000, L10000)
		fmt.Println("mono generated")
		Lm10 := criteriastruct(testmono_10, 0.001, alphabetStr)
		Lm100 := criteriastruct(testmono_100, 0.025, alphabetStr)
		Lm1000 := criteriastruct(testmono_1000, 0.15, alphabetStr)
		Lm10000 := criteriastruct(testmono_10000, 0.25, alphabetStr)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10[0], Lm10[1], Lm10[2], Lm10[3], Lm10[4], Lm10[5])
		fmt.Printf("L100:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm100[0], Lm100[1], Lm100[2], Lm100[3], Lm100[4], Lm100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm1000[0], Lm1000[1], Lm1000[2], Lm1000[3], Lm1000[4], Lm1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", Lm10000[0], Lm10000[1], Lm10000[2], Lm10000[3], Lm10000[4], Lm10000[5])
		fmt.Println("bi real")
		Lbi10 := criteria(T.text10, 0.001, alphabetStr)
		Lbi100 := criteria(T.text100, 0.025, alphabetStr)
		Lbi1000 := criteria(T.text1000, 0.15, alphabetStr)
		Lbi10000 := criteria(T.text10000, 0.25, alphabetStr)
		fmt.Printf("L10 : %v\nL100 : %v\nL1000 : %v\nL10000 : %v\n", Lbi10, Lbi100, Lbi1000, Lbi10000)
		fmt.Println("bi generated")
		LBI10 := criteriastruct(testbi_10, 0.001, alphabetStr)
		LBI100 := criteriastruct(testbi_100, 0.025, alphabetStr)
		LBI1000 := criteriastruct(testbi_1000, 0.15, alphabetStr)
		LBI10000 := criteriastruct(testbi_10000, 0.25, alphabetStr)
		fmt.Printf("L10:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10[0], LBI10[1], LBI10[2], LBI10[3], LBI10[4], LBI10[5])
		fmt.Printf("L100:\n\viginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI100[0], LBI100[1], LBI100[2], LBI100[3], LBI100[4], LBI100[5])
		fmt.Printf("L1000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI1000[0], LBI1000[1], LBI1000[2], LBI1000[3], LBI1000[4], LBI1000[5])
		fmt.Printf("L10000:\n\tviginere 1 :%v\n\tviginere 5 :%v\n\tviginere 10 :%v\n\taffine :%v\n\trandom :%v\n\treccurent :%v\n", LBI10000[0], LBI10000[1], LBI10000[2], LBI10000[3], LBI10000[4], LBI10000[5])
	}
}

func criteria(T []string, limit float64, alp string) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (structCriteria(t, limit, alp))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func criteriastruct(T textt, limit float64, alp string) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (structCriteria(t, limit, alp))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (structCriteria(t, limit, alp))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (structCriteria(t, limit, alp))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (structCriteria(t, limit, alp))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.randomm {
		buf1, buf2 := (structCriteria(t, limit, alp))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (structCriteria(t, limit, alp))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit5_1(T []string, limit float64, commonMonoGrams []string) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (criteria5_1(t, limit, commonMonoGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func bicrit5_1(T []string, limit float64, commonBiGrams []string) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (bicriteria5_1(t, limit, commonBiGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func bicrit5_1struct(T textt, limit float64, commonBiGrams []string) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (bicriteria5_1(t, limit, commonBiGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (bicriteria5_1(t, limit, commonBiGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (bicriteria5_1(t, limit, commonBiGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (bicriteria5_1(t, limit, commonBiGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.randomm {
		buf1, buf2 := (bicriteria5_1(t, limit, commonBiGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (bicriteria5_1(t, limit, commonBiGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit5_1struct(T textt, limit float64, commonMonoGrams []string) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (criteria5_1(t, limit, commonMonoGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (criteria5_1(t, limit, commonMonoGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (criteria5_1(t, limit, commonMonoGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (criteria5_1(t, limit, commonMonoGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.randomm {
		buf1, buf2 := (criteria5_1(t, limit, commonMonoGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (criteria5_1(t, limit, commonMonoGrams))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit3_0(T []string, limit float64, entropy float64) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (criteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func bicrit3_0(T []string, limit float64, entropy float64) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (bicriteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func bicrit3_0struct(T textt, limit float64, entropy float64) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (bicriteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (bicriteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (bicriteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (bicriteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.randomm {
		buf1, buf2 := (bicriteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (bicriteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit3_0struct(T textt, limit float64, entropy float64) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (criteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (criteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (criteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (criteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	for _, t := range T.randomm {
		buf1, buf2 := (criteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (criteria3_0(t, limit, entropy))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit1_3(T []string, limit float64, mono map[string]bool) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (criteria1_3(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func bicrit1_3(T []string, limit float64) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (bicriteria1_3(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func bicrit1_3struct(T textt, limit float64) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (bicriteria1_3(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (bicriteria1_3(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (bicriteria1_3(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (bicriteria1_3(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.randomm {
		buf1, buf2 := (bicriteria1_3(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (bicriteria1_3(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit1_3struct(T textt, limit float64, mono map[string]bool) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	var h1, h0 = 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (criteria1_3(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))
	for _, t := range T.vigenere5 {
		buf1, buf2 := (criteria1_3(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))
	for _, t := range T.vigenere10 {
		buf1, buf2 := (criteria1_3(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))
	for _, t := range T.affine {
		buf1, buf2 := (criteria1_3(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))
	for _, t := range T.randomm {
		buf1, buf2 := (criteria1_3(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))
	for _, t := range T.reccurent {
		buf1, buf2 := (criteria1_3(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit1_2(T []string, limit float64, mono map[string]int) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (criteria1_2(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func bicrit1_2(T []string, limit float64) []probs {
	l := float64(len(T))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T {
		buf1, buf2 := (bicriteria1_2(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func bicrit1_2struct(T textt, limit float64) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (bicriteria1_2(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (bicriteria1_2(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (bicriteria1_2(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (bicriteria1_2(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.randomm {
		buf1, buf2 := (bicriteria1_2(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (bicriteria1_2(t, limit))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit1_2struct(T textt, limit float64, mono map[string]int) []probs {
	l := float64(len(T.vigenere1))
	var res []probs
	h1, h0 := 0.0, 0.0
	for _, t := range T.vigenere1 {
		buf1, buf2 := (criteria1_2(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (criteria1_2(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (criteria1_2(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (criteria1_2(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.randomm {
		buf1, buf2 := (criteria1_2(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (criteria1_2(t, limit, mono))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func crit1_1bi(T []string, prohibited int, pg map[string]float64) []probs {
	var res []probs
	h1, h0 := 0.0, 0.0
	l := float64(len(T))
	for _, t := range T {
		buf1, buf2 := (bicriteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func crit1_0bi(T []string, pg map[string]bool) []probs {
	var res []probs
	h1, h0 := 0.0, 0.0
	l := float64(len(T))
	for _, t := range T {
		buf1, buf2 := (bicriteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func bicrit1_0struct(T textt, pg map[string]bool) []probs {
	var res []probs
	h1, h0 := 0.0, 0.0
	l := float64(len(T.vigenere1))
	for _, t := range T.vigenere1 {
		buf1, buf2 := (bicriteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere5))
	for _, t := range T.vigenere5 {
		buf1, buf2 := (bicriteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	l = float64(len(T.vigenere10))
	for _, t := range T.vigenere10 {
		buf1, buf2 := (bicriteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	l = float64(len(T.affine))
	for _, t := range T.affine {
		buf1, buf2 := (bicriteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	l = float64(len(T.randomm))
	for _, t := range T.randomm {
		buf1, buf2 := (bicriteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	l = float64(len(T.reccurent))
	for _, t := range T.reccurent {
		buf1, buf2 := (bicriteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func bicrit1_1struct(T textt, prohibited int, pg map[string]float64) []probs {
	var res []probs
	h1, h0 := 0.0, 0.0
	l := float64(len(T.vigenere1))
	for _, t := range T.vigenere1 {
		buf1, buf2 := (bicriteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere5))
	for _, t := range T.vigenere5 {
		buf1, buf2 := (bicriteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere10))
	for _, t := range T.vigenere10 {
		buf1, buf2 := (bicriteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	h1, h0 = 0.0, 0.0
	l = float64(len(T.affine))
	for _, t := range T.affine {
		buf1, buf2 := (bicriteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	h1, h0 = 0.0, 0.0
	l = float64(len(T.randomm))

	for _, t := range T.randomm {
		buf1, buf2 := bicriteria1_1(t, prohibited, pg)
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	h1, h0 = 0.0, 0.0
	l = float64(len(T.reccurent))

	for _, t := range T.reccurent {
		buf1, buf2 := (bicriteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func crit1_1(T []string, prohibited int, pg map[string]float64) []probs {
	var res []probs
	h1, h0 := 0.0, 0.0
	l := float64(len(T))
	for _, t := range T {
		buf1, buf2 := (criteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func crit1_0(T []string, pg map[string]bool) []probs {
	var res []probs
	h1, h0 := 0.0, 0.0
	l := float64(len(T))
	for _, t := range T {
		buf1, buf2 := (criteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func crit1_0struct(T textt, pg map[string]bool) []probs {
	var res []probs
	h1, h0 := 0.0, 0.0
	l := float64(len(T.vigenere1))
	for _, t := range T.vigenere1 {
		buf1, buf2 := (criteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))
	for _, t := range T.vigenere5 {
		buf1, buf2 := (criteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))

	for _, t := range T.vigenere10 {
		buf1, buf2 := (criteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))

	for _, t := range T.affine {
		buf1, buf2 := (criteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))

	for _, t := range T.randomm {
		buf1, buf2 := (criteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0
	l = float64(len(T.vigenere1))

	for _, t := range T.reccurent {
		buf1, buf2 := (criteria1_0(t, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})

	return res
}

func crit1_1struct(T textt, prohibited int, pg map[string]float64) []probs {
	var res []probs
	h1, h0 := 0.0, 0.0
	l := float64(len(T.vigenere1))
	for _, t := range T.vigenere1 {
		buf1, buf2 := (criteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere5 {
		buf1, buf2 := (criteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.vigenere10 {
		buf1, buf2 := (criteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.affine {
		buf1, buf2 := (criteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.randomm {
		buf1, buf2 := (criteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	for _, t := range T.reccurent {
		buf1, buf2 := (criteria1_1(t, prohibited, pg))
		h1 += float64(buf1)
		h0 += float64(buf2)
	}
	res = append(res, probs{h1 / l, h0 / l})
	h1, h0 = 0.0, 0.0

	return res
}

func readData(alp map[string]bool) []byte {
	file, _ := ioutil.ReadFile("./sourcetext.txt")

	file = bytes.ReplaceAll(bytes.ToLower(file), []byte("ґ"), []byte("г"))

	var formated []byte
	for _, v := range string(file) {
		if _, ok := alp[string(v)]; ok {
			formated = append(formated, []byte(string(v))...)
		}
	}

	return formated
}

func textGeneretor(size, count, step int, formText string) []string {
	i := 0
	var text []string
	for j := 0; j < count; j++ {
		text = append(text, formText[i:(i+size)])
		i += step
		i %= len(formText) - size
	}
	return text
}

func vigenere(text string, key []int, alp string, mono map[string]int) string {
	cipher := ""
	for i := 0; i < len(text); i += 2 {
		keyIndex := i % len(key)
		st := 2 * (mono[string(text[i:i+2])] + key[keyIndex]) % len(alp)
		cipher += string(alp[st : st+2])
	}

	return cipher
}

func bivigenere(text string, key []int, allbiggram []string, bigramToNumber map[string]int) string {
	cipher := ""
	for i := 0; i < len(text)-4; i += 4 {
		keyIndex := (i / 2) % len(key)
		cipher += allbiggram[(bigramToNumber[string(text[i:i+4])]+key[keyIndex])%len(allbiggram)]
	}

	return cipher
}

func affine(text string, key []int, mono map[string]int, alp string) string {
	a, b := key[0], key[1]
	cipher := ""
	for _, letter := range text {
		st := 2 * (a*mono[string(letter)] + b) % len(alp)
		cipher += string(alp[st : st+2])
	}
	return cipher
}

func biaffine(text string, key []int, bigramToNumber map[string]int, allbiggram []string) string {
	a, b := key[0], key[1]
	cipher := ""

	r, _ := regexp.Compile("..")

	for _, bi := range r.FindAllString(text, -1) {
		index := (a*bigramToNumber[bi] + b) % len(allbiggram)
		cipher += allbiggram[index]
	}

	return cipher
}

func generateKey(high, size int) []int {
	var key []int
	for i := 0; i < size; i++ {
		key = append(key, rand.Intn(high))
	}

	return key
}

func uniform(size int, alp string) string {
	var elems []int
	for i := 0; i < size; i++ {
		elems = append(elems, rand.Intn(len(alp)/2))
	}
	res := ""
	for _, v := range elems {
		res += string(alp[v*2 : 2*v+2])
	}

	return res
}

func biuniform(size int, allbigram []string) string {
	var elems []int
	for i := 0; i < size; i++ {
		elems = append(elems, rand.Intn(len(allbigram)))
	}
	res := ""
	for _, v := range elems {
		res += allbigram[v]
	}

	return res
}

func criteria1_0(text string, grams map[string]bool) (int, int) {
	for _, gram := range text {
		if _, ok := grams[string(gram)]; ok {
			return 1, 0
		}
		return 0, 0

	}

	return 0, 0
}

func bicriteria1_0(text string, bigrams map[string]bool) (int, int) {
	for i := 0; i < len(text)-4; i += 2 {
		if _, ok := bigrams[(text[i : i+4])]; ok {
			return 1, 0
		} else {
			return 0, 1
		}
	}

	return 0, 0
}

func criteria1_1(text string, prohibited int, prohibitedGrams map[string]float64) (int, int) {
	counter := 0

	for _, elem := range text {
		if _, ok := prohibitedGrams[string(elem)]; ok {
			counter++
		}
	}

	if counter >= prohibited {
		return 0, 1
	}
	return 1, 0
}

func bicriteria1_1(text string, prohibited int, prohibitedGrams map[string]float64) (int, int) {
	var bitext []string
	for i := 0; i < len(text)-2; i += 2 {
		bitext = append(bitext, text[i:i+4])
	}

	counter := 0
	for _, elem := range bitext {
		if _, ok := prohibitedGrams[string(elem)]; ok {
			counter++
		}
	}

	if counter >= prohibited {
		return 0, 1
	}

	return 1, 0
}

func criteria1_2(text string, limit float64, mono map[string]int) (int, int) {
	tlen := len(text)
	pFrequencies := make(map[string]float64)

	for key, _ := range mono {
		pFrequencies[key] = 0
	}

	for _, letter := range text {
		if _, ok := pFrequencies[string(letter)]; ok {
			pFrequencies[string(letter)]++
		}
	}

	for _, value := range pFrequencies {
		if value/float64(tlen) >= limit {
			return 1, 0
		}
	}

	return 0, 1
}

func bicriteria1_2(text string, limit float64) (int, int) {
	tlen := (len(text))

	pFrequencies := make(map[string]float64)

	for i := 0; i < len(text)-2; i += 2 {
		bigram := text[i : i+4]

		if _, ok := pFrequencies[bigram]; ok {
			pFrequencies[bigram]++
		}
	}

	for _, value := range pFrequencies {
		if value/(float64(tlen)) >= limit {
			return 1, 0
		}
	}

	return 0, 1
}

func criteria1_3(text string, limit float64, mono map[string]bool) (int, int) {
	pFrequencies := make(map[string]float64)

	for key, _ := range mono {
		pFrequencies[key] = 0
	}

	for _, letter := range text {
		if _, ok := pFrequencies[string(letter)]; ok {
			pFrequencies[string(letter)]++
		}
	}

	sum := 0.0

	for _, value := range pFrequencies {
		sum += value
	}

	if sum/float64(len(text)) > limit {
		return 1, 0
	}

	return 0, 1
}

func bicriteria1_3(text string, limit float64) (int, int) {
	pFrequencies := make(map[string]float64)

	for i := 0; i < len(text)-2; i += 2 {
		bigram := text[i : i+4]
		pFrequencies[bigram]++
	}

	for _, value := range pFrequencies {
		if value/float64(len(text)) >= (limit) {
			return 1, 0
		}
	}

	return 0, 1
}

func criteria3_0(text string, limit float64, monoEntropy float64) (int, int) {
	textEntropy := entropy(distribution(text), 1)

	if math.Abs(monoEntropy-textEntropy) > limit {
		return 1, 0
	}

	return 0, 1
}

func bicriteria3_0(text string, limit float64, Entropy float64) (int, int) {
	textEntropy := entropy(bidistribution(text), 2)

	if math.Abs(Entropy-textEntropy) > limit {
		return 1, 0
	}
	return 0, 1
}

func criteria5_1(text string, limit float64, commonMonoGrams []string) (int, int) {
	textCommonMonoGrams := make(map[string]float64)
	for _, letter := range commonMonoGrams {
		textCommonMonoGrams[letter] = 0
	}

	for _, letter := range commonMonoGrams {
		textCommonMonoGrams[letter]++
	}

	zeroValue := 0
	for _, v := range textCommonMonoGrams {
		if v == 0.0 {
			zeroValue++
		}
	}

	if zeroValue > int(limit) {
		return 1, 0
	}
	return 0, 1
}

func bicriteria5_1(text string, limit float64, commonBiGrams []string) (int, int) {
	textCommonBiGrams := make(map[string]float64)
	commonBiGramsmap := make(map[string]bool)

	for _, bi := range commonBiGrams {
		textCommonBiGrams[bi] = 0
		commonBiGramsmap[bi] = false
	}

	for i := 0; i < len(text)-2; i += 2 {
		bigram := text[i : i+4]
		if _, ok := commonBiGramsmap[bigram]; ok {
			textCommonBiGrams[bigram]++
		}
	}

	zeroValue := 0
	for _, v := range textCommonBiGrams {
		if v == 0.0 {
			zeroValue++
		}
	}

	if zeroValue > int(limit) {
		return 1, 0
	}

	return 0, 1
}

func structCriteria(text string, limit float64, alp string) (int, int) {
	randomText := uniform(len(text), alp)

	randomCoef := float64(len(text)) / float64(len(zipData([]byte(randomText))))
	textCoef := float64(len(text)) / float64(len(zipData([]byte(text))))

	if math.Abs(randomCoef-textCoef) > limit {
		return 1, 0
	}
	return 0, 1
}

func zipData(data []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(data)
	w.Close()

	return in.Bytes()
}

func getCommonGrams(grams map[string]float64, count int) []string {
	type t struct {
		letter string
		stat   float64
	}
	var ts []t
	for key, value := range grams {
		ts = append(ts, t{key, value})
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].stat > ts[j].stat
	})

	return func() []string {
		var res []string
		for _, v := range ts {
			res = append(res, v.letter)
		}
		return res
	}()
}

func distribution(text string) map[string]float64 {
	stats := make(map[string]float64)

	for _, letter := range text {
		stats[string(letter)] += 1
	}

	for key, _ := range stats {
		stats[key] /= float64(len(text))
	}

	return stats
}

func bidistribution(text string) map[string]float64 {
	stats := make(map[string]float64)

	for i := 0; i < len(text); i += 2 {
		bigram := text[i : i+2]
		stats[string(bigram)] += 1
	}

	for key, _ := range stats {
		stats[key] /= float64(len(text))
	}

	return stats
}

func entropy(freq map[string]float64, l float64) (entropy float64) {
	for _, value := range freq {
		entropy -= value * math.Log2(value) / l
	}

	return
}

func prohibitedGrams(grams map[string]float64, q float64) map[string]float64 {
	type t struct {
		letter string
		number float64
	}
	var ts []t
	for letter, num := range grams {
		ts = append(ts, t{letter: letter, number: num})
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].number > ts[j].number
	})

	newts := make(map[string]float64)
	ts = ts[int(math.Floor(q*float64(len(ts)))):]

	for _, t := range ts {
		newts[t.letter] = t.number
	}

	return newts
}

func reccurent(size int, alp string, monogram_to_number map[string]int) string {
	S := generateKey(len(alp)/2, 2)

	recurrent_sequence := string(alp[2*S[0]:2*S[0]+2]) + string(alp[2*S[1]:2*S[1]+2])

	for i := 4; i < size; i += 2 {
		prev_letter_index := monogram_to_number[string(recurrent_sequence[i-2:i])]
		prev_prev_letter_index := monogram_to_number[string(recurrent_sequence[i-4:i-2])]
		next_letter_index := (2*prev_letter_index + 2*prev_prev_letter_index) % len(alp)

		recurrent_sequence += string(alp[next_letter_index : next_letter_index+2])
	}

	return recurrent_sequence
}

func bireccurent(size int, allbigram []string, bigramtonumber map[string]int) string {
	S := generateKey(len(allbigram)/2, 2)

	recurrent_sequence := allbigram[S[0]] + allbigram[S[1]]

	for i := 4; i < size; i += 2 {
		prev_letter_index := bigramtonumber[recurrent_sequence[i-2:]]
		prev_prev_letter_index := bigramtonumber[recurrent_sequence[i-4:i-2]]
		next_letter_index := (prev_letter_index + prev_prev_letter_index) % len(allbigram)

		recurrent_sequence += allbigram[next_letter_index]
	}

	return recurrent_sequence
}

func canceledgrams(g map[string]float64) map[string]bool {
	var gg []string
	for gram, _ := range g {
		gg = append(gg, gram)
	}

	sort.Slice(gg, func(i, j int) bool {
		return gg[i] > gg[j]
	})
	fmt.Println("---------------", gg)
	res := make(map[string]bool)
	for i := 0; i < 10; i++ {
		res[gg[i]] = false
	}

	return res
}
