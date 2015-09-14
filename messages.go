type Message struct {
    MID                  int     `gorm:"primary_key"`
    Content              string
    Sender               User       `sql:"not null"`
    Timestamp            time.Time
    MessageRecipientID   int
    MessageRecipientType string
}