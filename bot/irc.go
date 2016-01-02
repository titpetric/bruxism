package bot

import (
	"errors"

	"github.com/fluffle/goirc/client"
)

// IRCServiceName is the service name for the IRC service.
const IRCServiceName string = "IRC"

// IRCMessage is a Message wrapper around client.Line.
type IRCMessage client.Line

// Channel returns the channel id for this message.
func (m *IRCMessage) Channel() string {
	i := client.Line(*m)
	return i.Target()
}

// UserName returns the user name for this message.
func (m *IRCMessage) UserName() string {
	return m.Nick
}

// UserID returns the user id for this message.
func (m *IRCMessage) UserID() string {
	return m.Nick
}

// UserAvatar returns the avatar url for this message.
func (m *IRCMessage) UserAvatar() string {
	return ""
}

// Message returns the message content for this message.
func (m *IRCMessage) Message() string {
	i := client.Line(*m)
	return i.Text()
}

// MessageID returns the message ID for this message.
func (m *IRCMessage) MessageID() string {
	return ""
}

// IsModerator returns whether or not the sender of this message is a moderator.
func (m *IRCMessage) IsModerator() bool {
	return false
}

// IRC is a Service provider for IRC.
type IRC struct {
	host        string
	nick        string
	channels    []string
	Conn        *client.Conn
	messageChan chan Message
}

// NewIRC creates a new IRC service.
func NewIRC(host, nick string, channels []string) *IRC {
	return &IRC{
		host:        host,
		nick:        nick,
		channels:    channels,
		messageChan: make(chan Message, 200),
	}
}

func (i *IRC) onMessage(conn *client.Conn, line *client.Line) {
	m := IRCMessage(*line)
	i.messageChan <- &m
}

func (i *IRC) onConnect(conn *client.Conn, line *client.Line) {
	for _, c := range i.channels {
		conn.Join(c)
	}
}

func (i *IRC) onDisconnect(conn *client.Conn, line *client.Line) {
	conn.ConnectTo(i.host)
}

// Name returns the name of the service.
func (i *IRC) Name() string {
	return IRCServiceName
}

// Open opens the service and returns a channel which all messages will be sent on.
func (i *IRC) Open() (<-chan Message, error) {
	i.Conn = client.SimpleClient(i.nick, "Septapus", "Septapus")
	i.Conn.Config().Version = "Septapus"
	i.Conn.Config().QuitMessage = ""

	i.Conn.HandleFunc("connected", i.onConnect)

	i.Conn.HandleFunc("disconnected", i.onDisconnect)

	i.Conn.HandleFunc(client.PRIVMSG, i.onMessage)

	go i.Conn.ConnectTo(i.host)

	return i.messageChan, nil
}

// IsMe returns whether or not a message was sent by the bot.
func (i *IRC) IsMe(message Message) bool {
	return message.UserName() == i.UserName()
}

// SendMessage sends a message.
func (i *IRC) SendMessage(channel, message string) error {
	i.Conn.Privmsg(channel, message)
	return nil
}

// DeleteMessage deletes a message.
func (i *IRC) DeleteMessage(channel, messageID string) error {
	return errors.New("Deleting messages not supported on IRC.")
}

// BanUser bans a user.
func (i *IRC) BanUser(channel, userID string, duration int) error {
	i.Conn.Kick(channel, userID)
	return nil
}

// UserName returns the bots name.
func (i *IRC) UserName() string {
	return i.Conn.Me().Nick
}

// SetPlaying will set the current game being played by the bot.
func (i *IRC) SetPlaying(game string) error {
	return errors.New("Set playing not supported on IRC.")
}

// Join will join a channel.
func (i *IRC) Join(join string) error {
	i.Conn.Join(join)
	return nil
}

// Typing sets that the bot is typing.
func (i *IRC) Typing(channel string) error {
	return errors.New("Typing not supported on IRC.")
}

// PrivateMessage will send a private message to a user.
func (i *IRC) PrivateMessage(userID, message string) error {
	return i.SendMessage(userID, message)
}

// SupportsMultiline returns whether the service supports multiline messages.
func (i *IRC) SupportsMultiline() bool {
	return false
}

// CommandPrefix returns the command prefix for the service.
func (i *IRC) CommandPrefix() string {
	return "!"
}

// IsPrivate returns whether or not a message was private.
func (i *IRC) IsPrivate(message Message) bool {
	return message.UserName() == message.Channel()
}
