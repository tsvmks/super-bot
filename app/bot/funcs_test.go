package bot

import (
	"bytes"
	"encoding/json"
	"github.com/radio-t/super-bot/app/bot/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFuncsBot_ReactionOnNewsRequest(t *testing.T) {
	mockHTTP := &mocks.HTTPClient{}
	b := NewFuncs(mockHTTP, "")

	articles := funcDesc{
			Title: "title1",
			Exampl:  "link1",
			returnType:   "VARCHAR",
		}
	articleJSON, err := json.Marshal(articles)
	require.NoError(t, err)

	mockHTTP.On("Do", mock.Anything).Return(&http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(articleJSON)),
	}, nil)

	require.Equal(
		t,
		Response{Text: "- [title1]\n\nСинтаксис:\n_link1_\nТип значения: \n- [все функции](https://help.krista.ru/kb/4217)", Send: true},
		b.OnMessage(Message{Text: "func!"}),
	)
}

func TestFuncsBot_ReactionOnNewsRequestAlt(t *testing.T) {
	mockHTTP := &mocks.HTTPClient{}
	b := NewFuncs(mockHTTP, "")

	article := funcDesc{
		Title: "title",
		Exampl:  "exampl",
	}
	articleJSON, err := json.Marshal(article)
	require.NoError(t, err)

	mockHTTP.On("Do", mock.Anything).Return(&http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(articleJSON)),
	}, nil)

	require.Equal(
		t,
		Response{Text: "- [title]\n\nСинтаксис:\n_exampl_\nТип значения: \n- [все функции](https://help.krista.ru/kb/4217)", Send: true},
		b.OnMessage(Message{Text: "func!"}),
	)
}

func TestFuncsBot_ReactionOnNewsRequestAlt_c(t *testing.T) {
	mockHTTP := &mocks.HTTPClient{}
	b := NewFuncs(mockHTTP, "")

	article := funcDesc{
		Title: "title",
		Exampl:  "exampl",
	}
	articleJSON, err := json.Marshal(article)
	require.NoError(t, err)

	mockHTTP.On("Do", mock.Anything).Return(&http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(articleJSON)),
	}, nil)

	require.Equal(
		t,
		Response{Text: "- [title]\n\nСинтаксис:\n_exampl_\nТип значения: \n- [все функции](https://help.krista.ru/kb/4217)", Send: true},
		b.OnMessage(Message{Text: "func! Concat"}),
	)
}

func TestFuncBot_ReactionOnUnexpectedMessage(t *testing.T) {
	mockHTTP := &mocks.HTTPClient{}
	b := NewFuncs(mockHTTP, "")
	require.Equal(t, Response{}, b.OnMessage(Message{Text: "unexpected"}))
}
