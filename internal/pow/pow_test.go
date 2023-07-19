package pow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladimirpekarski/wordofwisdom/internal/env"
	"github.com/vladimirpekarski/wordofwisdom/internal/lib/logger"
	"github.com/vladimirpekarski/wordofwisdom/internal/message"
)

func TestPow_GenerateChallenge(t *testing.T) {
	tests := []struct {
		name         string
		genStrFunc   func(int) (string, error)
		difficulty   int
		expErr       bool
		expChallenge message.Challenge
	}{
		{
			name: "generate challenge",
			genStrFunc: func(i int) (string, error) {
				return "abc", nil
			},
			difficulty: 4,
			expChallenge: message.Challenge{
				RandomStr:  "abc",
				HashPrefix: "0000",
			},
		},
		{
			name: "error from genStrFunc",
			genStrFunc: func(i int) (string, error) {
				return "", errors.New("some error")
			},
			difficulty: 4,
			expChallenge: message.Challenge{
				RandomStr:  "abc",
				HashPrefix: "0000",
			},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			powMock := NewMock(logger.New(env.Local))
			powMock.genStrFunc = test.genStrFunc

			ch, err := powMock.GenerateChallenge(test.difficulty)

			if !test.expErr {
				assert.NoError(t, err)
				assert.Equal(t, test.expChallenge, ch)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestPow_Solve(t *testing.T) {
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name        string
		calcHashFun func(string) string
		ctx         context.Context
		challenge   message.Challenge
		expErr      bool
		expSolution message.Solution
	}{
		{
			name: "solve challenge",
			challenge: message.Challenge{
				RandomStr:  "a",
				HashPrefix: "0000",
			},
			calcHashFun: func(s string) string {
				if s == "a2" {
					return "0000x"
				}

				return "x"
			},
			ctx: context.Background(),
			expSolution: message.Solution{
				Hash:  "0000x",
				Nonce: 2,
			},
		},
		{
			name:      "cancelled context",
			challenge: message.Challenge{},
			calcHashFun: func(s string) string {
				return "x"
			},
			ctx:    canceledCtx,
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			powMock := NewMock(logger.New(env.Local))
			powMock.calcHashFun = test.calcHashFun

			sol, err := powMock.Solve(test.ctx, test.challenge)

			if !test.expErr {
				assert.NoError(t, err)
				assert.Equal(t, test.expSolution, sol)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestPow_Validate(t *testing.T) {
	tests := []struct {
		name        string
		calcHashFun func(string) string
		challenge   message.Challenge
		solution    message.Solution
		expRes      bool
	}{
		{
			name: "validate, true",
			calcHashFun: func(s string) string {
				return "000ax"
			},
			challenge: message.Challenge{
				RandomStr:  "rand",
				HashPrefix: "000",
			},
			solution: message.Solution{
				Hash:  "000ax",
				Nonce: 2,
			},
			expRes: true,
		},
		{
			name: "validate, false",
			calcHashFun: func(s string) string {
				return "00ax"
			},
			challenge: message.Challenge{
				RandomStr:  "rand",
				HashPrefix: "000",
			},
			solution: message.Solution{
				Hash:  "000ax",
				Nonce: 2,
			},
			expRes: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			powMock := NewMock(logger.New(env.Local))
			powMock.calcHashFun = test.calcHashFun

			valid := powMock.Validate(test.challenge, test.solution)

			assert.Equal(t, valid, test.expRes)
		})
	}
}
