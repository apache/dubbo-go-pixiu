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

package server

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type ListenerManager struct {
	activeListener        []*model.Listener
	activeListenerService []*ListenerService
}

func CreateDefaultListenerManager(bs *model.Bootstrap) *ListenerManager {
	sl := bs.GetStaticListeners()
	var ls []*ListenerService
	for _, l := range bs.StaticResources.Listeners {
		listener := CreateListenerService(l, bs)
		ls = append(ls, listener)
	}

	return &ListenerManager{
		activeListener:        sl,
		activeListenerService: ls,
	}
}

func (lm *ListenerManager) addOrUpdateListener(l *model.Listener) {
	lm.activeListener = append(lm.activeListener, l)
}

func (lm *ListenerManager) StartListen() {
	for _, s := range lm.activeListenerService {
		go s.Start()
	}
}

func (lm *ListenerManager) addListenerService(ls *ListenerService) {
	lm.activeListenerService = append(lm.activeListenerService, ls)
}
