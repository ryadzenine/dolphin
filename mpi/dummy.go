package mpi

type Dummy map[string]Versioner

func (m Dummy) Queues() []string {
	names := make([]string, 0, len(m))
	for key := range m {
		names = append(names, key)
	}
	return names
}

func (m Dummy) Register(name string) bool {
	_, ok := m[name]
	if !ok {
		m[name] = nil
		return true
	}
	return false
}

func (m Dummy) Write(queue string, data Versioner) {
	m[queue] = data
}

func (m Dummy) ReadFirst(queue string) (data Versioner) {
	return m[queue]
}

func (m Dummy) ReadFirstAll() (data map[string]Versioner) {
	return m
}

func (m Dummy) ReadStates(versions map[string]int) map[string]Versioner {
	data := make(map[string]Versioner)
	for key, v := range m {
		if v.Version() > versions[key] {
			data[key] = v
		}
	}
	return data
}

func NewDummy() (queue Dummy) {
	return Dummy(make(map[string]Versioner))
}
