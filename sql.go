package stats

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Metadata map[string]string

func (m *Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Metadata) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	switch data := v.(type) {
	case string:
		return json.Unmarshal([]byte(data), m)
	case []byte:
		err := json.Unmarshal(data, m)
		return err
	default:
		return fmt.Errorf("cannot scan type %t into Matadata", v)
	}
}
