package discord

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/evan-buss/openbooks/core"
	"github.com/google/uuid"
)

// Config holds everything needed to run the Discord bot.
type Config struct {
	Token       string        // Discord bot token
	GuildID     string        // If set, commands register to this guild instantly (good for testing). Empty = global.
	DownloadDir string        // Directory where eBooks are saved (files land in <dir>/books).
	RateLimit   time.Duration // Minimum interval between searches.
	Debug       bool          // Enable verbose progress/debug logging on stdout.
	Log         bool          // Write raw IRC traffic to a log file under <dir>/logs.

	// IRC settings (sourced from the shared global flags).
	UserName  string
	Server    string
	EnableTLS bool
	SearchBot string
	Version   string
}

// selectPrefix identifies our select-menu component interactions and carries
// the session key after it: "request_select:<uuid>".
const selectPrefix = "request_select:"

// sessionTTL bounds how long a search result set is retained for selection.
const sessionTTL = 30 * time.Minute

type session struct {
	books   []core.BookDetail
	created time.Time
}

type bot struct {
	broker   *Broker
	booksDir string // absolute directory where downloads are saved

	mu       sync.Mutex
	sessions map[string]*session
}

// Run connects to IRC, opens the Discord session, registers the /request-book
// command, and runs until ctx is cancelled. The caller owns the lifecycle
// (e.g. signal handling); Run performs best-effort cleanup on exit.
//
// It is intended to run alongside the web server: when a token and guild are
// configured, server.Start launches this in a goroutine bound to its context.
func Run(ctx context.Context, cfg Config) error {
	if cfg.Token == "" {
		return errors.New("a Discord bot token is required")
	}

	broker, err := NewBroker(cfg)
	if err != nil {
		return fmt.Errorf("connecting to IRC: %w", err)
	}
	log.Printf("Discord: connected to IRC server %s as %s", cfg.Server, cfg.UserName)
	broker.Start(ctx)

	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		broker.Disconnect()
		return fmt.Errorf("creating discord session: %w", err)
	}
	// Slash commands and component interactions only need the Guilds intent.
	dg.Identify.Intents = discordgo.IntentsGuilds

	booksDir := filepath.Join(cfg.DownloadDir, "books")
	if abs, absErr := filepath.Abs(booksDir); absErr == nil {
		booksDir = abs
	}
	log.Printf("Discord: downloads will be saved to %s", booksDir)

	b := &bot{
		broker:   broker,
		booksDir: booksDir,
		sessions: make(map[string]*session),
	}

	dg.AddHandler(b.onInteraction)
	dg.AddHandler(func(_ *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Discord bot logged in as %s", r.User.String())
	})

	if err := dg.Open(); err != nil {
		broker.Disconnect()
		return fmt.Errorf("opening discord session: %w", err)
	}

	registered, err := dg.ApplicationCommandCreate(dg.State.User.ID, cfg.GuildID, &discordgo.ApplicationCommand{
		Name:        "request-book",
		Description: "Search for and download an eBook to the server",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Title and/or author to search for",
				Required:    true,
			},
		},
	})
	if err != nil {
		dg.Close()
		broker.Disconnect()
		return fmt.Errorf("registering /request-book command: %w", err)
	}
	if cfg.GuildID != "" {
		log.Printf("Discord: registered /request-book command in guild %s", cfg.GuildID)
	} else {
		log.Printf("Discord: registered /request-book command globally (may take up to an hour to appear)")
	}

	go b.reapSessions(ctx)
	log.Println("Discord bot is running.")

	<-ctx.Done()

	log.Println("Discord: shutting down…")
	if err := dg.ApplicationCommandDelete(dg.State.User.ID, cfg.GuildID, registered.ID); err != nil {
		log.Printf("Discord: failed to remove command: %v", err)
	}
	dg.Close()
	broker.Disconnect()
	return nil
}

func (b *bot) onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if i.ApplicationCommandData().Name == "request-book" {
			b.handleRequest(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if strings.HasPrefix(i.MessageComponentData().CustomID, selectPrefix) {
			b.handleSelect(s, i)
		}
	}
}

// handleRequest runs a search and presents the results as a select menu.
func (b *bot) handleRequest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	query := i.ApplicationCommandData().Options[0].StringValue()

	// Acknowledge immediately; searching can take a while.
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		log.Printf("defer request: %v", err)
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), searchTimeout+30*time.Second)
		defer cancel()

		books, err := b.broker.Search(ctx, query)
		if err != nil {
			b.editContent(s, i, fmt.Sprintf("❌ Search failed: %v", err))
			return
		}

		total := len(books)
		// Discord select menus allow at most 25 options.
		if len(books) > 25 {
			books = books[:25]
		}

		key := uuid.NewString()
		b.putSession(key, books)

		options := make([]discordgo.SelectMenuOption, len(books))
		for idx, bk := range books {
			options[idx] = discordgo.SelectMenuOption{
				Label:       truncate(displayTitle(bk), 100),
				Description: truncate(describe(bk), 100),
				Value:       strconv.Itoa(idx),
			}
		}

		content := fmt.Sprintf("Found **%d** results for *%s*.", total, query)
		if total > len(books) {
			content += fmt.Sprintf(" Showing the first %d — choose one to download:", len(books))
		} else {
			content += " Choose one to download:"
		}

		components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    selectPrefix + key,
						Placeholder: "Select a book…",
						Options:     options,
					},
				},
			},
		}

		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content:    &content,
			Components: &components,
		}); err != nil {
			log.Printf("edit search response: %v", err)
		}
	}()
}

// handleSelect downloads the chosen book and delivers it.
func (b *bot) handleSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	key := strings.TrimPrefix(data.CustomID, selectPrefix)

	books, ok := b.takeSession(key)
	if !ok {
		b.updateMessage(s, i, "⌛ This selection has expired. Please run `/request-book` again.")
		return
	}
	if len(data.Values) == 0 {
		b.updateMessage(s, i, "❌ No selection received.")
		return
	}
	idx, err := strconv.Atoi(data.Values[0])
	if err != nil || idx < 0 || idx >= len(books) {
		b.updateMessage(s, i, "❌ Invalid selection.")
		return
	}
	book := books[idx]

	// Replace the menu with a status message.
	b.updateMessage(s, i, fmt.Sprintf("⬇️ Downloading **%s**…", truncate(displayTitle(book), 200)))

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout+30*time.Second)
		defer cancel()

		path, err := b.broker.Download(ctx, book.Full)
		if err != nil {
			b.followup(s, i, fmt.Sprintf("❌ Download failed: %v", err))
			return
		}

		name := filepath.Base(path)
		absPath := path
		if a, absErr := filepath.Abs(path); absErr == nil {
			absPath = a
		}
		log.Printf("Saved %q to %s", name, absPath)

		// Saved to the server; also upload the file to Discord. If the upload is
		// rejected (e.g. too large, missing permissions), fall back to a message
		// pointing at the on-disk copy.
		if err := b.followupWithFile(s, i, name, path); err != nil {
			log.Printf("uploading %q to Discord: %v", name, err)
			b.followup(s, i, fmt.Sprintf("✅ Saved **%s** on the server at:\n`%s`\n⚠️ Couldn't upload it to Discord: %v", name, absPath, err))
		}
	}()
}

// editContent replaces the content of a deferred interaction response.
func (b *bot) editContent(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &content}); err != nil {
		log.Printf("edit content: %v", err)
	}
}

// updateMessage updates the message a component is attached to, clearing its
// components.
func (b *bot) updateMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: []discordgo.MessageComponent{},
		},
	}); err != nil {
		log.Printf("update message: %v", err)
	}
}

// followup posts a plain follow-up message to an interaction.
func (b *bot) followup(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{Content: content}); err != nil {
		log.Printf("followup message: %v", err)
	}
}

// followupWithFile uploads the file at path as a follow-up attachment. It
// returns an error (rather than only logging) so the caller can fall back to a
// path-only message when the upload is rejected.
func (b *bot) followupWithFile(s *discordgo.Session, i *discordgo.InteractionCreate, name, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("✅ Downloaded **%s** (also saved on the server).", name),
		Files:   []*discordgo.File{{Name: name, Reader: f}},
	})
	return err
}

func (b *bot) putSession(key string, books []core.BookDetail) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.sessions[key] = &session{books: books, created: time.Now()}
}

func (b *bot) takeSession(key string) ([]core.BookDetail, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, ok := b.sessions[key]
	if !ok {
		return nil, false
	}
	delete(b.sessions, key)
	return s.books, true
}

func (b *bot) reapSessions(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.mu.Lock()
			for k, s := range b.sessions {
				if time.Since(s.created) > sessionTTL {
					delete(b.sessions, k)
				}
			}
			b.mu.Unlock()
		}
	}
}

// displayTitle returns a human-friendly label for a result, falling back to the
// raw line when the parser couldn't extract a title.
func displayTitle(b core.BookDetail) string {
	if b.Title != "" {
		if b.Author != "" {
			return fmt.Sprintf("%s - %s", b.Author, b.Title)
		}
		return b.Title
	}
	return strings.TrimPrefix(b.Full, "!")
}

// describe builds the select-option description line (format/size/server).
func describe(b core.BookDetail) string {
	parts := make([]string, 0, 3)
	if b.Format != "" {
		parts = append(parts, b.Format)
	}
	if b.Size != "" {
		parts = append(parts, b.Size)
	}
	if b.Server != "" {
		parts = append(parts, "@"+b.Server)
	}
	return strings.Join(parts, " · ")
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 1 {
		return string(r[:max])
	}
	return string(r[:max-1]) + "…"
}
