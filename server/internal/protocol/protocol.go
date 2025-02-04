package protocol

const (
    CmdHello     = "HELLO"
    CmdChallenge = "CHALLENGE"
    CmdSolution  = "SOLUTION"
    CmdQuote     = "QUOTE"
    CmdError     = "ERROR"
)

type Message struct {
    Command string
    Payload string
}

func (m *Message) String() string {
    if m.Payload == "" {
        return m.Command
    }
    return m.Command + " " + m.Payload
} 