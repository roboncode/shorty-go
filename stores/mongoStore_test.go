package stores

import (
	"testing"
)

/*
	mockThing := Thing{}

	mockDAL := &mocks.DataAccessLayer{}
	mockDAL.On("Insert", mockThing).Return(nil)

	actual := Bar(mockDAL)

	mockDAL.AssertExpectations(t)

	expect.Equal(t, mockThing, actual, "should return a Thing")
*/

func TestMongoStore_Create(t *testing.T) {
	//link := models.Link{}
	//
	//mockStore := &mocks.Store{}
	//mockStore.On("Create", link).Return(nil)
	//
	//actual := MongoStore(mockStore)
	//
	//link, err := mockStore.Create("abcde", "https://roboncode.com")
	////if err != nil {
	////	panic(err)
	////}
	//
	//assert.Equal(t, link.Code, "abcde", "MongoStore link.Code did not match")
	//assert.Equal(t, link.LongUrl, "https://roboncode.com", "MongoStore link.LongUrl did not match")
}
