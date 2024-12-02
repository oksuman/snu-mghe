package mkbfv

import (
	"strconv"
	"testing"

	"github.com/oksuman/snu-mghe/mkrlwe"
)

func BenchmarkMKBFV(b *testing.B) {

	defaultParams := []ParametersLiteral{PN14QP439, PN15QP880}

	for _, defaultParam := range defaultParams {
		params := NewParametersFromLiteral(defaultParam)

		userList := make([]string, *maxUsers)
		idset := mkrlwe.NewIDSet()

		for i := range userList {
			userList[i] = "user" + strconv.Itoa(i)
			idset.Add(userList[i])
		}

		testContext, _ := genTestParams(params, idset)

		for numUsers := 2; numUsers <= *maxUsers; numUsers *= 2 {
			benchMulAndRelin(testContext, userList[:numUsers], b)
			benchPrevMulAndRelin(testContext, userList[:numUsers], b)
		}

		/*
			for numUsers := 2; numUsers <= maxUsers; numUsers *= 2 {
				benchRotate(testContext, userList[:numUsers], b)
			}
		*/
	}
}

func benchMulAndRelin(testContext *testParams, userList []string, b *testing.B) {

	numUsers := len(userList)
	msgList := make([]*Message, numUsers)
	ctList := make([]*Ciphertext, numUsers)

	rlkSet := testContext.rlkSet
	eval := testContext.evaluator

	for i := range userList {
		msgList[i], ctList[i] = newTestVectors(testContext, userList[i], 0, 2)
	}

	ct0 := ctList[0]
	ct1 := ctList[0]

	for i := range userList {
		ct0 = eval.AddNew(ct0, ctList[i])
		ct1 = eval.SubNew(ct1, ctList[i])
	}

	b.Run(GetTestName(testContext.params, "MKMulAndRelin: "+strconv.Itoa(numUsers)+"/ "), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			eval.MulRelinNew(ct0, ct1, rlkSet)
		}

	})
}

func benchPrevMulAndRelin(testContext *testParams, userList []string, b *testing.B) {

	numUsers := len(userList)
	msgList := make([]*Message, numUsers)
	ctList := make([]*Ciphertext, numUsers)

	rlkSet := testContext.oldRlkSet
	eval := testContext.evaluator

	for i := range userList {
		msgList[i], ctList[i] = newTestVectors(testContext, userList[i], 0, 2)
	}

	ct0 := ctList[0]
	ct1 := ctList[0]

	for i := range userList {
		ct0 = eval.AddNew(ct0, ctList[i])
		ct1 = eval.SubNew(ct1, ctList[i])
	}

	b.Run(GetTestName(testContext.params, "MKPrevMulAndRelin: "+strconv.Itoa(numUsers)+"/ "), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			eval.PrevMulRelinNew(ct0, ct1, rlkSet)
		}

	})
}

func benchRotate(testContext *testParams, userList []string, b *testing.B) {

	numUsers := len(userList)
	msgList := make([]*Message, numUsers)
	ctList := make([]*Ciphertext, numUsers)

	rtkSet := testContext.rtkSet
	eval := testContext.evaluator

	for i := range userList {
		msgList[i], ctList[i] = newTestVectors(testContext, userList[i], 0, 2)
	}

	ct := ctList[0]

	for i := range userList {
		ct = eval.AddNew(ct, ctList[i])
	}

	b.Run(GetTestName(testContext.params, "MKRotate: "+strconv.Itoa(numUsers)+"/ "), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			eval.RotateNew(ct, 2, rtkSet)
		}

	})
}
