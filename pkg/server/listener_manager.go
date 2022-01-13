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
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"runtime/debug"
)

// wrapListenerService wrap listener service and its configuration.
type wrapListenerService struct {
	listener.ListenerService
	cfg *model.Listener
}

// ListenerManager the listener manager
type ListenerManager struct {
	activeListener        []*model.Listener
	bootstrap             *model.Bootstrap
	activeListenerService []*wrapListenerService
}

// CreateDefaultListenerManager create listener manager from config
func CreateDefaultListenerManager(bs *model.Bootstrap) *ListenerManager {
	sl := bs.GetStaticListeners()
	var listeners []*wrapListenerService
	for _, lsCof := range bs.StaticResources.Listeners {
		ls, err := listener.CreateListenerService(lsCof, bs)
		if err != nil {
			logger.Error("CreateDefaultListenerManager %s error: %v", lsCof.Name, err)
		}
		listeners = append(listeners, &wrapListenerService{ls, lsCof})
	}

	return &ListenerManager{
		activeListener:        sl,
		activeListenerService: listeners,
		bootstrap:             bs,
	}
}

func (lm *ListenerManager) AddOrUpdateListener(lsConf *model.Listener) error {
	//todo add sync lock for concurrent using
	if theListener := lm.getListener(lsConf.Name); theListener != nil {
		lm.updateListener(theListener, lsConf)
		return nil
	}
	ls, err := listener.CreateListenerService(lsConf, lm.bootstrap)
	if err != nil {
		return err
	}
	lm.startListenerServiceAsync(ls)
	lm.addListenerService(ls, lsConf)
	lm.activeListener = append(lm.activeListener, lsConf)
	return nil
}

func (lm *ListenerManager) updateListener(listener *model.Listener, to *model.Listener) {
	//todo update listener and service
	panic("not implement")
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

func (lm *ListenerManager) startListenerServiceAsync(s listener.ListenerService) chan<- struct{} {
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
		err := s.Start()
		if err != nil {

		}
	}()
	return done
}

func (lm *ListenerManager) addListenerService(ls listener.ListenerService, lsConf *model.Listener) {
	lm.activeListenerService = append(lm.activeListenerService, &wrapListenerService{ls, lsConf})
}

func (lm *ListenerManager) GetListenerService(name string) listener.ListenerService {
	for i := range lm.activeListenerService {
		if lm.activeListenerService[i].cfg.Name == name {
			return lm.activeListenerService[i]
		}
	}
	return nil
}

func (lm *ListenerManager) RemoveListener(names []string) {
	//todo implement remove Listener and ListenerService
}
