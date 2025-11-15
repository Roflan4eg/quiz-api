package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Roflan4eg/quiz-api/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGetQuestion_Success(t *testing.T) {
	ctx, st := suite.New(t)

	questionText := gofakeit.Question()
	questionData := map[string]interface{}{
		"text": questionText,
	}

	body, _ := json.Marshal(questionData)
	req, err := http.NewRequestWithContext(ctx, "POST", st.BaseURL+"/questions", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var question map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&question)
	require.NoError(t, err)

	require.NotEmpty(t, question["id"])
	assert.Equal(t, questionText, question["text"])

	questionID := int(question["id"].(float64))

	req, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/questions/%d", st.BaseURL, questionID), nil)
	require.NoError(t, err)

	resp, err = st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedQuestion map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&fetchedQuestion)
	require.NoError(t, err)

	assert.Equal(t, question["id"], fetchedQuestion["id"])
	assert.Equal(t, question["text"], fetchedQuestion["text"])
}

func TestCreateQuestion_Duplicate(t *testing.T) {
	ctx, st := suite.New(t)

	questionText := gofakeit.Question()
	questionData := map[string]interface{}{
		"text": questionText,
	}

	body, _ := json.Marshal(questionData)
	req, err := http.NewRequestWithContext(ctx, "POST", st.BaseURL+"/questions", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var question1 map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&question1)
	require.NotEmpty(t, question1["id"])

	body, _ = json.Marshal(questionData)
	req, err = http.NewRequestWithContext(ctx, "POST", st.BaseURL+"/questions", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var question2 map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&question2)
	require.NotEmpty(t, question2["id"])

	assert.NotEqual(t, question1["id"], question2["id"])
}

func TestCreateQuestion_EmptyText(t *testing.T) {
	ctx, st := suite.New(t)

	questionData := map[string]interface{}{
		"text": "",
	}

	body, _ := json.Marshal(questionData)
	req, err := http.NewRequestWithContext(ctx, "POST", st.BaseURL+"/questions", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetQuestion_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	req, err := http.NewRequestWithContext(ctx, "GET", st.BaseURL+"/questions/99999", nil)
	require.NoError(t, err)

	resp, err := st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCreateQuestion_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	cases := []struct {
		name        string
		text        string
		expectedErr string
	}{
		{
			name:        "Create question with empty text",
			text:        "",
			expectedErr: "Question text is required",
		},
		{
			name:        "Create question with whitespace only",
			text:        "   ",
			expectedErr: "Question text is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			questionData := map[string]interface{}{
				"text": tc.text,
			}

			body, _ := json.Marshal(questionData)
			req, err := http.NewRequestWithContext(ctx, "POST", st.BaseURL+"/questions", bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := st.Client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var errorResp map[string]interface{}
			if err = json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
				if errorMsg, exists := errorResp["error"]; exists {
					assert.Contains(t, errorMsg.(string), tc.expectedErr)
				}
			}
		})
	}
}
