package ts3

type ServerGroup struct {
	ID   int
	Name string
}

func (c *TS3Client) ListServerGroups() ([]ServerGroup, error) {
	rows, err := c.query("servergrouplist")
	if err != nil {
		return nil, err
	}

	groups := make([]ServerGroup, 0, len(rows))
	for _, row := range rows {
		id := atoiDefault(row["sgid"], 0)
		if id <= 0 {
			continue
		}
		groups = append(groups, ServerGroup{
			ID:   id,
			Name: row["name"],
		})
	}

	return groups, nil
}
