package message

type Challenge struct {
	RandomStr  string
	HashPrefix string
}

type Solution struct {
	Hash  string
	Nonce int
}

type BookRecord struct {
	Quote            string
	Author           string
	PassedValidation bool
}
