package main

func (m *Match) Run() {
	defer func(m *Match) {
		m.teardown()
	}(m)

	err := m.preflight()
	if err != nil {
		m.Log.Error("Error in preflight", "err", err)
	}

	for m.Server.IsRunning() && !m.Ended {
		m.Log.Debug("Waiting for command from server")
		cmd, err := m.Server.Read()
		if err != nil {
			m.Log.Error("Error while reading command from server", "err", err)
			return
		}
		m.parseCommand(cmd)
	}

	// Wait for exit to get handled.
	<-m.Server.OnExit()

	if !m.Ended {
		m.Log.Warn("Server exited without proper game end", "exit", m.Server.Exit, "err", m.Server.Error)
	}
}
