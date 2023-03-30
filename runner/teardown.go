package main

func (m *Match) teardown() (err error) {
	if m.observer != nil {
		m.logger.Debug("Closing observer file")
		err = m.observer.Close()
		if err != nil {
			return
		}
	}

	for name, process := range m.Players {
		m.logger.Debug("Killing player", "player", name)
		if process.IsRunning() {
			err := process.Kill()
			if err != nil {
				m.logger.Error("Could not kill player", "err", err)
			}
		}
	}

	return
}
