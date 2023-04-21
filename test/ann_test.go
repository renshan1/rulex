package test

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

// go test -timeout 30s -run ^Test_ANN_wine_demo github.com/i4de/rulex/test -v -count=1

func Test_ANN_wine_demo(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	data, err := load("./data/wine.data")
	if err != nil {
		panic(err)
	}

	for i := range data {
		deep.Standardize(data[i].Input)
	}
	data.Shuffle()

	fmt.Printf("have %d entries\n", len(data))

	neural := deep.NewNeural(&deep.Config{
		Inputs:     len(data[0].Input),
		Layout:     []int{2, 3},
		Activation: deep.ActivationSigmoid,
		Mode:       deep.ModeMultiClass,
		Weight:     deep.NewNormal(1, 0),
		Bias:       true,
	})
	// t.Log(neural.String())
	//trainer := training.NewTrainer(training.NewSGD(0.005, 0.5, 1e-6, true), 50)
	//trainer := training.NewBatchTrainer(training.NewSGD(0.005, 0.1, 0, true), 50, 300, 16)
	//trainer := training.NewTrainer(training.NewAdam(0.1, 0, 0, 0), 50)
	trainer := training.NewBatchTrainer(training.NewAdam(0.1, 0, 0, 0), 50, len(data)/2, 12)
	//data, heldout := data.Split(0.5)
	trainer.Train(neural, data, data, 5000)
	testData1 := []float64{13.05, 1.65, 2.55, 18, 98, 2.45, 2.43, .29, 1.44, 4.25, 1.12, 2.51, 1105}
	testData2 := []float64{12.16, 1.61, 2.31, 22.8, 90, 1.78, 1.69, .43, 1.56, 2.45, 1.33, 2.26, 495}
	testData3 := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	result1 := neural.Predict(testData1)
	result2 := neural.Predict(testData2)
	result3 := neural.Predict(testData3)
	t.Log(result1)
	t.Log(result2)
	t.Log(result3)
}

func load(path string) (training.Examples, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(bufio.NewReader(f))

	var examples training.Examples
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		examples = append(examples, toExample(record))
	}

	return examples, nil
}

func toExample(in []string) training.Example {
	res, err := strconv.ParseFloat(in[0], 64)
	if err != nil {
		panic(err)
	}
	resEncoded := onehot(3, res)
	var features []float64
	for i := 1; i < len(in); i++ {
		res, err := strconv.ParseFloat(in[i], 64)
		if err != nil {
			panic(err)
		}
		features = append(features, res)
	}

	return training.Example{
		Response: resEncoded,
		Input:    features,
	}
}

func onehot(classes int, val float64) []float64 {
	res := make([]float64, classes)
	res[int(val)-1] = 1
	return res
}
