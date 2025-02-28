package protocol

// Константы команд протокола
const (
    // CmdHello - начальная команда от клиента
    CmdHello     = "HELLO"
    // CmdChallenge - команда отправки вызова от сервера
    CmdChallenge = "CHALLENGE"
    // CmdSolution - команда отправки решения от клиента
    CmdSolution  = "SOLUTION"
    // CmdQuote - команда отправки цитаты от сервера
    CmdQuote     = "QUOTE"
    // CmdError - команда сообщения об ошибке
    CmdError     = "ERROR"
)

// Message представляет структуру сообщения протокола
type Message struct {
    // Command - команда сообщения
    Command string
    // Payload - полезная нагрузка сообщения
    Payload string
}

// String возвращает строковое представление сообщения
func (m *Message) String() string {
    if m.Payload == "" {
        return m.Command
    }
    return m.Command + " " + m.Payload
} 