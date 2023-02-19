package ltemodulecontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceSuccess(t *testing.T) {
	input := `OK

	+CUSD: 1,"Balans 84.59 grn, bonus 0.00.
	***
	Zberigayte pryvatnist' vashogo nomeru z Drugym Nomerom
	1.Tak
	",15`

	expected := 84.59
	actual, err := ParseAccBalance(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBalanceError(t *testing.T) {
	input := `OK

	+CUSD: 1,"Bans none bonus 0.00.
	***
	Zberigayte pryvatnist' vashogo nomeru z Drugym Nomerom
	1.Tak
	",15`

	expected := 0.0
	actual, err := ParseAccBalance(input)

	assert.Error(t, err)
	assert.Equal(t, expected, actual)
}

func TestInetBalanceSuccess(t *testing.T) {
	input := `	OK

+CUSD: 0,"300hv na vsi mobil'ni nomery, rezervnyy bezlimit na lifecell, 7.55 GB. Nastupna oplata paketu poslug 14.03.23.

",15`

	expected := 7.55
	actual, err := ParseInetBalance(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestInetBalanceError(t *testing.T) {
	input := `	OK

+CUSD: 0,"300hv na vsi mobil'ni nomery, rezervnyy bezlimit na lifecell, no-value. Nastupna oplata paketu poslug 14.03.23.

",15`

	expected := 0.0
	actual, err := ParseInetBalance(input)

	assert.Error(t, err)
	assert.Equal(t, expected, actual)
}

func TestTariffSuccess(t *testing.T) {
	input := `OK

	+CUSD: 1,"Vash nomer:380931272950
	Vash taryf:Prosto Life za 120grn/ 4 tyzhni
	1.Detali
	2.Pereviryty datu aktyvatsiyi
	3.Zavantazhuyte BiP
	 ",15`

	expected := "Prosto Life za 120grn/ 4 tyzhni"
	actual, err := ParseTariff(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTariffError(t *testing.T) {
	input := `OK

	+CUSD: 1,"Vash nomer:380931272950
	rosto Life za 120grn/ 4 tyzhni
	1.Detali
	2.Pereviryty datu aktyvatsiyi
	3.Zavantazhuyte BiP
	 ",15`

	expected := ""
	actual, err := ParseTariff(input)

	assert.Error(t, err)
	assert.Equal(t, expected, actual)
}

func TestPhoneNumberSuccess(t *testing.T) {
	input := `OK

	+CUSD: 1,"Vash nomer:380931272950
	Vash taryf:Prosto Life za 120grn/ 4 tyzhni
	1.Detali
	2.Pereviryty datu aktyvatsiyi
	3.Zavantazhuyte BiP
	 ",15`

	expected := "380931272950"
	actual, err := ParsePhoneNumber(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestPhoneNumberError(t *testing.T) {
	input := `OK

	+CUSD: 1,"Vash nomer:absent
	Vash taryf:Prosto Life za 120grn/ 4 tyzhni
	1.Detali
	2.Pereviryty datu aktyvatsiyi
	3.Zavantazhuyte BiP
	 ",15`

	expected := ""
	actual, err := ParsePhoneNumber(input)

	assert.Error(t, err)
	assert.Equal(t, expected, actual)
}
