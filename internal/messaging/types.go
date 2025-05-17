package messaging

type DropSessionEmailPayload struct {
	GuardianName  string
	GuardianEmail string
	Children      []string
	Secret        string
	Date          string
	ChurchName    string
}
