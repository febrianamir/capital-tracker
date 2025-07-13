package app

import (
	"capital-tracker/lib"
	"capital-tracker/lib/constant"
	"capital-tracker/lib/style"
	"capital-tracker/response"
	"fmt"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (h *Handler) Update_ListWatchlist(app *App, msg tea.KeyMsg) (cmds []tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		app.Screen = constant.ModeMenu

	case tea.KeyEnter:
		cmds = h.startListWatchlistApi(app)
	}

	return cmds
}

func (h *Handler) View_ListWatchlist(app *App) string {
	if app.ListWatchlist.IsLoading {
		return fmt.Sprintf("\n  %s Loading...\n\n", app.Spinner.View())
	}

	var builder strings.Builder

	// each item in renderedContents represent one line
	renderedContents := []string{}

	mapCoinResponse := map[string]float64{}
	for _, coin := range app.ListWatchlist.CoinListResponse {
		mapCoinResponse[coin.ID] = coin.CurrentPrice
	}

	headerFormat := "%-20s %-16s"
	header := fmt.Sprintf(headerFormat, "Name", "Market Price")
	renderedContents = append(renderedContents, style.FontBold(style.ColorCyan(header)))

	dataFormat := "%-20s %-16s"

	for _, token := range app.ListWatchlist.Watchlists {
		renderedContents = append(renderedContents, fmt.Sprintf(dataFormat, fmt.Sprintf("%s (%s)", token.DisplayName, token.Symbol), fmt.Sprintf("$%s", lib.FormatPrice("%g", mapCoinResponse[token.ApiId]))))
	}

	for _, renderedContent := range renderedContents {
		builder.WriteString(lib.PrintLine(renderedContent))
	}

	return builder.String()
}

func (h *Handler) startListWatchlistApi(app *App) (cmds []tea.Cmd) {
	app.ListWatchlist.IsLoading = true

	tokenIds := []string{}
	for _, watchlist := range app.ListWatchlist.Watchlists {
		tokenIds = append(tokenIds, watchlist.ApiId)
	}
	queryParams := map[string]string{
		"vs_currency": "usd",
		"ids":         strings.Join(tokenIds, ","),
		"precision":   "2",
	}

	cmds = append(cmds, app.Spinner.Tick, func() tea.Msg {
		coins, err := lib.DoRequest[response.CoinList](http.MethodGet, "/coins/markets", queryParams)
		if err != nil {
			return AppResponseMsg{Error: err}
		}

		return AppResponseMsg{CoinListResponse: coins}
	})
	return cmds
}
