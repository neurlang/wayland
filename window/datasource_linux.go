package window

import wl "github.com/neurlang/wayland/wl"
import "github.com/neurlang/wayland/wlclient"

type DataSource struct {
	src        *wl.DataSource
	CopyBuffer string
}

func (ds *DataSource) Offer(str string) {
	_ = ds.src.Offer(str) // Ignore error during offer
}

func (ds *DataSource) AddListener(l wlclient.DataSourceListener) {
	wlclient.DataSourceAddListener(ds.src, l)
}

func (ds *DataSource) RemoveListener(l wlclient.DataSourceListener) {
	wlclient.DataSourceRemoveListener(ds.src, l)
}
