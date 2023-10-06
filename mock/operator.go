package mock

import (
	"context"
	"gopkg.in/irc.v4"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

type Config struct {
	Server           string
	SearchResultPath string
	BookResultPath   string
}

// Operator listens to an IRC channel and responds to search / download requests
type Operator struct {
	config     *Config
	client     *irc.Client
	log        *log.Logger
	dccManager *DccManager
}

func NewOperator(config *Config) *Operator {
	conn, err := net.Dial("tcp", config.Server)
	if err != nil {
		panic(err)
	}

	operator := &Operator{
		config:     config,
		dccManager: NewDccManager(),
		log:        log.New(os.Stdout, "OPERATOR: ", 0),
	}

	client := irc.NewClient(conn, irc.ClientConfig{
		Nick:           "search",
		Pass:           "",
		User:           "search",
		Name:           "search",
		EnableISupport: true,
		EnableTracker:  true,
		PingFrequency:  time.Second * 30,
		PingTimeout:    time.Second * 30,
		SendLimit:      0,
		SendBurst:      0,
		Handler:        irc.HandlerFunc(operator.Handler),
	})

	operator.client = client

	return operator
}

func (o *Operator) Handler(client *irc.Client, message *irc.Message) {
	// Skip 353 messages
	if message.Command == "353" {
		return
	}

	if strings.HasPrefix(message.Trailing(), "@search") {
		query := strings.SplitN(message.Trailing(), " ", 2)[1]
		o.log.Printf("Search for '%s'\n", query)

		dccString := o.dccManager.ServeFile("SearchBot_results_for_ the great gatsby.txt.zip")

		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
		err := client.Writef(`:search PRIVMSG %s :%s`, message.Name, dccString)
		if err != nil {
			log.Println(err)
		}
	}

	if strings.HasPrefix(message.Trailing(), "!") {
		query := strings.SplitN(message.Trailing(), " ", 2)[1]
		o.log.Printf("Download request for '%s'\n", query)

		dccString := o.dccManager.ServeFile("ebook.epub")

		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
		err := client.Writef(`:search PRIVMSG %s :%s`, message.Name, dccString)
		if err != nil {
			log.Println(err)
		}
	}
}

func (o *Operator) StartListening(ctx context.Context) error {
	go func() {
		time.Sleep(time.Second * 2)
		o.client.Write("JOIN #ebooks")
		o.log.Println("Joined #ebooks channel")
	}()

	o.log.Printf("Connected to %s\n", o.config.Server)
	return o.client.RunContext(ctx)
}
