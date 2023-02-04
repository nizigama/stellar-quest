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
	masterAccountPublicKey string = "GBDFBHA3V27NAVHDUTPUXFBDXPSLJ3JAONQG3BGWL7YR42E4TMT3RYWA"
	masterAccountSecretKey string = "SASXFE3GMAIFZHE4YCCSOLOJRPES7HNS4XL6UICXHNLGEDSEP2KO4EH4"

	firstSignerAccountPublicKey string = "GBMI45L5XJHO3EHK4XGJIRQNST5AZXGPTT245HVP4KCD4EU5UQKFEYNT"
	firstSignerAccountSecretKey string = "SAHVM2C7HYRAQP6ZFBOIF7OELH7JGIOEYGEIRMROSICL3KHQVGETIJU5"

	secondSignerAccountPublicKey string = "GCQWKCEQ5OFCJHEYHXPFA52VI345F4ZHO4LRKSYER6XBKICGUBEI6ZDW"
	secondSignerAccountSecretKey string = "SAOBQIBJG2U43RHTVIKD3EWYZVCPHOP5MTRWXYPW26NATWS6M77FNMHL"

	receiverAccountPublicKey string = "GAGVO3IXPL5EZ7RXKQLT7MQSZ2AKEINIZVXJATQRM7F3EDHGEUWCERV5"
	receiverAccountSecretKey string = "SDPSSPG7TFK5776FMLSXUHOQOSNNA23I74TNLSM7HTLEQPUMMKO2U4QS"

	networkPassphrase string = network.TestNetworkPassphrase
)

var (
	networkClient *horizonclient.Client = horizonclient.DefaultTestNetClient
)

func getMasterAccount() (horizon.Account, error) {

	kp, _ := keypair.ParseFull(masterAccountSecretKey)

	request := horizonclient.AccountRequest{
		AccountID: kp.Address(),
	}

	return networkClient.AccountDetail(request)
}

func getMasterAccountKeypair() *keypair.Full {
	kp, _ := keypair.ParseFull(masterAccountSecretKey)

	return kp
}

func createControlSignedTransaction(account horizon.Account, accountKeypair *keypair.Full) (*txnbuild.Transaction, error) {

	thresholdOperation := &txnbuild.SetOptions{
		MasterWeight:    txnbuild.NewThreshold(txnbuild.Threshold(1)),
		LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(5)),
		MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(5)),
		HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(5)),
	}

	firstSignersOperation := &txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: firstSignerAccountPublicKey,
			Weight:  *txnbuild.NewThreshold(txnbuild.Threshold(2)),
		},
	}

	secondSignersOperation := &txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: secondSignerAccountPublicKey,
			Weight:  *txnbuild.NewThreshold(txnbuild.Threshold(2)),
		},
	}

	transaction, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations: []txnbuild.Operation{
			thresholdOperation,
			firstSignersOperation,
			secondSignersOperation,
		},
		BaseFee: txnbuild.MinBaseFee,
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewInfiniteTimeout(),
		},
	})

	if err != nil {
		return nil, err
	}

	return transaction.Sign(networkPassphrase, accountKeypair)
}

func createPaymentSignedTransaction(account horizon.Account, signersKeypairs ...*keypair.Full) (*txnbuild.Transaction, error) {

	transferOperation := &txnbuild.Payment{
		Destination: receiverAccountPublicKey,
		Amount:      "8000",
		Asset:       txnbuild.NativeAsset{},
	}

	transaction, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations: []txnbuild.Operation{
			transferOperation,
		},
		BaseFee: txnbuild.MinBaseFee,
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewInfiniteTimeout(),
		},
	})

	if err != nil {
		return nil, err
	}

	return transaction.Sign(networkPassphrase, signersKeypairs...)
}

func main() {

	account, err := getMasterAccount()

	if err != nil {
		log.Println("error getting master account details", err)
		return
	}

	masterKeypair := getMasterAccountKeypair()
	firstSignerKeypair := keypair.MustParseFull(firstSignerAccountSecretKey)
	secondSignerKeypair := keypair.MustParseFull(secondSignerAccountSecretKey)

	controlTransaction, err := createControlSignedTransaction(account, masterKeypair)

	if err != nil {
		log.Println("error creating a signed transaction for the control feature", err)
		return
	}

	_, err = networkClient.SubmitTransaction(controlTransaction)

	if err != nil {
		log.Println("error submitting control transaction to the network", err)
		return
	}

	fmt.Println("successfully set signers, weight and threshold")

	transferTransaction, err := createPaymentSignedTransaction(account, masterKeypair, firstSignerKeypair, secondSignerKeypair)

	if err != nil {
		log.Println("error creating a signed transaction for the transfer feature", err)
		return
	}

	_, err = networkClient.SubmitTransaction(transferTransaction)

	if err != nil {
		log.Println("error submitting transfer transaction to the network", err)
		return
	}

	fmt.Println("successfully transfered 8000 XLM")

}
