package internal

type HostFilter struct {
	hosts []*Host
}

func (f *HostFilter) ExcludeHidden() *HostFilter {
	hosts := make([]*Host, 0)
	for _, h := range f.hosts {
		if !h.Hidden {
			hosts = append(hosts, h)
		}
	}
	f.hosts = hosts
	return f
}

func (f *HostFilter) GetHosts() []*Host {
	return f.hosts
}

func (f *HostFilter) GetHostByName(name string) *Host {
	for _, h := range f.hosts {
		if h.Name == name {
			return h
		}
	}
	return nil
}
