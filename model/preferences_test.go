/*
 * Copyright (c) 2020, nwillc@gmail.com
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package model

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/nwillc/snipgo/mocks"
	"github.com/nwillc/snipgo/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

const testPrefFile = "../test/files/preferences.json"

type PreferencesTestSuite struct {
	suite.Suite
	ctx *services.Context
	badFilename  string
	goodFilename string
}

func TestPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesTestSuite))
}

func (suite *PreferencesTestSuite) SetupTest() {
	suite.ctx = services.NewDefaultContext()
	suite.badFilename = "foo"
	suite.goodFilename = testPrefFile
}

//func (suite *PreferencesTestSuite) TestBadHomeDir() {
//	defaultUserHomeGet := userHomeGet
//	userHomeGet = failingHomeDirMock{}
//	defer func() {
//		userHomeGet = defaultUserHomeGet
//		recover()
//	}()
//	_, _ = ReadPreferences("")
//	suite.T().Errorf("expected panic")
//}

func (suite *PreferencesTestSuite) TestNonExistPrefs() {
	_, ok := ReadPreferences(suite.ctx, suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestExistPrefs() {
	_, ok := ReadPreferences(suite.ctx, suite.goodFilename)
	assert.Nil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestMalformedFile() {
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	tempFile.WriteString("not json")

	_, ok := ReadPreferences(suite.ctx, tempFile.Name())
	if assert.NotNil(suite.T(), ok) {
		assert.Errorf(suite.T(), ok, "json marshal failed")
	}
}

func (suite *PreferencesTestSuite) TestWrite() {
	p := Preferences{DefaultFile: "foo"}
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	err = p.Write(suite.ctx, tempFile.Name())
	assert.Nil(suite.T(), err)
	read, err := ReadPreferences(suite.ctx, tempFile.Name())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), p.DefaultFile, read.DefaultFile)
}

func (suite *PreferencesTestSuite) TestWriteMarshal() {
	p := Preferences{DefaultFile: "foo"}
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockJson := mocks.NewMockJson(mockCtrl)
	mockJson.EXPECT().
		Marshal(gomock.Any()).
		Return([]byte{}, fmt.Errorf("mock error")).
		Times(1)
	var mockContext = services.Context{
		JSON:   mockJson,
		Os:     suite.ctx.Os,
		IoUtil: suite.ctx.IoUtil,
	}
	err = p.Write(&mockContext, tempFile.Name())
	if assert.NotNil(suite.T(), err) {
		assert.Errorf(suite.T(), err, "json marshal failure")
	}
}


