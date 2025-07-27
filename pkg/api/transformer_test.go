package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlitchTipSlackMessage(t *testing.T) {
	jsonStr := `{"alias": "GlitchTip", "text": "GlitchTip Alert", "attachments": [{"title": "*errors.errorString: Failed to setup database connection", "title_link": "https://xxx.com/devguard/issues/5", "text": null, "image_url": null, "color": "#e52b50", "fields": [{"title": "Project", "value": "devguard-api", "short": true}, {"title": "Environment", "value": "dev", "short": true}, {"title": "Release", "value": "0.11.1-439-g8f91aaa5-dirty", "short": false}], "mrkdown_in": ["text"]}], "sections": [{"activityTitle": "*errors.errorString: Failed to setup database connection", "activitySubtitle": "[View Issue DEVGUARD-API-5](https://xxx.com/devguard/issues/5)"}]}`
	msg, err := convertRawJsonToMatrixMessage(jsonStr, GlitchTip, mappingCodeGlitchTip)
	assert.NoError(t, err)

	// Expected message with the improved schema
	expectedPlain := `🔴 ERROR GlitchTip Alert
📦 devguard-api (dev)
🏷️ 0.11.1-439-g8f91aaa5-dirty

🐛 *errors.errorString: Failed to setup database connection
🔗 https://xxx.com/devguard/issues/5`

	expectedHtml := `<b>🔴 ERROR GlitchTip Alert</b><br/><code>devguard-api</code> <i>(dev)</i><br/>🏷️ <code>0.11.1-439-g8f91aaa5-dirty</code><br/><br/>🐛 *errors.errorString: Failed to setup database connection<br/>🔗 <a href="https://xxx.com/devguard/issues/5">View Issue</a>`

	assert.Equal(t, &MatrixMessage{
		Plain: expectedPlain,
		Html:  expectedHtml,
	}, msg)
}

func TestGlitchTipWarningMessage(t *testing.T) {
	jsonStr := `{"alias": "GlitchTip", "text": "GlitchTip Warning", "attachments": [{"title": "Deprecated API usage detected", "title_link": "https://xxx.com/devguard/issues/12", "text": null, "image_url": null, "color": "#ff9500", "fields": [{"title": "Project", "value": "frontend-app", "short": true}, {"title": "Environment", "value": "production", "short": true}], "mrkdown_in": ["text"]}]}`
	msg, err := convertRawJsonToMatrixMessage(jsonStr, GlitchTip, mappingCodeGlitchTip)
	assert.NoError(t, err)

	// Expected message with warning severity
	expectedPlain := `⚠️ WARNING GlitchTip Warning
📦 frontend-app (production)

🐛 Deprecated API usage detected
🔗 https://xxx.com/devguard/issues/12`

	expectedHtml := `<b>⚠️ WARNING GlitchTip Warning</b><br/><code>frontend-app</code> <i>(production)</i><br/><br/>🐛 Deprecated API usage detected<br/>🔗 <a href="https://xxx.com/devguard/issues/12">View Issue</a>`

	assert.Equal(t, &MatrixMessage{
		Plain: expectedPlain,
		Html:  expectedHtml,
	}, msg)
}

func TestBotKubeWebHookMessage(t *testing.T) {
	jsonStr := `{"source":"k8s-recommendation-events","data":{"APIVersion":"v1","Action":"","Cluster":"David-Test","Count":0,"Kind":"Pod","Level":"success","Messages":null,"Name":"nginx","Namespace":"default","Reason":"","Recommendations":["The 'latest' tag used in 'nginx' image of Pod 'default/nginx' container 'nginx' should be avoided."],"Resource":"v1/pods","TimeStamp":"2025-07-11T07:24:02Z","Title":"v1/pods created","Type":"create","Warnings":null},"timeStamp":"0001-01-01T00:00:00Z"}`
	msg, err := convertRawJsonToMatrixMessage(jsonStr, Botkube, mappingCodeBotKube)
	assert.NoError(t, err)

	// Expected message with the improved schema
	expectedPlain := `✅ SUCCESS Kubernetes Created
📦 **Pod/nginx** in **default**@David-Test

💡 Recommendations:
• The 'latest' tag used in 'nginx' image of Pod 'default/nginx' container 'nginx' should be avoided.`

	expectedHtml := `<b>✅ SUCCESS Kubernetes Created</b><br/>📦 <b>Pod/nginx</b> in <b>default</b>@<code>David-Test</code><br/><b>💡 Recommendations:</b><ul><li>The 'latest' tag used in 'nginx' image of Pod 'default/nginx' container 'nginx' should be avoided.</li></ul>`

	assert.Equal(t, &MatrixMessage{
		Plain: expectedPlain,
		Html:  expectedHtml,
	}, msg)
}

func TestBotKubeErrorWebHookMessage(t *testing.T) {
	jsonStr := `{"source":"k8s-err-events","data":{"APIVersion":"v1","Action":"","Cluster":"development","Count":1,"Kind":"Pod","Level":"error","Messages":["Failed to pull image \"nonexistentimage:latest\": failed to pull and unpack image \"docker.io/library/nonexistentimage:latest\": failed to resolve reference \"docker.io/library/nonexistentimage:latest\": pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed"],"Name":"error-pod","Namespace":"default","Reason":"Failed","Recommendations":null,"Resource":"v1/pods","TimeStamp":"2025-07-26T18:40:25Z","Title":"v1/pods error","Type":"error","Warnings":null},"timeStamp":"0001-01-01T00:00:00Z"}`
	msg, err := convertRawJsonToMatrixMessage(jsonStr, Botkube, mappingCodeBotKube)
	assert.NoError(t, err)

	// Expected message with the improved schema for error events
	expectedPlain := `🔴 ERROR Kubernetes Error
📦 **Pod/error-pod** in **default**@development

📋 Messages:
• Failed to pull image "nonexistentimage:latest": failed to pull and unpack image "docker.io/library/nonexistentimage:latest": failed to resolve reference "docker.io/library/nonexistentimage:latest": pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed`

	expectedHtml := `<b>🔴 ERROR Kubernetes Error</b><br/>📦 <b>Pod/error-pod</b> in <b>default</b>@<code>development</code><br/><br/><b>📋 Messages:</b><ul><li>Failed to pull image "nonexistentimage:latest": failed to pull and unpack image "docker.io/library/nonexistentimage:latest": failed to resolve reference "docker.io/library/nonexistentimage:latest": pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed</li></ul>`

	assert.Equal(t, &MatrixMessage{
		Plain: expectedPlain,
		Html:  expectedHtml,
	}, msg)
}
