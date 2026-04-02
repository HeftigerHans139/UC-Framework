package ts3

import "fmt"

type Channel struct {
	ID   int
	Name string
}

func (c *TS3Client) ListChannels() ([]Channel, error) {
	rows, err := c.query("channellist")
	if err != nil {
		return nil, err
	}

	channels := make([]Channel, 0, len(rows))
	for _, row := range rows {
		channels = append(channels, Channel{
			ID:   atoiDefault(row["cid"], 0),
			Name: row["channel_name"],
		})
	}
	return channels, nil
}

func (c *TS3Client) RenameChannel(channelID int, newName string) error {
	if channelID <= 0 {
		return fmt.Errorf("invalid channel id")
	}
	if newName == "" {
		return fmt.Errorf("channel name must not be empty")
	}
	_, err := c.query("channeledit cid=" + intToString(channelID) + " channel_name=" + escapeTS3(newName))
	return err
}

func (c *TS3Client) SetChannelAccess(channelID int, neededJoinPower int, neededSubscribePower int) error {
	if channelID <= 0 {
		return fmt.Errorf("invalid channel id")
	}
	if neededJoinPower < 0 || neededSubscribePower < 0 {
		return fmt.Errorf("invalid channel access power")
	}
	if _, err := c.query(
		"channeladdperm cid=" + intToString(channelID) +
			" permsid=i_channel_needed_join_power permvalue=" + intToString(neededJoinPower),
	); err != nil {
		return err
	}

	if _, err := c.query(
		"channeladdperm cid=" + intToString(channelID) +
			" permsid=i_channel_needed_subscribe_power permvalue=" + intToString(neededSubscribePower),
	); err != nil {
		return err
	}

	return nil
}
