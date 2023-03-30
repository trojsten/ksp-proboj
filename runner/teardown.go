package main

func (m *Match) teardown() {
	var err error
	if m.observer != nil {
		m.logger.Debug("Closing observer file")
		err = m.observer.Close()
		if err != nil {
			m.logger.Error("Could not close observer file", "err", err)
		}
	}

	if m.Server.IsRunning() {
		m.logger.Debug("Killing server")
		err = m.Server.Kill()
		if err != nil {
			m.logger.Error("Could not kill server")
		}
	}

	for name, process := range m.Players {
		m.logger.Debug("Killing player", "player", name)
		if process.IsRunning() {
			err := process.Kill()
			if err != nil {
				m.logger.Error("Could not kill player", "player", name, "err", err)
			}
		}
	}

	return
}
