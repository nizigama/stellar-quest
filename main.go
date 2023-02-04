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
	assetCode                 string = "WEIRTKN"
	trustLimit                string = "18"
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

func createSignedTransaction(issuerKeypair *keypair.Full, trustingAccount horizon.Account) (*txnbuild.Transaction, error) {

	asset, _ := txnbuild.CreditAsset{
		Code:   assetCode,
		Issuer: issuerKeypair.Address(),
	}.ToChangeTrustAsset()

	op := txnbuild.ChangeTrust{
		Line:  asset,
		Limit: trustLimit,
	}

	sentAsset, err := asset.ToAsset()

	if err != nil {
		return nil, err
	}

	sendAssetsOp := txnbuild.Payment{
		Destination:   trustingAccount.AccountID,
		Asset:         sentAsset,
		Amount:        "18",
		SourceAccount: issuerKeypair.Address(),
	}

	createAccountOp := txnbuild.CreateAccount{
		Destination: issuerKeypair.Address(),
		Amount:      minimumCreateAccountTxFee,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &trustingAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&createAccountOp,
				&op,
				&sendAssetsOp,
			},
			BaseFee:       txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()}, // Use a real timeout in production!
		},
	)

	if err != nil {
		return nil, err
	}

	signer := getTrustingAccountKeypair()

	return tx.Sign(network.TestNetworkPassphrase, signer, issuerKeypair)
}

func main() {

	issuer := createIssuerAccountKeypair()

	trustingAcc, err := getTrustingAccountDetails()

	if err != nil {
		log.Println("error getting trusting account details", err)
		return
	}

	transaction, err := createSignedTransaction(issuer, trustingAcc)

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
