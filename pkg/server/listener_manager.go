package server

import "github.com/apache/dubbo-go-pixiu/pkg/model"

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
		activeListenerService: make([]*ListenerService, 0),
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
