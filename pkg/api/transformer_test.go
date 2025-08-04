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

💡 The 'latest' tag used in 'nginx' image of Pod 'default/nginx' container 'nginx' should be avoided.`

	expectedHtml := `<b>✅ SUCCESS Kubernetes Created</b><br/>📦 <b>Pod/nginx</b> in <b>default</b>@<code>David-Test</code><br/><br/>💡 The 'latest' tag used in 'nginx' image of Pod 'default/nginx' container 'nginx' should be avoided.`

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

📋 Failed to pull image "nonexistentimage:latest": failed to pull and unpack image "docker.io/library/nonexistentimage:latest": failed to resolve reference "docker.io/library/nonexistentimage:latest": pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed`

	expectedHtml := `<b>🔴 ERROR Kubernetes Error</b><br/>📦 <b>Pod/error-pod</b> in <b>default</b>@<code>development</code><br/><br/>📋 Failed to pull image "nonexistentimage:latest": failed to pull and unpack image "docker.io/library/nonexistentimage:latest": failed to resolve reference "docker.io/library/nonexistentimage:latest": pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed`

	assert.Equal(t, &MatrixMessage{
		Plain: expectedPlain,
		Html:  expectedHtml,
	}, msg)
}

func TestDevGuardDependencyVulnerabilityMessage(t *testing.T) {
	jsonStr := `{"organization":{"id":"3fd43312-7ada-4be5-a5bf-7fe340a3be8a","name":"WetterOnline","slug":"wetteronline"},"project":{"id":"ce298219-814c-4e54-bc58-a0a5b2b08973","name":"Example Project","slug":"example-project"},"asset":{"id":"856ab205-cc6d-4c49-878d-4047102ffa33","name":"Example Asset","slug":"example-asset"},"assetVersion":{"name":"example-version","slug":"example-version"},"payload":[{"id":"dep-vuln-001","cve":{"cve":"CVE-2021-44228","description":"Apache Log4j2 <=2.14.1 JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints.","cvss":10,"severity":"critical"},"cveID":"CVE-2021-44228","componentPurl":"pkg:maven/org.apache.logging.log4j/log4j-core@2.14.1","componentFixedVersion":"2.15.0","riskAssessment":95,"rawRiskAssessment":9.8,"priority":1}],"type":"dependencyVulnerabilities"}`
	msg, err := convertRawJsonToMatrixMessage(jsonStr, DevGuard, mappingCodeDevGuard)
	assert.NoError(t, err)

	// Expected message for single critical vulnerability with CVE details and raw risk
	expectedPlain := `🛡️ DevGuard Security Scan
📁 Example Project / Example Asset (example-version)

🔴 CRITICAL CVE-2021-44228 (CVSS: 10) [Risk: 9.8]
📦 log4j-core → Fix: 2.15.0
💬 Apache Log4j2 <=2.14.1 JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints.`

	expectedHtml := `<b>🛡️ DevGuard Security Scan</b><br/>📁 <b>Example Project</b> / <b>Example Asset</b> <i>(example-version)</i><br/><br/>🔴 CRITICAL <code>CVE-2021-44228</code> <i>(CVSS: 10)</i> <i>[Risk: 9.8]</i><br/>📦 <code>log4j-core</code> → Fix: <code>2.15.0</code><br/>💬 <i>Apache Log4j2 <=2.14.1 JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints.</i>`

	assert.Equal(t, &MatrixMessage{
		Plain: expectedPlain,
		Html:  expectedHtml,
	}, msg)
}

func TestDevGuardMultipleVulnerabilitiesMessage(t *testing.T) {
	jsonStr := `{"organization":{"name":"WetterOnline","slug":"wetteronline"},"project":{"name":"Example Project","slug":"example-project"},"asset":{"name":"Example Asset","slug":"example-asset"},"assetVersion":{"name":"v1.0.0"},"payload":[{"cve":{"cvss":10},"cveID":"CVE-2021-44228","componentPurl":"pkg:maven/org.apache.logging.log4j/log4j-core@2.14.1","riskAssessment":95,"rawRiskAssessment":9.8,"priority":1},{"cve":{"cvss":7.5},"cveID":"CVE-2022-12345","componentPurl":"pkg:npm/lodash@4.0.0","riskAssessment":75,"rawRiskAssessment":7.2,"priority":2},{"cve":{"cvss":4.5},"cveID":"CVE-2023-56789","componentPurl":"pkg:npm/express@3.0.0","riskAssessment":45,"rawRiskAssessment":4.1,"priority":3}],"type":"dependencyVulnerabilities"}`
	msg, err := convertRawJsonToMatrixMessage(jsonStr, DevGuard, mappingCodeDevGuard)
	assert.NoError(t, err)

	// Expected message for multiple vulnerabilities with CVSS and raw risk scores
	expectedPlain := `🛡️ DevGuard Security Scan
📁 Example Project / Example Asset (v1.0.0)

📊 3 vulnerabilities detected:
🔴 1 Critical
🟠 1 High
🟡 1 Medium

🔍 Top vulnerabilities:
• 🔴 CRITICAL CVE-2021-44228 in log4j-core (CVSS: 10) [Risk: 9.8]
• 🟠 HIGH CVE-2022-12345 in lodash (CVSS: 7.5) [Risk: 7.2]
• 🟡 MEDIUM CVE-2023-56789 in express (CVSS: 4.5) [Risk: 4.1]`

	expectedHtml := `<b>🛡️ DevGuard Security Scan</b><br/>📁 <b>Example Project</b> / <b>Example Asset</b> <i>(v1.0.0)</i><br/><br/>📊 <b>3 vulnerabilities detected:</b><br/>🔴 1 Critical<br/>🟠 1 High<br/>🟡 1 Medium<br/><br/>🔍 <b>Top vulnerabilities:</b><ul><li>🔴 CRITICAL <code>CVE-2021-44228</code> in <code>log4j-core</code> <i>(CVSS: 10)</i> <i>[Risk: 9.8]</i></li><li>🟠 HIGH <code>CVE-2022-12345</code> in <code>lodash</code> <i>(CVSS: 7.5)</i> <i>[Risk: 7.2]</i></li><li>🟡 MEDIUM <code>CVE-2023-56789</code> in <code>express</code> <i>(CVSS: 4.5)</i> <i>[Risk: 4.1]</i></li></ul>`

	assert.Equal(t, &MatrixMessage{
		Plain: expectedPlain,
		Html:  expectedHtml,
	}, msg)
}
