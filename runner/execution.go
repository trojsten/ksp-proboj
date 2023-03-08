package main

func (m *Match) Run() {
	err := m.preflight()
	if err != nil {
		m.logger.Error("Error in preflight", "err", err)
		return
	}

	for m.Server.IsRunning() {
		m.logger.Debug("Waiting for command from server")
		cmd, err := m.Server.Read()
		if err != nil {
			m.logger.Error("Error while reading command from server", "err", err)
			return
		}
		m.parseCommand(cmd)
	}

	// Wait for exit to get handled.
	<-m.Server.OnExit()

	if !m.ended {
		m.logger.Warn("Server exited without proper game end", "exit", m.Server.Exit, "err", m.Server.Error)
	}

	for name, process := range m.Players {
		m.logger.Debug("Killing player", "player", name)
		err := process.Kill()
		if err != nil {
			m.logger.Error("Could not kill player", "err", err)
		}
	}
}
