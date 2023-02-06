package logging

import "testing"

func TestLog_LevelDispatch_LogAll(t *testing.T) {
	l := NewLogger(Config{
		Level: LevelAll,
	})
	suite := []LogLevel{
		LevelFatal,
		LevelError,
		LevelWarning,
		LevelInfo,
		LevelDebug,
		LevelAll,
	}
	for _, lvl := range suite {
		t.Run(lvl.String(), func(t *testing.T) {
			if !l.log(lvl, "wat") {
				t.Errorf("%q should be logged with LevelAll", lvl)
			}
		})
	}
}

func TestLog_LevelDispatch_LogWarningsAndAbove(t *testing.T) {
	l := NewLogger(Config{
		Level: LevelWarning,
	})
	suite := map[LogLevel]bool{
		LevelFatal:   true,
		LevelError:   true,
		LevelWarning: true,
		LevelInfo:    false,
		LevelDebug:   false,
		LevelAll:     false,
		LevelTrace:   false,
	}
	for lvl, want := range suite {
		t.Run(lvl.String(), func(t *testing.T) {
			got := l.log(lvl, "wat")
			if want != got {
				t.Errorf("%q: want %v, got %v", lvl, want, got)
			}
		})
	}
}
