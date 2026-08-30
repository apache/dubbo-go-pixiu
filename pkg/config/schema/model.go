/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package schema

import "github.com/apache/dubbo-go-pixiu/pkg/common/copyutil"

const (
	APIVersionV1Alpha1 = "pixiu.apache.org/v1alpha1"

	KindConfigSet = "ConfigSet"
	KindResource  = "Resource"
	KindMethod    = "Method"
	KindListener  = "Listener"
	KindCluster   = "Cluster"
)

type Lifecycle string

const (
	LifecycleDraft     Lifecycle = "draft"
	LifecyclePublished Lifecycle = "published"
	LifecycleArchived  Lifecycle = "archived"
)

// ConfigObject is the stable envelope for every Admin-managed object. Spec is
// intentionally dynamic; its semantics come from the matching ObjectSchema.
type ConfigObject struct {
	APIVersion string         `json:"apiVersion" yaml:"apiVersion"`
	Kind       string         `json:"kind" yaml:"kind"`
	Metadata   ObjectMetadata `json:"metadata" yaml:"metadata"`
	Spec       map[string]any `json:"spec" yaml:"spec"`
	Status     *ObjectStatus  `json:"status,omitempty" yaml:"status,omitempty"`
}

type ObjectMetadata struct {
	ID              string            `json:"id" yaml:"id"`
	Name            string            `json:"name" yaml:"name"`
	Revision        int64             `json:"revision,omitempty" yaml:"revision,omitempty"`
	Lifecycle       Lifecycle         `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	OwnerReferences []ObjectReference `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
}

type ObjectReference struct {
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string `json:"kind" yaml:"kind"`
	ID         string `json:"id" yaml:"id"`
}

type ObjectStatus struct {
	Phase            string      `json:"phase,omitempty" yaml:"phase,omitempty"`
	ObservedRevision int64       `json:"observedRevision,omitempty" yaml:"observedRevision,omitempty"`
	Conditions       []Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
}

type Condition struct {
	Type    string `json:"type" yaml:"type"`
	Status  string `json:"status" yaml:"status"`
	Reason  string `json:"reason,omitempty" yaml:"reason,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

// ConfigSet is a versionable snapshot boundary. Publishing a ConfigSet keeps
// Resource/Method and Listener/Cluster changes in one revision.
type ConfigSet struct {
	APIVersion string         `json:"apiVersion" yaml:"apiVersion"`
	Kind       string         `json:"kind" yaml:"kind"`
	Metadata   ObjectMetadata `json:"metadata" yaml:"metadata"`
	Objects    []ConfigObject `json:"objects" yaml:"objects"`
}

func (o ConfigObject) Clone() ConfigObject {
	cloned := o
	cloned.Metadata = o.Metadata.clone()
	cloned.Spec = copyutil.CloneStringAnyMap(o.Spec)
	if o.Status != nil {
		status := *o.Status
		status.Conditions = append([]Condition(nil), o.Status.Conditions...)
		cloned.Status = &status
	}
	return cloned
}

func (m ObjectMetadata) clone() ObjectMetadata {
	cloned := m
	cloned.Labels = copyutil.CloneStringMap(m.Labels)
	cloned.Annotations = copyutil.CloneStringMap(m.Annotations)
	cloned.OwnerReferences = append([]ObjectReference(nil), m.OwnerReferences...)
	return cloned
}

func (s ConfigSet) Clone() ConfigSet {
	cloned := s
	cloned.Metadata = s.Metadata.clone()
	cloned.Objects = make([]ConfigObject, len(s.Objects))
	for i := range s.Objects {
		cloned.Objects[i] = s.Objects[i].Clone()
	}
	return cloned
}
