package networkdata

import "fmt"

func GetNetworkDataName(bits, binSize, N int, id string, iteration int) string {
	id = CombineIdIteration(id, iteration)
	return fmt.Sprintf("nodes_data_b%d_k%d_%d_%s.txt", bits, binSize, N, id)
}

func CombineIdIteration(id string, iteration int) string {
	if iteration >= 0 {
		iterstr := fmt.Sprintf("i%d", iteration)
		if len(id) > 0 {
			id = id + "-" + iterstr
		} else {
			id = iterstr
		}
	}
	return id
}
