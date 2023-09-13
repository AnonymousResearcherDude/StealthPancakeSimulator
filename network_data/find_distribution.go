package networkdata

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func isZipf() {
	file, err := os.Open("workload.bin")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	freqMap := make(map[int]int)

	for {
		var frequency uint32
		var id uint32
		err := binary.Read(file, binary.LittleEndian, &frequency)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading frequency: ", err)
			return
		}

		err = binary.Read(file, binary.LittleEndian, &id)
		if err != nil {
			fmt.Println("Error reading number: ", err)
			return
		}
		freqMap[int(id)] = int(frequency)
	}

	frequencies := make([]int, 0, len(freqMap))
	for _, freq := range freqMap {
		frequencies = append(frequencies, freq)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(frequencies)))

	points := make(plotter.XYs, len(frequencies))
	for i, freq := range frequencies {
		points[i].X = math.Log(float64(i + 1)) // Rank
		points[i].Y = math.Log(float64(freq))  // Frequency
	}

	maxRank := float64(len(frequencies))
	zipfPoints := make(plotter.XYs, len(frequencies))
	for i := range zipfPoints {
		zipfPoints[i].X = math.Log(float64(i + 1))         // Rank
		zipfPoints[i].Y = math.Log(maxRank / float64(i+1)) // Ideal Zipf's Law
	}

	// Create a new plot
	p := plot.New()

	p.Title.Text = "Zipf's Law Test - generated chunks"
	p.X.Label.Text = "Log(Rank)"
	p.Y.Label.Text = "Log(Frequency)"

	s, err := plotter.NewScatter(points)
	if err != nil {
		log.Fatal(err)
	}
	s.GlyphStyle.Color = plotutil.Color(0)

	zipfLine, err := plotter.NewLine(zipfPoints)
	if err != nil {
		log.Fatal(err)
	}
	zipfLine.LineStyle.Color = plotutil.Color(1)

	p.Add(s, zipfLine)

	if err := p.Save(8*vg.Inch, 8*vg.Inch, "zipf_plot.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Plot saved as zipf_plot_test.png")
}
