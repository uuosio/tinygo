package main

import (
	"encoding/hex"

	"github.com/uuosio/chain"

	"github.com/drand/kyber"
	"github.com/drand/kyber/pairing/bn256"
	"github.com/drand/kyber/sign/bls"
)

func Debug(msg string) {
	chain.Println(msg)
}

func Sign(scalar kyber.Scalar, msg []byte) ([]byte, error) {
	scheme := bls.NewSchemeOnG1(bn256.NewSuiteG2())
	return scheme.Sign(scalar, msg)
}

func PrivateKeyFromHex(s string) (kyber.Scalar, error) {
	seed, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	suite := bn256.NewSuiteG2()
	scalar := suite.Scalar().SetBytes(seed)
	return scalar, nil
}

func Verify(pub kyber.Point, msg, sig []byte) error {
	scheme := bls.NewSchemeOnG1(bn256.NewSuiteG2())
	Debug("verifing...")
	//	return nil
	return scheme.Verify(pub, msg, sig)
}

func PublicKey(scalar kyber.Scalar) kyber.Point {
	suite := bn256.NewSuiteG2()
	return suite.Point().Mul(scalar, nil)
}

func main() {
	scalar, _ := PrivateKeyFromHex("2ccd7331e1f3d02da4df42302c14f9c16a6efdd1a7b6b31265b114083eb41b2d")
	{
		sig, err := Sign(scalar, []byte("hello"))
		if err != nil {
			panic(err)
		}
		chain.Println("++///:", sig)
	}
	return
}
