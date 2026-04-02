package ts3

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (c *TS3Client) query(command string) ([]map[string]string, error) {
	c.queryMu.Lock()
	defer c.queryMu.Unlock()

	c.lock.RLock()
	conn := c.conn
	reader := c.reader
	writer := c.writer
	c.lock.RUnlock()

	if conn == nil || reader == nil || writer == nil {
		return nil, errors.New("ts3 query not connected")
	}

	if _, err := writer.WriteString(command + "\n"); err != nil {
		c.MarkDisconnected(err)
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		c.MarkDisconnected(err)
		return nil, err
	}

	var payloadLines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			c.MarkDisconnected(err)
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "error ") {
			errInfo := parseKV(strings.TrimPrefix(line, "error "))
			if errInfo["id"] != "0" {
				return nil, fmt.Errorf("ts3 error id=%s msg=%s", errInfo["id"], unescapeTS3(errInfo["msg"]))
			}
			break
		}
		payloadLines = append(payloadLines, line)
	}

	if len(payloadLines) == 0 {
		return []map[string]string{}, nil
	}

	items := strings.Split(strings.Join(payloadLines, ""), "|")
	result := make([]map[string]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		result = append(result, parseKV(item))
	}
	return result, nil
}

func parseKV(s string) map[string]string {
	out := map[string]string{}
	parts := strings.Fields(s)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		out[kv[0]] = unescapeTS3(kv[1])
	}
	return out
}

func escapeTS3(s string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		" ", "\\s",
		"/", "\\/",
		"|", "\\p",
		"\a", "\\a",
		"\b", "\\b",
		"\f", "\\f",
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\v", "\\v",
	)
	return replacer.Replace(s)
}

func unescapeTS3(s string) string {
	replacer := strings.NewReplacer(
		"\\s", " ",
		"\\/", "/",
		"\\p", "|",
		"\\a", string(rune(7)),
		"\\b", string(rune(8)),
		"\\f", string(rune(12)),
		"\\n", "\n",
		"\\r", "\r",
		"\\t", "\t",
		"\\v", "\v",
		"\\\\", "\\",
	)
	return replacer.Replace(s)
}

func intToString(v int) string {
	return strconv.Itoa(v)
}

func atoiDefault(v string, d int) int {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return d
	}
	return n
}
