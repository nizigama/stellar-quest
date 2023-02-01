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
	originPublicKey                    string = "GAKN4XJFTOFQTJS52CMTWWYCRJOIMBDAUTAAVQQF43WOIL3MU6PR3OG4"
	originSecretKey                    string = "SAJKZNLAHAT6ENZXYFGWJ6VWQJUUMB3KVR6VZFSR7FZUHRVY5IWWGMLM"
	networkPassphrase                  string = network.TestNetworkPassphrase
	minimumAmountForCreatingNewAccount string = "10"
)

var (
	networkClient *horizonclient.Client = horizonclient.DefaultTestNetClient
)

func getOriginAccount() (horizon.Account, error) {
	kp, _ := keypair.ParseFull(originSecretKey)

	accountDetailsRequest := horizonclient.AccountRequest{AccountID: kp.Address()}

	return networkClient.AccountDetail(accountDetailsRequest)
}

func getOriginKeypair() *keypair.Full {
	kp, _ := keypair.ParseFull(originSecretKey)

	return kp
}

func generateNewAccountKeypair() *keypair.Full {
	kp, _ := keypair.Random()

	return kp
}

func createSignedTransaction(originAcc horizon.Account, destinationAccountAddr string, amount string) (*txnbuild.Transaction, error) {

	paymentOperation := txnbuild.Payment{
		Destination: destinationAccountAddr,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}

	createAccountOperation := txnbuild.CreateAccount{
		Destination: destinationAccountAddr,
		Amount:      minimumAmountForCreatingNewAccount,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &originAcc,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&createAccountOperation,
				&paymentOperation,
			},
			BaseFee:       txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()}, // Use a real timeout in production!
		},
	)

	if err != nil {
		return nil, err
	}

	signer := getOriginKeypair()

	tx, err = tx.Sign(networkPassphrase, signer)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func main() {

	senderAcc, err := getOriginAccount()

	if err != nil {
		log.Println("error getting sender account details", err)
		return
	}

	receiver := generateNewAccountKeypair()

	amount := "450"

	transaction, err := createSignedTransaction(senderAcc, receiver.Address(), amount)

	if err != nil {
		log.Println("error creating payment transaction", err)
		return
	}

	tx, err := networkClient.SubmitTransaction(transaction)

	if err != nil {
		log.Println("error submit payment transaction", err)
		return
	}

	fmt.Printf("Payment transaction successful\n%v successfully sent to %s\nTransaction fees: %v\n", amount, receiver.Address(), tx.FeeCharged)
}
