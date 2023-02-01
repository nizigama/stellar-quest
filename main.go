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
	parentPublicKey                  string = "GAWQLJVH2RNNNDRWQLXAF5SC6IDB67QUNPNJVTQ3K272GOIEIH3AQMNM"
	parentSecretKey                  string = "SDH4O75YQP2RKQ65XOTQCILDAOKT5NQSTEQZ2GPVF66Y7IYRCU6LD77V"
	newAccountMinimumStartingBalance string = "10"
	networkPassphrase                string = network.TestNetworkPassphrase
)

var (
	networkClient *horizonclient.Client = horizonclient.DefaultTestNetClient
)

func main() {

	parentAcc, err := getParentAccount()

	if err != nil {
		log.Println("error getting source account details", err)
		return
	}

	newAcc := generateNewAccountCredentials()

	signedTransaction, err := createNewAccountTransaction(parentAcc, newAcc.Address())

	if err != nil {
		log.Println("error creating new account transaction", err)
		return
	}

	tx, err := networkClient.SubmitTransaction(signedTransaction)

	if err != nil {
		log.Println("error submitting transaction", err)
		return
	}

	fmt.Printf("New account details:\nAddress: %s\nSecret: %s\nTransaction fees: %v\n", newAcc.Address(), newAcc.Seed(), tx.FeeCharged)
}

func getParentAccount() (horizon.Account, error) {

	kp, _ := keypair.ParseFull(parentSecretKey)

	request := horizonclient.AccountRequest{AccountID: kp.Address()}

	return networkClient.AccountDetail(request)
}

func getParentAccountKeypair() *keypair.Full {

	kp, _ := keypair.ParseFull(parentSecretKey)

	return kp
}

func generateNewAccountCredentials() *keypair.Full {
	kp, _ := keypair.Random()

	return kp
}

func createNewAccountTransaction(originAcc horizon.Account, destinationAccountAddr string) (*txnbuild.Transaction, error) {

	operation := txnbuild.CreateAccount{
		Destination: destinationAccountAddr,
		Amount:      newAccountMinimumStartingBalance,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &originAcc,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&operation},
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()}, // Use a real timeout in production!
		},
	)

	if err != nil {
		return nil, err
	}

	originKp := getParentAccountKeypair()

	return tx.Sign(networkPassphrase, originKp)
}
