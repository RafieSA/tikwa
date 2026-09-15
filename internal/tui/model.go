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
	viewFilter
	viewRequestInput
	viewLoading
	viewSuccess
	viewError
	viewRequestSuccess
)

type Model struct {
	state       viewState
	cursor      int // menu utama 0-2
	filterCursor int // filter 0-4
	input       textinput.Model
	requestInput textinput.Model
	progress    progress.Model
	filePath    string
	profile     engine.Profile
	filter      engine.Filter
	errMsg      string
	successPath string
	requestPath string
	requestIdea string
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

	ri := textinput.New()
	ri.Placeholder = "misal: vintage BW 90an, warm sunset"
	ri.CharLimit = 200
	ri.Width = 50

	p := progress.New(progress.WithDefaultGradient(), progress.WithWidth(40))

	return Model{
		state:        viewMenu,
		input:        ti,
		requestInput: ri,
		progress:     p,
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

func doEnhance(path string, prof engine.Profile, filter engine.Filter) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		out := engine.OutputPath(path, prof)
		// sisipkan filter di nama file: video-tiktok-dramatis.mp4
		if filter != engine.FilterNatural && filter != engine.FilterOriginal {
			// ganti -tiktok.mp4 jadi -tiktok-dramatis.mp4
			out = strings.Replace(out, "-"+prof.Name+".", "-"+prof.Name+"-"+string(filter)+".", 1)
			if !strings.Contains(out, string(filter)) {
				// fallback kalau OutputPath tidak match
				ext := ".mp4"
				base := out[:len(out)-4]
				out = base + "-" + string(filter) + ext
			}
		}
		if filter == engine.FilterOriginal {
			out = strings.Replace(out, "-"+prof.Name+".", "-"+prof.Name+"-original.", 1)
		}
		err := engine.Run(ctx, path, out, prof, filter)
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
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state == viewMenu {
				return m, tea.Quit
			}
			if m.state == viewRequestSuccess || m.state == viewSuccess || m.state == viewError {
				m.state = viewMenu
				m.errMsg = ""
				return m, nil
			}
			// dari filter/request/loading, q balik menu
			if m.state == viewFilter || m.state == viewRequestInput {
				m.state = viewMenu
				m.errMsg = ""
				return m, nil
			}
		case "esc":
			switch m.state {
			case viewInput:
				m.state = viewMenu
				return m, nil
			case viewFilter:
				m.state = viewInput
				m.input.Focus()
				return m, textinput.Blink
			case viewRequestInput:
				m.state = viewFilter
				return m, nil
			case viewError, viewSuccess, viewRequestSuccess:
				m.state = viewMenu
				m.errMsg = ""
				return m, nil
			}
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
			case "enter":
				if m.cursor == 0 {
					m.state = viewInput
					m.profile = engine.TikTokProfile()
					m.input.Focus()
					return m, textinput.Blink
				}
				if m.cursor == 1 {
					m.state = viewInput
					m.profile = engine.WhatsAppProfile()
					m.input.Focus()
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
				m.state = viewFilter
				m.filterCursor = 0
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd

		case viewFilter:
			switch msg.String() {
			case "up", "k":
				if m.filterCursor > 0 {
					m.filterCursor--
				}
			case "down", "j":
				if m.filterCursor < 4 {
					m.filterCursor++
				}
			case "1", "2", "3", "4", "5":
				idx := int(msg.String()[0] - '1')
				m.filterCursor = idx
				return m.handleFilterSelect()
			case "enter":
				return m.handleFilterSelect()
			}

		case viewRequestInput:
			switch msg.String() {
			case "enter":
				idea := strings.TrimSpace(m.requestInput.Value())
				if idea == "" {
					m.errMsg = "❌ Ide tidak boleh kosong — tulis misal: vintage BW 90an"
					m.state = viewError
					return m, nil
				}
				path, err := engine.SaveRequest(idea)
				if err != nil {
					m.errMsg = "❌ " + err.Error()
					m.state = viewError
					return m, nil
				}
				m.requestPath = path
				m.requestIdea = idea
				m.requestInput.SetValue("")
				m.state = viewRequestSuccess
				return m, nil
			}
			var cmd tea.Cmd
			m.requestInput, cmd = m.requestInput.Update(msg)
			return m, cmd

		case viewError, viewSuccess, viewRequestSuccess:
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

func (m Model) handleFilterSelect() (tea.Model, tea.Cmd) {
	filters := engine.AllFilters()
	f := filters[m.filterCursor]
	if f == engine.FilterRequest {
		m.state = viewRequestInput
		m.requestInput.Focus()
		return m, textinput.Blink
	}
	m.filter = f
	m.state = viewLoading
	m.percent = 0
	return m, tea.Batch(tick(), doEnhance(m.filePath, m.profile, m.filter))
}

func (m Model) View() string {
	logo := StyleLogo.Render("▶  TIKWA  ◀")
	sub := StyleSubtitle.Render("Bikin Video HP Tidak Burem  •  TikTok • WhatsApp")
	header := lipgloss.JoinVertical(lipgloss.Center, logo, sub)
	header = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, header)

	switch m.state {
	case viewMenu:
		return m.viewMenu(header)
	case viewInput:
		return m.viewInput(header)
	case viewFilter:
		return m.viewFilter(header)
	case viewRequestInput:
		return m.viewRequestInput(header)
	case viewLoading:
		return m.viewLoading(header)
	case viewSuccess:
		return m.viewSuccess(header)
	case viewRequestSuccess:
		return m.viewRequestSuccess(header)
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
			StyleHelp.Render("Enter = lanjut pilih filter  •  Esc = kembali"),
	)
	box = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, box)
	b.WriteString(box)

	if strings.TrimSpace(m.input.Value()) == "" {
		empty := StyleEmptyBox.Render("📭 Contoh: /Users/rafie/Movies/VID_20260915.mp4\natau drag & drop file ke terminal")
		empty = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, empty)
		b.WriteString("\n" + empty)
	}
	return b.String()
}

func (m Model) viewFilter(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")

	title := StyleSubtitle.Render(fmt.Sprintf("📁 %s  •  %s (%dx%d)", m.filePath, strings.ToUpper(m.profile.Name), m.profile.Width, m.profile.Height))
	title = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title)
	b.WriteString(title + "\n")

	prompt := lipgloss.NewStyle().Bold(true).Foreground(ColorGold).Render("Pilih Filter:")
	prompt = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, prompt)
	b.WriteString(prompt + "\n")

	filters := engine.AllFilters()
	var rows []string
	for i, f := range filters {
		key := StyleMenuKey.Render(fmt.Sprintf("[%d]", i+1))
		lbl := StyleMenuItem.Render(f.Label())
		row := lipgloss.JoinHorizontal(lipgloss.Center, key, lbl)
		if i == m.filterCursor {
			row = StyleSelected.Render(" " + row + " ")
		} else {
			row = "  " + row + "  "
		}
		rows = append(rows, row)
	}
	menu := lipgloss.JoinVertical(lipgloss.Left, rows...)
	box := StyleMenuBox.Render(menu)
	box = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, box)
	b.WriteString(box + "\n")

	help := StyleHelp.Render("↑↓ pilih  •  1-5 atau Enter  •  Esc kembali  •  q menu")
	help = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, help)
	b.WriteString(help)
	return b.String()
}

func (m Model) viewRequestInput(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")

	title := lipgloss.NewStyle().Bold(true).Foreground(ColorGold).Render("💡 Request Filter Impian")
	title = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title)
	b.WriteString(title + "\n")

	box := StyleMenuBox.Render(
		"Tulis ide filter lu (contoh: vintage BW 90an, warm sunset):\n" + m.requestInput.View() + "\n\n" +
			StyleHelp.Render("Enter = kirim  •  Esc = kembali"),
	)
	box = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, box)
	b.WriteString(box)

	hint := StyleEmptyBox.Render("📝 Ide lu akan disimpan di ~/.tikwa-requests.txt\n& bisa lu kirim juga ke GitHub Issues:\nhttps://github.com/RafieSA/tikwa/issues/new")
	hint = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, hint)
	b.WriteString("\n" + hint)
	return b.String()
}

func (m Model) viewLoading(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")
	info := StyleSubtitle.Render(fmt.Sprintf("⏳ Memoles: %s → %s [%s]", m.filePath, engine.OutputPath(m.filePath, m.profile), m.filter))
	info = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, info)
	b.WriteString(info + "\n\n")

	bar := m.progress.ViewAs(m.percent)
	bar = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, bar)
	b.WriteString(bar + "\n")

	pct := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, StyleHelp.Render(fmt.Sprintf("%.0f%% • filter: %s • mohon tunggu...", m.percent*100, m.filter)))
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

func (m Model) viewRequestSuccess(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")
	msg := StyleSuccessBox.Render(fmt.Sprintf("💡 Request dicatat: \"%s\"", m.requestIdea))
	msg = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, msg)
	b.WriteString(msg + "\n")
	info := StyleSubtitle.Render(fmt.Sprintf("Disimpan di: %s", m.requestPath))
	info = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, info)
	b.WriteString(info + "\n\n")

	link := StyleMenuBox.Render("Kirim juga ke GitHub biar gue bikin:\nhttps://github.com/RafieSA/tikwa/issues/new?title=[Filter Request] " + m.requestIdea)
	link = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, link)
	b.WriteString(link + "\n")

	help := StyleHelp.Render("Enter = kembali ke menu  •  q = keluar")
	help = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, help)
	b.WriteString(help)
	return b.String()
}

func (m Model) viewError(header string) string {
	var b strings.Builder
	b.WriteString(header + "\n")
	errText := m.errMsg
	if len(errText) > 400 {
		errText = errText[:400] + "…"
	}
	msg := StyleErrorBox.Render(errText)
	msg = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, fmt.Sprintf("%s", msg))
	b.WriteString(msg + "\n\n")
	help := StyleHelp.Render("Enter/Esc = kembali  •  cek path & coba lagi")
	help = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, help)
	b.WriteString(help)
	return b.String()
}
