package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// record[2]  = Timestamp
// record[4]  = um
// record[6]  = dm
// record[8]  = Rt
type data struct {
	Timestamp int `json:"ts"`
	Rt        int `json:"rt"`
}

func generateGraph(dataSetNumber int, graph map[string][]data) {
	for i := 0; i < dataSetNumber; i++ {
		var path = fmt.Sprintf("/home/mahmoud/Desktop/research/auto-scaler/data/MSCallGraph_%d.csv", i)
		fmt.Println("parse file: ", path)
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		csvReader := csv.NewReader(file)
		csvReader.Read()
		if err != nil {
			return
		} // ignore header
		for {
			record, err := csvReader.Read()
			if err != nil {
				if err.Error() != "EOF" {
					panic(err)
				}
				break
			}
			timestamp, err := strconv.ParseInt(record[2], 10, 64)
			if err != nil {
				panic(err)
			}
			um := record[4]
			dm := record[6]
			rt, err := strconv.ParseInt(record[8], 10, 64)
			if err != nil {
				panic(err)
			}

			key := fmt.Sprintf("%s-%s", um, dm)
			_, exists := graph[key]
			if !exists {
				graph[key] = make([]data, 0)
			}
			graph[key] = append(graph[key], data{
				Timestamp: int(timestamp),
				Rt:        int(rt),
			})
		}
		file.Close()
	}

}

func sortByTimestamp(graph map[string][]data) {
	for _, val := range graph {
		sort.SliceStable(val, func(i, j int) bool {
			return val[i].Timestamp < val[j].Timestamp
		})
	}
}

func exportGraphAndEdges(path string, graph map[string][]data) {
	file, err := os.Create(path + "graph.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"um", "dm", "edgeValue"})
	for key, val := range graph {
		parts := strings.Split(key, "-")
		um := parts[0]
		dm := parts[1]
		marshal, err := json.Marshal(val)
		if err != nil {
			panic(err)
		}
		err = csvWriter.Write([]string{um, dm, string(marshal)})
		if err != nil {
			panic(err)
		}
	}
	csvWriter.Flush()

}
func exportUniqueMicroServices(path string, graph map[string][]data) {
	file, err := os.Create(path + "./uniqueMicroservices.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"microservice"})
	var uniqueMicroservices = make(map[string]int)

	for key := range graph {
		parts := strings.Split(key, "-")
		um := parts[0]
		dm := parts[1]
		uniqueMicroservices[um] = 1
		uniqueMicroservices[dm] = 1
	}

	for key := range uniqueMicroservices {
		csvWriter.Write([]string{key})
	}
	csvWriter.Flush()

}
func exportGraphToCsv(path string, graph map[string][]data) {
	exportUniqueMicroServices(path, graph)
	exportGraphAndEdges(path, graph)
}

func main() {
	dataSetNumber := 3
	var graph = make(map[string][]data)
	fmt.Println("start parsing files ...")
	generateGraph(dataSetNumber, graph)

	fmt.Println("start sorting by timestamps")
	sortByTimestamp(graph)

	fmt.Println("exporting data ...")
	exportGraphToCsv("./", graph)
}
