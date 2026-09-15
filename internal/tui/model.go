package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/RafieSA/tikwa/internal/engine"
)

// State halaman TUI
type viewState int

const (
	viewMenu viewState = iota
	viewInput
	viewLoading
	viewSuccess
	viewError
)

type Model struct {
	state       viewState
	cursor      int // 0=TikTok, 1=WhatsApp, 2=Keluar
	input       textinput.Model
	progress    progress.Model
	filePath    string
	profile     engine.Profile
	errMsg      string
	successPath string
	percent     float64
	width       int
	height      int
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "/Users/rafie/Videos/video.mp4"
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 50

	p := progress.New(progress.WithDefaultGradient(), progress.WithWidth(40))

	return Model{
		state:    viewMenu,
		input:    ti,
		progress: p,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Messages
type tickMsg time.Time
type doneMsg struct{ path string; err error }

func tick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func doEnhance(path string, prof engine.Profile) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		out := engine.OutputPath(path, prof)
		err := engine.Run(ctx, path, out, prof)
		if err != nil {
			return doneMsg{"", err}
		}
		return doneMsg{out, nil}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == viewMenu {
				return m, tea.Quit
			}
			// dari state lain, q = balik ke menu
			m.state = viewMenu
			m.errMsg = ""
			return m, nil
		case "esc":
			m.state = viewMenu
			m.errMsg = ""
			return m, nil
		}

		switch m.state {
		case viewMenu:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < 2 {
					m.cursor++
				}
			case "1":
				m.cursor = 0
				m.state = viewInput
				m.profile = engine.TikTokProfile()
				m.input.Focus()
				return m, textinput.Blink
			case "2":
				m.cursor = 1
				m.state = viewInput
				m.profile = engine.WhatsAppProfile()
				m.input.Focus()
				return m, textinput.Blink
			case "3", "enter":
				if m.cursor == 0 {
					m.cursor = 0
					m.state = viewInput
					m.profile = engine.TikTokProfile()
					return m, textinput.Blink
				}
				if m.cursor == 1 {
					m.cursor = 1
					m.state = viewInput
					m.profile = engine.WhatsAppProfile()
					return m, textinput.Blink
				}
				if m.cursor == 2 {
					return m, tea.Quit
				}
			}

		case viewInput:
			switch msg.String() {
			case "enter":
				path := strings.TrimSpace(m.input.Value())
				if path == "" {
					m.errMsg = "📭 Belum ada file dipilih — ketik path video dulu"
					m.state = viewError
					return m, nil
				}
				if err := engine.ValidateInput(path); err != nil {
					m.errMsg = "❌ " + err.Error()
					m.state = viewError
					return m, nil
				}
				if err := engine.CheckFFmpeg(); err != nil {
					m.errMsg = "❌ " + err.Error()
					m.state = viewError
					return m, nil
				}
				m.filePath = path
				m.state = viewLoading
				m.percent = 0
				return m, tea.Batch(tick(), doEnhance(path, m.profile))
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd

		case viewError, viewSuccess:
			if msg.String() == "enter" {
				m.state = viewMenu
				m.errMsg = ""
				return m, nil
			}
		}

	case tickMsg:
		if m.state == viewLoading {
			if m.percent < 0.92 {
				m.percent += 0.07
				return m, tick()
			}
			return m, tick()
		}

	case doneMsg:
		if msg.err != nil {
			m.errMsg = "❌ " + msg.err.Error()
			m.state = viewError
			return m, nil
		}
		m.successPath = msg.path
		m.percent = 1.0
		m.state = viewSuccess
		return m, nil
	}

	// update progress sub-model
	if m.state == viewLoading {
		cmd := m.progress.SetPercent(m.percent)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	// Logo — selalu tampil di atas
	logo := StyleLogo.Render("▶  TIKWA  ◀")
	sub := StyleSubtitle.Render("Bikin Video HP Tidak Burem  •  TikTok • WhatsApp")
	header := lipgloss.JoinVertical(lipgloss.Center, logo, sub)
	header = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, header)

	switch m.state {
	case viewMenu:
		return m.viewMenu(header)
	case viewInput:
		return m.viewInput(header)
	case viewLoading:
		return m.viewLoading(header)
	case viewSuccess:
		return m.viewSuccess(header)
	case viewError:
		return m.viewError(header)
	default:
		return header
	}
}

func (m Model) viewMenu(header string) string {
	items := []struct{ key, label string }{
		{"1", "Untuk TikTok  — 1080×1920 vertikal, tajam & terang"},
		{"2", "Untuk WhatsApp — 720p kecil, tetap jernih (<25MB)"},
		{"q", "Keluar"},
	}
	var b strings.Builder
	b.WriteString(header + "\n")

	var rows []string
	for i, it := range items {
		key := StyleMenuKey.Render("[" + it.key + "]")
		lbl := StyleMenuItem.Render(it.label)
		row := lipgloss.JoinHorizontal(lipgloss.Center, key, lbl)
		if i == m.cursor {
			row = StyleSelected.Render(" " + row + " ")
		} else {
			row = "  " + row + "  "
		}
		rows = append(rows, row)
	}
	menu := lipgloss.JoinVertical(lipgloss.Left, rows...)
	box := StyleMenuBox.Render(menu)
	box = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, box)

	help := StyleHelp.Render("↑↓ pilih  •  1/2 atau Enter  •  q keluar")
	help = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, help)

	b.WriteString(box + "\n" + help)
	return b.String()
}

func (m Model) viewInput(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")

	title := lipgloss.NewStyle().Bold(true).Foreground(ColorTikTokCyan).Render(
		fmt.Sprintf("→ Mode: %s (%dx%d)", strings.ToUpper(m.profile.Name), m.profile.Width, m.profile.Height),
	)
	title = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title)
	b.WriteString(title + "\n")

	box := StyleMenuBox.Render(
		"Path video:\n" + m.input.View() + "\n\n" +
			StyleHelp.Render("Enter = mulai poles  •  Esc = kembali"),
	)
	box = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, box)
	b.WriteString(box)

	// Empty state hint jika input masih kosong
	if strings.TrimSpace(m.input.Value()) == "" {
		empty := StyleEmptyBox.Render("📭 Contoh: /Users/rafie/Movies/VID_20260915.mp4\natau drag & drop file ke terminal")
		empty = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, empty)
		b.WriteString("\n" + empty)
	}
	return b.String()
}

func (m Model) viewLoading(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")
	info := StyleSubtitle.Render(fmt.Sprintf("⏳ Memoles: %s → %s", m.filePath, engine.OutputPath(m.filePath, m.profile)))
	info = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, info)
	b.WriteString(info + "\n\n")

	bar := m.progress.ViewAs(m.percent)
	bar = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, bar)
	b.WriteString(bar + "\n")

	pct := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, StyleHelp.Render(fmt.Sprintf("%.0f%% • mohon tunggu...", m.percent*100)))
	b.WriteString(pct)
	return b.String()
}

func (m Model) viewSuccess(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")
	msg := StyleSuccessBox.Render(fmt.Sprintf("✅ Berhasil! Hasil: %s", m.successPath))
	msg = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, msg)
	b.WriteString(msg + "\n\n")
	help := StyleHelp.Render("Enter = kembali ke menu  •  q = keluar")
	help = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, help)
	b.WriteString(help)
	return b.String()
}

func (m Model) viewError(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")
	// potong error biar tidak kepanjangan di TUI
	errText := m.errMsg
	if len(errText) > 300 {
		errText = errText[:300] + "…"
	}
	msg := StyleErrorBox.Render(errText)
	msg = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, fmt.Sprintf("%s", msg))
	b.WriteString(msg + "\n\n")
	help := StyleHelp.Render("Enter/Esc = kembali  •  cek path & coba lagi")
	help = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, help)
	b.WriteString(help)
	return b.String()
}
