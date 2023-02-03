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
	mainAccountPublicKey       string = "GAS46PNBXBPV3WRO2IEQB5NUB7IQNAZEJ6MY4HDGZJJ24JEZWZ6TRDOT"
	mainAccountSecretKey       string = "SCCMBR2FZBFAZ7RZIA2762BXXOOVGRHYZOZJQOCVCJKTVA7JCVOHZFM3"
	networkPassphrase          string = network.TestNetworkPassphrase
	minimumCreateAccountAmount string = "10"
)

var (
	networkClient *horizonclient.Client = horizonclient.DefaultTestNetClient
)

func getMainAccountDetails() (horizon.Account, error) {

	kp, _ := keypair.Parse(mainAccountSecretKey)

	req := horizonclient.AccountRequest{
		AccountID: kp.Address(),
	}

	return networkClient.AccountDetail(req)
}

func getMainAccountKeypair() *keypair.Full {

	kp, _ := keypair.ParseFull(mainAccountSecretKey)

	return kp
}

func createNewRandomAccountKeypair() *keypair.Full {
	kp, _ := keypair.Random()

	return kp
}

func createSignedTransaction(newAccountAddr string, creatorAccountKp *keypair.Full, creatorAccount horizon.Account) (*txnbuild.Transaction, error) {

	createAccountOperation := txnbuild.CreateAccount{
		Destination: newAccountAddr,
		Amount:      minimumCreateAccountAmount,
	}

	mergeAccountOperation := txnbuild.AccountMerge{
		Destination:   newAccountAddr,
		SourceAccount: creatorAccountKp.Address(),
	}

	transaction, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &creatorAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&createAccountOperation,
				&mergeAccountOperation,
			},
			BaseFee:       txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()}, // Use a real timeout in production!
		},
	)

	if err != nil {
		return nil, err
	}

	return transaction.Sign(networkPassphrase, creatorAccountKp)
}

func main() {

	mainAcc, err := getMainAccountDetails()

	if err != nil {
		log.Println("failed getting the account details of the main account", err)
		return
	}

	mainAccKp := getMainAccountKeypair()

	newAcc := createNewRandomAccountKeypair()

	transaction, err := createSignedTransaction(newAcc.Address(), mainAccKp, mainAcc)

	if err != nil {
		log.Println("failed creating a signed transaction", err)
		return
	}

	tx, err := networkClient.SubmitTransaction(transaction)

	if err != nil {
		log.Println("failed submitting transaction", err)
		return
	}

	fmt.Printf("Account deletion successful\nTransaction fees: %v", tx.FeeCharged)
}
