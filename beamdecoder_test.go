package beam

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func readFloatsFromFile(s string) ([]float32, error) {
	// Open the file
	f, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the file
	var floats []float32
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		value, err := strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			return nil, err
		}
		float := float32(value)
		floats = append(floats, float)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return floats, nil
}

func TestNewBeamDecoder(t *testing.T) {
	beamSize := 5
	bd := NewBeamDecoder(beamSize)
	if bd.beamSize != beamSize {
		t.Errorf("Expected beam size %d, got %d", beamSize, bd.beamSize)
	}
}

func TestDecodeSimpleInput(t *testing.T) {
	bd := NewBeamDecoder(2)
	predictedSeq := [][]float32{
		{0.1, 0.9, 0.0, 0.5},
		{0.8, 0.1, 0.1, 0.3},
		{0.0, 0.2, 0.8, 0.9},
	}
	charsList := []string{"A", "B", "C", "-"}

	labels, probs, finalLabels, _ := bd.Decode(predictedSeq, charsList)

	expectedLabels := []string{"BAC", "BACA"}
	expectedProbs := []float32{-3.4, -3}
	expectedFinalLabels := [][]int{{1, 0, 2}, {1, 0, 2, 0}}

	if !reflect.DeepEqual(labels, expectedLabels) {
		t.Errorf("Expected labels %v, got %v", expectedLabels, labels)
	}
	if !reflect.DeepEqual(probs, expectedProbs) {
		t.Errorf("Expected probabilities %v, got %v", expectedProbs, probs)
	}
	if !reflect.DeepEqual(finalLabels, expectedFinalLabels) {
		t.Errorf("Expected final labels %v, got %v", expectedFinalLabels, finalLabels)
	}
}

func TestDecodeEmptyInput(t *testing.T) {
	bd := NewBeamDecoder(2)

	predictedSeq := [][]float32{}
	charsList := []string{"A", "B", "C", "-"}

	labels, probs, finalLabels, _ := bd.Decode(predictedSeq, charsList)

	if len(labels) != 0 || len(probs) != 0 || len(finalLabels) != 0 {
		t.Errorf("Expected empty output, got labels: %v, probs: %v, finalLabels: %v", labels, probs, finalLabels)
	}
}

func TestDecodeRepeatedCharacters(t *testing.T) {
	bd := NewBeamDecoder(3)

	predictedSeq := [][]float32{
		{0.1, 0.5, 0.1, 0.5, 0.1, 0.5, 0.1, 0.5, 0.1, 0.5},
		{0.2, 0.4, 0.2, 0.4, 0.2, 0.4, 0.2, 0.4, 0.2, 0.4},
		{0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3},
		{0.4, 0.2, 0.4, 0.2, 0.4, 0.2, 0.4, 0.2, 0.4, 0.2},
		{0.5, 0.1, 0.5, 0.1, 0.5, 0.1, 0.5, 0.1, 0.5, 0.1},
	}

	charsList := []string{"A", "B", "C", "D", "E", "-"}

	labels, probs, finalLabels, err := bd.Decode(predictedSeq, charsList)

	if err != nil {
		t.Errorf("Error but normal output expected")
	}

	expectedLabels := []string{"EAEAEAEAEA", "EAEAEAEADA", "EAEAEAEAEB"}
	expectedProbs := []float32{-5, -4.9, -4.9}
	expectedFinalLabels := [][]int{{4, 0, 4, 0, 4, 0, 4, 0, 4, 0}, {4, 0, 4, 0, 4, 0, 4, 0, 3, 0}, {4, 0, 4, 0, 4, 0, 4, 0, 4, 1}}

	if !reflect.DeepEqual(labels, expectedLabels) {
		t.Errorf("Expected labels %v, got %v", expectedLabels, labels)
	}
	if !reflect.DeepEqual(probs, expectedProbs) {
		t.Errorf("Expected probabilities %v, got %v", expectedProbs, probs)
	}
	if !reflect.DeepEqual(finalLabels, expectedFinalLabels) {
		t.Errorf("Expected final labels %v, got %v", expectedFinalLabels, finalLabels)
	}
}

func TestDecodeEmptyPredictedSequence(t *testing.T) {
	bd := NewBeamDecoder(2)
	predictedSeq := [][]float32{}
	charsList := []string{"A", "B", "C", "<end>"}

	_, _, _, err := bd.Decode(predictedSeq, charsList)

	if err == nil {
		t.Errorf("Expected an error for empty predicted sequence, got nil")
	}
	if err.Error() != "predicted sequence is empty" {
		t.Errorf("Expected error message 'predicted sequence is empty', got '%v'", err)
	}
}

func TestDecodeDiffLenError(t *testing.T) {
	bd := NewBeamDecoder(3)

	predictedSeq := [][]float32{
		{0.1, 0.5, 0.1, 0.5, 0.1, 0.5, 0.1, 0.5, 0.1, 0.5},
		{0.2, 0.4, 0.2, 0.4, 0.2, 0.4, 0.2, 0.4, 0.2, 0.4},
		{0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3},
		{0.4, 0.2, 0.4, 0.2, 0.4, 0.2, 0.4, 0.2, 0.4, 0.2},
		{0.5, 0.1, 0.5, 0.1, 0.5, 0.1, 0.5, 0.1, 0.5, 0.1},
	}

	charsList := []string{"A", "B", "C", "-"}

	_, _, _, err := bd.Decode(predictedSeq, charsList)

	if err == nil {
		t.Errorf("Expected an error for unequal sequence lengths, got nil")
	}
	if err.Error() != "len(charList) < len(predictedSeq)" {
		t.Errorf("Expected error message 'len(charList) < len(predictedSeq)', got '%v'", err)
	}
}

func TestDecodeRealExample(t *testing.T) {
	ff, err := readFloatsFromFile("inputs.txt")
	if err != nil {
		log.Fatal(err)
	}

	charList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "E", "K", "M", "H", "O", "P", "C", "T", "Y", "X", "-"}
	var scores [][]float32

	bd := NewBeamDecoder(1)

	for i := 0; i < 23; i++ {
		id_start := i * 18
		id_end := id_start + 18
		scores = append(scores, ff[id_start:id_end])
	}
	labels, probs, finalLabels, _ := bd.Decode(scores, charList)

	expectedLabels := []string{"K281AP156"}
	expectedProbs := []float32{-184.5512}
	expectedFinalLabels := [][]int{{13, 2, 8, 1, 10, 17, 1, 5, 6}}

	if !reflect.DeepEqual(labels, expectedLabels) {
		t.Errorf("Expected labels %v, got %v", expectedLabels, labels)
	}
	if !reflect.DeepEqual(probs, expectedProbs) {
		t.Errorf("Expected probabilities %v, got %v", expectedProbs, probs)
	}
	if !reflect.DeepEqual(finalLabels, expectedFinalLabels) {
		t.Errorf("Expected final labels %v, got %v", expectedFinalLabels, finalLabels)
	}
}
