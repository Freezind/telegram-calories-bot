package unit

import (
	"testing"

	"github.com/freezind/telegram-calories-bot/internal/services"
)

func TestGreet(t *testing.T) {
	got := services.Greet()
	want := "Hello! 👋"
	if got != want {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
}
