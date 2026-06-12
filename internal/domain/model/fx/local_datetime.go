package fx

import (
	"encoding/json"
	"time"
)

// LocalDateTime はレスポンス用に "yyyy-MM-dd HH:mm:ss" 形式でシリアライズされる時刻型。
// JavaのTradeEntry.contractAt @JsonFormat(pattern="yyyy-MM-dd HH:mm:ss") に対応。
type LocalDateTime struct {
	time.Time
}

func (d LocalDateTime) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Format("2006-01-02 15:04:05"))
}
