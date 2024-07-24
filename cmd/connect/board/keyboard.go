package board

import (
	"github.com/gdamore/tcell/v2"
)

// pollEvents starts a goroutine to handle terminal events.
func (b *Board) pollEvents() chan struct{} {
	quit := make(chan struct{})

	go func() {
		for {
			event := b.screen.PollEvent()

			// Check if we received a key event.
			ev, isEventKey := event.(*tcell.EventKey)
			if !isEventKey {
				continue
			}

			keyType := ev.Key()

			// Allow the user to quit the game at any time.
			if keyType == tcell.KeyRune {
				if ev.Rune() == rune('q') {
					close(quit)
					return
				}
			}

			// Allow the user to clear the modal.
			if keyType == tcell.KeyEscape {
				if b.modalUp {
					b.closeModal()
				}
			}

			// Only the blue player can control the piece.
			if b.currentTurn == colorRed {
				b.screen.Beep()
				continue
			}

			switch keyType {
			case tcell.KeyRune:
				switch ev.Rune() {
				case rune('n'):
					b.newGame()
					if b.currentTurn == colorRed {
						b.runAISupport()
					}

				case rune(' '):
					if isWinner := b.dropPiece(true); !isWinner {
						b.runAISupport()
					}
				}

			case tcell.KeyLeft:
				b.movePlayerPiece(dirLeft)

			case tcell.KeyRight:
				b.movePlayerPiece(dirRight)

			case tcell.KeyEnter, tcell.KeyDown:
				if isWinner := b.dropPiece(true); !isWinner {
					b.runAISupport()
				}
			}
		}
	}()

	if b.currentTurn == colorRed {
		go b.runAISupport()
	}

	return quit
}
