package rating

import (
	"sort"
)

var arPairVariants [][2][]int;

func GeneratePairingVariants() {
	var arPlayers []int;
	for i := 0; i < 8; i++ {
		arPlayers = append(arPlayers, i);
	}
	var arVariantsUnpaired = permutations(arPlayers);
	for _, arCheckedVariantUnpaired := range arVariantsUnpaired {

		var arCheckedVariant [2][]int;
		for i, _ := range arCheckedVariantUnpaired {
			iTeam := i % 2;
			arCheckedVariant[iTeam] = append(arCheckedVariant[iTeam], arCheckedVariantUnpaired[i]);
		}
		sort.Ints(arCheckedVariant[0]);
		sort.Ints(arCheckedVariant[1]);

		bVariantExists := VariantExists(arCheckedVariant);
		if (!bVariantExists) {
			arPairVariants = append(arPairVariants, arCheckedVariant);
		}
	}
}

func VariantExists(arCheckedVariant [2][]int) bool {
	for _, arVariant := range arPairVariants {
		if ((ArraysMatch(arCheckedVariant[0], arVariant[0]) && ArraysMatch(arCheckedVariant[1], arVariant[1])) || (ArraysMatch(arCheckedVariant[1], arVariant[0]) && ArraysMatch(arCheckedVariant[0], arVariant[1]))) {
			return true;
		}
	}
	return false;
}

func ArraysMatch(ar1 []int, ar2 []int) bool {
	iLen := len(ar1);
	if (iLen != len(ar2)) {
		return false;
	}
	for i := 0; i < iLen; i++ {
		if (ar1[i] != ar2[i]) {
			return false;
		}
	}
	return true;
}

func permutations(arr []int) [][]int { //https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int){
		if n == 1{
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++{
				helper(arr, n - 1)
				if n % 2 == 1{
					tmp := arr[i]
					arr[i] = arr[n - 1]
					arr[n - 1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n - 1]
					arr[n - 1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
