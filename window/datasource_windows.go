package window

import "github.com/neurlang/wayland/wlclient"

type DataSource struct {
	CopyBuffer string
}

func (s *DataSource) RemoveListener(l wlclient.DataSourceListener) {

}

func (s *DataSource) Offer(s2 string) {

}

func (s *DataSource) AddListener(l wlclient.DataSourceListener) {

}
