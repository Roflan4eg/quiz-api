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

func TestCreateGetDeleteAnswer_Success(t *testing.T) {
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
	questionID := int(question["id"].(float64))

	answerText := gofakeit.Word()
	answerData := map[string]interface{}{
		"text":    answerText,
		"user_id": gofakeit.UUID(),
	}

	body, _ = json.Marshal(answerData)
	req, err = http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/questions/%d/answers", st.BaseURL, questionID), bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var answer map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&answer)
	require.NoError(t, err)

	require.NotEmpty(t, answer["id"])
	assert.Equal(t, answerText, answer["text"])
	assert.Equal(t, float64(questionID), answer["question_id"])
	assert.Equal(t, answerData["user_id"], answer["user_id"])

	answerID := int(answer["id"].(float64))

	req, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/answers/%d", st.BaseURL, answerID), nil)
	require.NoError(t, err)

	resp, err = st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedAnswer map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&fetchedAnswer)
	require.NoError(t, err)

	assert.Equal(t, answer["id"], fetchedAnswer["id"])
	assert.Equal(t, answer["text"], fetchedAnswer["text"])
	assert.Equal(t, answer["question_id"], fetchedAnswer["question_id"])
	assert.Equal(t, answer["user_id"], fetchedAnswer["user_id"])

	req, err = http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/answers/%d", st.BaseURL, answerID), nil)
	require.NoError(t, err)

	resp, err = st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	req, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/answers/%d", st.BaseURL, answerID), nil)
	require.NoError(t, err)

	resp, err = st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCreateAnswer_NonExistentQuestion(t *testing.T) {
	ctx, st := suite.New(t)

	answerData := map[string]interface{}{
		"text":    gofakeit.Word(),
		"user_id": gofakeit.UUID(),
	}

	body, _ := json.Marshal(answerData)
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/questions/99999/answers", st.BaseURL), bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestMultipleAnswers_SameUserSameQuestion(t *testing.T) {
	ctx, st := suite.New(t)

	questionData := map[string]interface{}{
		"text": gofakeit.Question(),
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
	questionID := int(question["id"].(float64))

	userID := gofakeit.UUID()

	for i := 0; i < 3; i++ {
		answerData := map[string]interface{}{
			"text":    fmt.Sprintf("Answer %d: %s", i+1, gofakeit.Word()),
			"user_id": userID,
		}

		body, _ = json.Marshal(answerData)
		req, err = http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/questions/%d/answers", st.BaseURL, questionID), bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = st.Client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode, "Should allow multiple answers from same user")

		var answer map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&answer)
		require.NoError(t, err)

		assert.Equal(t, userID, answer["user_id"])
		assert.Equal(t, float64(questionID), answer["question_id"])
	}
}

func TestGetAnswer_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	req, err := http.NewRequestWithContext(ctx, "GET", st.BaseURL+"/answers/99999", nil)
	require.NoError(t, err)

	resp, err := st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteAnswer_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	req, err := http.NewRequestWithContext(ctx, "DELETE", st.BaseURL+"/answers/99999", nil)
	require.NoError(t, err)

	resp, err := st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCreateAnswer_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	questionData := map[string]interface{}{
		"text": gofakeit.Question(),
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
	questionID := int(question["id"].(float64))

	cases := []struct {
		name        string
		text        string
		userID      string
		expectedErr string
	}{
		{
			name:        "Create answer with empty text",
			text:        "",
			userID:      gofakeit.UUID(),
			expectedErr: "Answer text is required",
		},
		{
			name:        "Create answer with whitespace only",
			text:        "   ",
			userID:      gofakeit.UUID(),
			expectedErr: "Answer text is required",
		},
		{
			name:        "Create answer with empty user_id",
			text:        gofakeit.Word(),
			userID:      "",
			expectedErr: "User ID is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			answerData := map[string]interface{}{
				"text":    tc.text,
				"user_id": tc.userID,
			}

			body, _ = json.Marshal(answerData)
			req, err = http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/questions/%d/answers", st.BaseURL, questionID), bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err = st.Client.Do(req)
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

func TestCreateAnswer_MultipleUsersSameQuestion(t *testing.T) {
	ctx, st := suite.New(t)

	questionData := map[string]interface{}{
		"text": gofakeit.Question(),
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
	questionID := int(question["id"].(float64))

	userIDs := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

	for i, userID := range userIDs {
		answerData := map[string]interface{}{
			"text":    fmt.Sprintf("Answer from user %d: %s", i+1, gofakeit.Word()),
			"user_id": userID,
		}

		body, _ = json.Marshal(answerData)
		req, err = http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/questions/%d/answers", st.BaseURL, questionID), bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = st.Client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode, "Should allow answers from different users")

		var answer map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&answer)
		require.NoError(t, err)

		assert.Equal(t, userID, answer["user_id"])
		assert.Equal(t, float64(questionID), answer["question_id"])
	}
}
