package utils

import (
	"encoding/hex"
	"math/rand"
	"mp-generate-test-data/internal/constants"
)

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandAddress(n int) string {
	return "0x" + RandStringRunes(n)
}

func RandTokenPayment() string {
	return constants.TOKEN_PAYMENT[rand.Intn(len(constants.TOKEN_PAYMENT))]
}

func RandVerifiedContract() string {
	return constants.VER_CONTRACT_ADDR[rand.Intn(len(constants.VER_CONTRACT_ADDR))]
}

// +1 => exclude 0
func RandNumberOfAssets() int {
	return rand.Intn(constants.MAX_NUM_ASSETS) + 1
}

func RandErcType() int {
	return constants.ERC_TYPE[rand.Intn(len(constants.ERC_TYPE))]
}

func RandQuantity() int {
	return rand.Intn(constants.MAX_QUANTITY)
}

func RandWithMinMax(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandPrice() float64 {
	return RandWithMinMax(1, 6)
}

// I want 9/10 order in orders was matched
func RandMatchedOrder() bool {
	return RandMatchedWithPreDesire(0.9)
}

func RandMatchedWithPreDesire(portion float64) bool {
	return portion >= rand.Float64()
}

func RandEndAt(startedAt int64) int64 {
	return startedAt + rand.Int63n(constants.DURATION_60_DAYS_SECOND)
}
