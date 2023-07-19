package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/vladimirpekarski/wordofwisdom/internal/message"
	"golang.org/x/exp/rand"
	"golang.org/x/exp/slog"
	"strings"
	"time"
)

type Pow struct {
	Log *slog.Logger
}

func New(log *slog.Logger) Pow {
	return Pow{
		Log: log,
	}
}

func (p Pow) GenerateChallenge(n, difficulty int) (message.Challenge, error) {
	randomStr, err := p.generateString(difficulty)
	if err != nil {
		return message.Challenge{}, fmt.Errorf("failed to generate string: %w", err)
	}

	hasPrefix := strings.Repeat("0", difficulty)
	ch := message.Challenge{
		RandomStr:  randomStr,
		HashPrefix: hasPrefix,
	}

	return ch, nil
}

func (p Pow) Solve(ch message.Challenge) message.Solution {
	start := time.Now()
	nonce := 0

	for {
		hash := p.calcHash(fmt.Sprintf("%s%d", ch.RandomStr, nonce))
		if strings.HasPrefix(hash, ch.HashPrefix) {
			p.Log.Info("solved",
				slog.Int("elapsed_time, ms", int(time.Since(start).Milliseconds())),
				slog.Int("nonce", nonce))
			return message.Solution{
				Hash:  hash,
				Nonce: nonce,
			}
		}
		p.Log.Debug("wrong solving", slog.String("wrong_hash", hash))
		nonce++
	}
}

func (p Pow) Validate(ch message.Challenge, sl message.Solution) bool {
	expectedHash := p.calcHash(fmt.Sprintf("%s%d", ch.RandomStr, sl.Nonce))
	return strings.HasPrefix(expectedHash, ch.HashPrefix) && expectedHash == sl.Hash
}

func (p Pow) generateString(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", fmt.Errorf("failed to generate string: %w", err)
	}

	return hex.EncodeToString(bytes), err
}

func (p Pow) calcHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}