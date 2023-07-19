package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/exp/rand"
	"golang.org/x/exp/slog"
	"strings"
	"time"
)

type Pow struct {
	Log *slog.Logger
}

type Challenge struct {
	RandomStr  string
	HashPrefix string
}

type Solution struct {
	Hash  string
	Nonce int
}

type Record struct {
	Quote  string
	Author string
}

func New(log *slog.Logger) Pow {
	return Pow{
		Log: log,
	}
}

func (p Pow) GenerateChallenge(n, difficulty int) (Challenge, error) {
	randomStr, err := p.generateString(difficulty)
	if err != nil {
		return Challenge{}, fmt.Errorf("failed to generate string: %w", err)
	}

	hasPrefix := strings.Repeat("0", difficulty)
	ch := Challenge{
		RandomStr:  randomStr,
		HashPrefix: hasPrefix,
	}

	return ch, nil
}

func (p Pow) Solve(ch Challenge) Solution {
	start := time.Now()
	nonce := 0

	for {
		hash := p.calcHash(fmt.Sprintf("%s%d", ch.RandomStr, nonce))
		if strings.HasPrefix(hash, ch.HashPrefix) {
			p.Log.Info("solved",
				slog.Int("elapsed_time, ms", int(time.Since(start).Milliseconds())),
				slog.Int("nonce", nonce))
			return Solution{
				Hash:  hash,
				Nonce: nonce,
			}
		}
		p.Log.Debug("wrong solving", slog.String("wrong_hash", hash))
		nonce++
	}
}

func (p Pow) Validate(ch Challenge, sl Solution) bool {
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
