package main

func (m *Match) teardown() {
	var err error

	m.Log.Debug("Closing observer file.")
	err = m.Observer.Close()
	if err != nil {
		m.Log.Error("Could not close observer file.", "err", err)
	}

	if m.Server.IsRunning() {
		m.Log.Debug("Killing server")
		err = m.Server.Kill()
		if err != nil {
			m.Log.Error("Could not kill server")
		}
	}

	for name, process := range m.Players {
		m.Log.Debug("Killing player", "player", name)
		if process.IsRunning() {
			err := process.Kill()
			if err != nil {
				m.Log.Error("Could not kill player", "player", name, "err", err)
			}
		}
	}

	signalMatchEnd(m)
}
