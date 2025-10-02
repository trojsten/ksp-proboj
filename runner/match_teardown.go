package main

// teardown performs cleanup operations when a match ends.
// It ensures all resources are properly released, including closing the observer
// file, terminating the server process, and killing all player processes.
// This method is called automatically when a match completes or fails.
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
		m.Server.WaitForEnd()
	}

	for name, process := range m.Players {
		m.Log.Debug("Killing player", "player", name)
		if process.IsRunning() {
			err := process.Kill()
			if err != nil {
				m.Log.Error("Could not kill player", "player", name, "err", err)
			}
		}
		process.WaitForEnd()
	}

	signalMatchEnd(m)
}
