package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

const (
	trustingAccountPublicKey  string = "GD4RJTVWN55XW4SLC4GUBIMQVEXE6TV455UGIFTJH6MYJ5PGP7JLLJ6B"
	trustingAccountSecretKey  string = "SAWAVABGVSR2KZFZWRRNBI4RZ7LK7ORG5E3YXLR5KBFJDAQJ4JLJ6SLM"
	assetCode                 string = "TKN"
	trustLimit                string = "12"
	minimumCreateAccountTxFee string = "10"
)

var (
	networkClient *horizonclient.Client = horizonclient.DefaultTestNetClient
)

func createIssuerAccountKeypair() *keypair.Full {
	kp, _ := keypair.Random()

	return kp
}

func getTrustingAccountDetails() (horizon.Account, error) {

	kp, _ := keypair.ParseFull(trustingAccountSecretKey)

	ar := horizonclient.AccountRequest{AccountID: kp.Address()}

	return networkClient.AccountDetail(ar)
}

func getTrustingAccountKeypair() *keypair.Full {
	kp, _ := keypair.ParseFull(trustingAccountSecretKey)

	return kp
}

func createSignedTransaction(issuerAddr string, trustingAccount horizon.Account) (*txnbuild.Transaction, error) {

	asset, _ := txnbuild.CreditAsset{
		Code:   assetCode,
		Issuer: issuerAddr,
	}.ToChangeTrustAsset()

	op := txnbuild.ChangeTrust{
		Line:  asset,
		Limit: trustLimit,
	}

	createAccountOp := txnbuild.CreateAccount{
		Destination: issuerAddr,
		Amount:      minimumCreateAccountTxFee,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &trustingAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&createAccountOp,
				&op,
			},
			BaseFee:       txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()}, // Use a real timeout in production!
		},
	)

	if err != nil {
		return nil, err
	}

	signer := getTrustingAccountKeypair()

	return tx.Sign(network.TestNetworkPassphrase, signer)
}

func main() {

	issuer := createIssuerAccountKeypair()

	trustingAcc, err := getTrustingAccountDetails()

	if err != nil {
		log.Println("error getting trusting account details", err)
		return
	}

	transaction, err := createSignedTransaction(issuer.Address(), trustingAcc)

	if err != nil {
		log.Println("error creating change trust transaction", err)
		return
	}

	tx, err := networkClient.SubmitTransaction(transaction)

	if err != nil {
		log.Println("error submitting change trust transaction", horizonclient.GetError(err).Problem.Extras)
		return
	}

	fmt.Printf("Asset successfully created\nFees: %v\n", tx.FeeCharged)
}
