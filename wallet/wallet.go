// Package wallet defines common datastructures and an interface for cryptocurrency wallets
package wallet

import (
	"github.com/brave-intl/bat-go/utils/altcurrency"
	"github.com/shopspring/decimal"
)

// Info contains information about a wallet like associated identifiers, the denomination,
// the last known balance and provider
type Info struct {
	ID          string                   `json:"paymentId" valid:"uuidv4,optional"`
	Provider    string                   `json:"provider" valid:"in(uphold)"`
	ProviderID  string                   `json:"providerId" valid:"uuidv4"`
	AltCurrency *altcurrency.AltCurrency `json:"altcurrency" valid:"-"`
	PublicKey   string                   `json:"publicKey,omitempty" valid:"hexadecimal,optional"`
	LastBalance *Balance                 `json:"balances,omitempty" valid:"-"`
}

// TransactionInfo contains information about a transaction like the denomination, amount in probi,
// destination address, status and identifier
type TransactionInfo struct {
	Probi       decimal.Decimal          `json:"probi"`
	AltCurrency *altcurrency.AltCurrency `json:"altcurrency"`
	Destination string                   `json:"address"`
	Fee         decimal.Decimal          `json:"fee"`
	Status      string                   `json:"status"`
	ID          string                   `json:"id"`
}

// Balance holds balance information for a wallet
type Balance struct {
	TotalProbi       decimal.Decimal
	SpendableProbi   decimal.Decimal
	ConfirmedProbi   decimal.Decimal
	UnconfirmedProbi decimal.Decimal
}

// Wallet is an interface for a cryptocurrency wallet
type Wallet interface {
	GetWalletInfo() Info
	// Transfer moves funds out of the associated wallet and to the specific destination
	Transfer(altcurrency altcurrency.AltCurrency, probi decimal.Decimal, destination string) (*TransactionInfo, error)
	// VerifyTransaction verifies that the base64 encoded transaction is valid
	// NOTE VerifyTransaction must guard against transactions that seek to exploit parser differences
	// such as including additional fields that are not understood by local implementation but may
	// be understood by the upstream wallet provider.
	VerifyTransaction(transactionB64 string) (*TransactionInfo, error)
	// SubmitTransaction submits the base64 encoded transaction for verification but does not move funds
	SubmitTransaction(transactionB64 string, confirm bool) (*TransactionInfo, error)
	// ConfirmTransaction confirms a previously submitted transaction, moving funds
	ConfirmTransaction(id string) (*TransactionInfo, error)
	// GetBalance returns the last known balance, if refresh is true then the current balance is fetched
	GetBalance(refresh bool) (*Balance, error)
}

// IsInsufficientBalance is a helper method for determining if an error indicates insufficient balance
// to perform the requested action
func IsInsufficientBalance(err error) bool {
	type insufficientBalance interface {
		InsufficientBalance() bool
	}
	te, ok := err.(insufficientBalance)
	return ok && te.InsufficientBalance()
}

// IsUnauthorized is a helper method for determining if an error indicates the wallet unauthorized
// to perform the requested action
func IsUnauthorized(err error) bool {
	type unauthorized interface {
		Unauthorized() bool
	}
	te, ok := err.(unauthorized)
	return ok && te.Unauthorized()
}

// IsInvalidSignature is a helper method for determining if an error indicates there was an invalid signature
func IsInvalidSignature(err error) bool {
	type invalidSignature interface {
		InvalidSignature() bool
	}
	te, ok := err.(invalidSignature)
	return ok && te.InvalidSignature()
}
