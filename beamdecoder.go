package beam

import (
	"cmp"
	"errors"
	"slices"
)

type BeamDecoder struct {
	beamSize int
}

type sq struct {
	seq   []int
	score float32
}

func NewBeamDecoder(beamSize int) *BeamDecoder {
	return &BeamDecoder{beamSize: beamSize}
}

func (bd *BeamDecoder) Decode(predictedSeq [][]float32, charsList []string) ([]string, []float32, [][]int, error) {

	if len(predictedSeq) == 0 {
		return nil, nil, nil, errors.New("predicted sequence is empty")
	}

	if len(charsList) == 0 {
		return nil, nil, nil, errors.New("character list is empty")
	}

	if len(charsList) < len(predictedSeq) {
		return nil, nil, nil, errors.New("len(charList) < len(predictedSeq)")
	}

	labels := []string{}
	finalProb := []float32{}
	finalLabels := [][]int{}

	sequences := []sq{{[]int{}, float32(0.0)}}

	allSeq := [][]float32{}
	singlePrediction := predictedSeq

	for i := range singlePrediction[0] {
		singleSeq := []float32{}

		for j := range singlePrediction {
			singleSeq = append(singleSeq, singlePrediction[j][i])
		}
		allSeq = append(allSeq, singleSeq)
	}

	for _, row := range allSeq {
		allCandidates := []sq{}
		for i := range sequences {

			seq, score := sequences[i].seq, sequences[i].score
			for j, char := range row {

				var new_seq []int
				new_seq = append(new_seq, seq...)
				new_seq = append(new_seq, j)

				candidate := sq{new_seq, score - char}

				allCandidates = append(allCandidates, candidate)
			}
		}

		slices.SortFunc(allCandidates, func(a, b sq) int {
			return cmp.Compare(a.score, b.score)
		})

		sequences = allCandidates[:bd.beamSize]
	}

	fullPredLabels := [][]int{}
	probs := []float32{}
	for _, i := range sequences {
		predictedLabels := i.seq
		withoutRepeating := []int{}
		currentChar := predictedLabels[0]
		if currentChar != len(charsList)-1 {
			withoutRepeating = append(withoutRepeating, currentChar)
		}
		for _, c := range predictedLabels {
			if (currentChar == c) || (c == len(charsList)-1) {
				if c == len(charsList)-1 {
					currentChar = c
				}
				continue
			}
			withoutRepeating = append(withoutRepeating, c)
			currentChar = c
		}

		fullPredLabels = append(fullPredLabels, withoutRepeating)
		probs = append(probs, i.score)
	}

	for i, label := range fullPredLabels {
		decodedLabel := ""
		for _, j := range label {
			decodedLabel += charsList[j]
		}
		labels = append(labels, decodedLabel)
		finalProb = append(finalProb, probs[i])
		finalLabels = append(finalLabels, label)
	}

	return labels, finalProb, finalLabels, nil
}
