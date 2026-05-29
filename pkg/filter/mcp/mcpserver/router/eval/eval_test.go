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

package eval

import (
	"os"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func loadSample(t *testing.T) Dataset {
	data, err := os.ReadFile("sample_dataset.yaml")
	require.NoError(t, err)
	var ds Dataset
	require.NoError(t, yaml.Unmarshal(data, &ds))
	return ds
}

func TestEvaluate_SampleDatasetRecall(t *testing.T) {
	ds := loadSample(t)
	require.NotEmpty(t, ds.Tools)
	require.NotEmpty(t, ds.Cases)

	res := Evaluate(ds, model.SchemaWeights{}, 3)

	// With deterministic tag/capability/description matching the expected tools
	// should rank within the top 3 for every case.
	assert.GreaterOrEqual(t, res.MeanRecall, 0.8, "Recall@3 should be >= 0.8")
	assert.Equal(t, len(ds.Cases), res.Cases)
}

func TestScore(t *testing.T) {
	recall, precision := score([]string{"a", "b"}, []string{"a", "c"})
	assert.Equal(t, 0.5, recall)    // 1 of 2 expected found
	assert.Equal(t, 0.5, precision) // 1 of 2 predicted correct
}

func TestScore_EmptyExpected(t *testing.T) {
	recall, precision := score([]string{"a"}, nil)
	assert.Equal(t, 1.0, recall)
	assert.Equal(t, 1.0, precision)
}

func TestEvaluate_EmptyDataset(t *testing.T) {
	res := Evaluate(Dataset{}, model.SchemaWeights{}, 5)
	assert.Equal(t, 0, res.Cases)
	assert.Equal(t, 0.0, res.MeanRecall)
}
