package utils

import (
	"log"
	"math/rand"
	"reflect"
	"time"
)

func shuffle (toShuffle[]string) []string{
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := range toShuffle {
		randInd := r.Intn(len(toShuffle))
		toShuffle[i], toShuffle[randInd] = toShuffle[randInd], toShuffle[i]
	}
	return toShuffle
}


func ShuffleGroups(prevGroups [][]string) ([]string) {
	shuffled := make([]string, len(prevGroups[0]))
	copy(shuffled, prevGroups[0])

	isSimilar := true
	count := 0
	maxIteration := 1000
	
	for isSimilar && count < maxIteration {
		shuffled = shuffle(shuffled)
		for _, arr := range prevGroups {
			if reflect.DeepEqual(arr, shuffled) {
				count++
				if count == 1000 {
					log.Println("Reached max shuffle iteration")
				}
				break
			}
			isSimilar = false
		}
	}
	
	return shuffled
}

func SortToGroups(all []string) [][]string {
	lenAll := len(all)
	var numGroups int
	if lenAll < 10 {
		numGroups = 3
	} else if lenAll > 9 && lenAll < 15 {
		numGroups = 4
	} else if lenAll > 14 && lenAll < 21 {
		numGroups = 5
	} else {
		numGroups = 6
	}

	out := make([][]string, numGroups)

	for i, name := range all {
		out[i % numGroups] = append(out[i%numGroups], name)
	}

	

	for len(out) < 7 {
		out = append(out, []string{""})
	}

	return out

}