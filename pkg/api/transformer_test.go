package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlitchTipSlackMessage(t *testing.T) {
	jsonStr := `{"alias": "GlitchTip", "text": "GlitchTip Alert", "attachments": [{"title": "*errors.errorString: Failed to setup database connection", "title_link": "https://xxx.com/devguard/issues/5", "text": null, "image_url": null, "color": "#e52b50", "fields": [{"title": "Project", "value": "devguard-api", "short": true}, {"title": "Environment", "value": "dev", "short": true}, {"title": "Release", "value": "0.11.1-439-g8f91aaa5-dirty", "short": false}], "mrkdown_in": ["text"]}], "sections": [{"activityTitle": "*errors.errorString: Failed to setup database connection", "activitySubtitle": "[View Issue DEVGUARD-API-5](https://xxx.com/devguard/issues/5)"}]}`
	msg, err := convertRawJsonToMatrixMessage(jsonStr, GlitchTip)
	assert.NoError(t, err)
	assert.Equal(t, msg, &MatrixMessage{
		Plain: "GlitchTip Alert",
		Html:  "<b>GlitchTip Alert:</b> devguard-api: *errors.errorString: Failed to setup database connection (<a href=\"https://xxx.com/devguard/issues/5\">View Issue<a/>)",
	})
}

func TestBotKubeWebHookMessage(t *testing.T) {
	jsonStr := `{"source":"k8s-recommendation-events","data":{"APIVersion":"v1","Action":"","Cluster":"David-Test","Count":0,"Kind":"Pod","Level":"success","Messages":null,"Name":"nginx","Namespace":"default","Reason":"","Recommendations":["The 'latest' tag used in 'nginx' image of Pod 'default/nginx' container 'nginx' should be avoided."],"Resource":"v1/pods","TimeStamp":"2025-07-11T07:24:02Z","Title":"v1/pods created","Type":"create","Warnings":null},"timeStamp":"0001-01-01T00:00:00Z"}`
	msg, err := convertRawJsonToMatrixMessage(jsonStr, Botkube)
	assert.NoError(t, err)
	assert.Equal(t, msg, &MatrixMessage{
		Plain: "Botkube Alert",
		Html:  "<b>Botkube Alert:</b> k8s-recommendation-events - Pod - create - [\"The 'latest' tag used in 'nginx' image of Pod 'default/nginx' container 'nginx' should be avoided.\"]",
	})
}
