package db

// GetConfig returns a config value by key.
func (d *DB) GetConfig(key string) (string, error) {
	var val string
	err := d.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&val)
	return val, err
}

// SetConfig sets a config value.
func (d *DB) SetConfig(key, value string) error {
	_, err := d.Exec("INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)", key, value)
	return err
}
