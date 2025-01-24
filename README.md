# Beam Decoder for NLP

## Usage example:
```
beam_size := 2    
bd := NewBeamDecoder(beam_size)

predictedSeq := [][]float32{
    {0.1, 0.9, 0.0, 0.5},
    {0.8, 0.1, 0.1, 0.3},
    {0.0, 0.2, 0.8, 0.9},
}

charsList := []string{"A", "B", "C", "-"}

labels, probs, chars, _:= bd.Decode(predictedSeq, charsList)

fmt.Println(labels) // "BAC", "BACA"
fmt.Println(probs) //-3.4, -3
fmt.Println(chars) // {1, 0, 2}, {1, 0, 2, 0}
```
