package mpi

type Dummy map[string]Versionable

func (m Dummy) Queues() []string {
  names := make([]string, 0, len(m))
  for key, _ := range m {
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

func (m Dummy) Write(queue string, data Versionable) {
  m[queue] = data
}

func (m Dummy) ReadFirst(queue string) (data Versionable) {
  return m[queue]
}

func (m Dummy) ReadFirstAll() (data map[string]Versionable) {
  return m
}

func (m Dummy) ReadStates(versions map[string]int) map[string]Versionable {
  data := make(map[string]Versionable)
  for key, v := range m {
    if v.Version() > versions[key] {
      data[key] = v
    }
  }
  return data
}

func NewDummy(capacity int) (queue Dummy) {
  return Dummy(make(map[string]Versionable))
}
