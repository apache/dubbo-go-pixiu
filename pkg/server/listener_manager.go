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
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"runtime/debug"
)

type ListenerManager struct {
	activeListener        []*model.Listener
	activeListenerService []*ListenerService
	bootstrap             *model.Bootstrap
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
		bootstrap:             bs,
	}
}

func (lm *ListenerManager) AddOrUpdateListener(l *model.Listener) {
	if theListener := lm.getListener(l.Name); theListener != nil {
		//todo to implement update
		panic("to implement")
	}
	listener := CreateListenerService(l, lm.bootstrap)
	lm.startListenerServiceAsync(listener)
	lm.addListenerService(listener)
	lm.activeListener = append(lm.activeListener, l)
}

func (lm *ListenerManager) getListener(name string) *model.Listener {
	for _, l := range lm.activeListener {
		if l.Name == name {
			return l
		}
	}
	return nil
}

func (lm *ListenerManager) StartListen() {
	for _, s := range lm.activeListenerService {
		lm.startListenerServiceAsync(s)
	}
}

func (lm *ListenerManager) startListenerServiceAsync(s *ListenerService) chan<- struct{} {
	done := make(chan struct{})
	go func() {
		defer func() {
			panicErr := recover()
			if panicErr != nil {
				logger.Errorf("recover from panic %v", panicErr)
				debug.PrintStack()
			}
			close(done)
		}()
		s.Start()
	}()
	return done
}

func (lm *ListenerManager) addListenerService(ls *ListenerService) {
	lm.activeListenerService = append(lm.activeListenerService, ls)
}

func (lm *ListenerManager) GetListenerService(name string) *ListenerService {
	for i := range lm.activeListenerService {
		if lm.activeListenerService[i].cfg.Name == name {
			return lm.activeListenerService[i]
		}
	}
	return nil
}
