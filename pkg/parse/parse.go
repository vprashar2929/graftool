package parse

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
)

func ParseEpoch(timestamp interface{}) time.Time {
	flt, _, err := big.ParseFloat(fmt.Sprint("", timestamp), 10, 0, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)

	}
	i, _ := flt.Int64()
	return time.Unix(i, 0)
}

func ParseQuery(query, namespace, node, job, interval string) string {
	if strings.Contains(query, "$node") || strings.Contains(query, "$job") || strings.Contains(query, "$namespace") || strings.Contains(query, "$interval") || strings.Contains(query, "$__rate_interval") {
		replacer := strings.NewReplacer("$node", node, "$job", job, "$namespace", namespace, "$interval", interval, "$__rate_interval", interval)
		resquery := replacer.Replace(query)
		return resquery
	}
	return query
}
