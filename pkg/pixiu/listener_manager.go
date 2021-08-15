package pixiu

import "github.com/apache/dubbo-go-pixiu/pkg/model"

type ListenerManager struct {
	activeListener   []*model.Listener
	activeListenerService []*ListenerService
}

func CreateDefaultListenerManager(bs *model.Bootstrap) *ListenerManager {
	sl := bs.GetStaticListeners()
	return &ListenerManager{
		activeListener:   sl,
		activeListenerService: make([]*ListenerService, 0),
	}
}

func (lm *ListenerManager) addOrUpdateListener(l *model.Listener) {
	lm.activeListener = append(lm.activeListener, l)
}

func (lm *ListenerManager) StartListen() {
	listeners := lm.activeListener

	for _, s := range listeners {
		ls := CreateListenerService(s)
		lm.addListenerService(ls)
		go ls.Start()
	}
}

func (lm *ListenerManager) addListenerService(ls * ListenerService) {
	lm.activeListenerService = append(lm.activeListenerService, ls)
}