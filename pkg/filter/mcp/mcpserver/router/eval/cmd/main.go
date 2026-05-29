/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Command routereval runs the offline SchemaMatcher evaluation over a YAML
// dataset and prints Recall@K / Precision@K for a range of K values.
//
// Usage:
//
//	go run ./pkg/filter/mcp/mcpserver/router/eval/cmd -dataset sample_dataset.yaml -k 3,5,10
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

import (
	"gopkg.in/yaml.v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router/eval"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func main() {
	datasetPath := flag.String("dataset", "sample_dataset.yaml", "path to evaluation dataset YAML")
	kList := flag.String("k", "3,5,10", "comma-separated K values to evaluate")
	flag.Parse()

	data, err := os.ReadFile(*datasetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read dataset: %v\n", err)
		os.Exit(1)
	}

	var ds eval.Dataset
	if err := yaml.Unmarshal(data, &ds); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse dataset: %v\n", err)
		os.Exit(1)
	}

	weights := model.SchemaWeights{} // zeros => matcher applies defaults

	fmt.Printf("Dataset: %s  (%d tools, %d cases)\n", *datasetPath, len(ds.Tools), len(ds.Cases))
	fmt.Printf("%-6s %-12s %-14s\n", "K", "Recall@K", "Precision@K")
	for _, kStr := range strings.Split(*kList, ",") {
		k, err := strconv.Atoi(strings.TrimSpace(kStr))
		if err != nil || k <= 0 {
			continue
		}
		res := eval.Evaluate(ds, weights, k)
		fmt.Printf("%-6d %-12.3f %-14.3f\n", res.K, res.MeanRecall, res.MeanPrecision)
	}
}
