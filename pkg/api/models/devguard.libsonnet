local input = std.parseJson(std.extVar("input"));

// Helper functions for safe field access
local getField(obj, field, default=null) = 
  if std.objectHas(obj, field) then obj[field] else default;

// Extract main sections
local organization = getField(input, "organization", {});
local project = getField(input, "project", {});
local asset = getField(input, "asset", {});
local assetVersion = getField(input, "assetVersion", {});
local payload = getField(input, "payload", []);
local eventType = getField(input, "type", "unknown");

// Organization info
local orgName = getField(organization, "name", "Unknown Organization");
local orgSlug = getField(organization, "slug", "unknown");

// Project info  
local projectName = getField(project, "name", "Unknown Project");
local projectSlug = getField(project, "slug", "unknown");

// Asset info
local assetName = getField(asset, "name", "Unknown Asset");
local assetSlug = getField(asset, "slug", "unknown");

// Asset version info
local versionName = getField(assetVersion, "name", "unknown");

// Helper function to format severity based on risk assessment
local formatSeverity(riskAssessment) =
  if riskAssessment >= 90 then "🔴 CRITICAL"
  else if riskAssessment >= 70 then "🟠 HIGH"
  else if riskAssessment >= 40 then "🟡 MEDIUM"
  else if riskAssessment >= 10 then "🟢 LOW"
  else "ℹ️ INFO";

// Helper function to format risk assessment with proper rounding
local formatRiskAssessment(rawRisk) =
  if rawRisk > 0 then
    // Convert to string and use a more reliable string truncation approach
    local str = std.toString(rawRisk);
    local dotPos = std.findSubstr(".", str);
    if std.length(dotPos) > 0 then
      local beforeDot = str[0:dotPos[0]];
      local afterDot = str[dotPos[0]+1:];
      // Take only first decimal place and handle edge cases
      if std.length(afterDot) >= 1 then
        local firstDecimal = std.parseInt(afterDot[0:1]);
        // Simple round if there are more decimals
        if std.length(afterDot) > 1 && std.parseInt(afterDot[1:2]) >= 5 then
          if firstDecimal == 9 then
            std.toString(std.parseInt(beforeDot) + 1) + ".0"
          else
            beforeDot + "." + std.toString(firstDecimal + 1)
        else
          beforeDot + "." + std.toString(firstDecimal)
      else
        str
    else
      str + ".0"
  else "0";

// Helper function to extract package name from PURL
local extractPackageName(purl) =
  if purl != null then
    local parts = std.split(purl, "/");
    if std.length(parts) > 0 then
      local lastPart = parts[std.length(parts) - 1];
      local nameVersion = std.split(lastPart, "@");
      if std.length(nameVersion) > 0 then nameVersion[0] else "unknown"
    else "unknown"
  else "unknown";

// Process vulnerabilities
local vulnCount = std.length(payload);
local hasVulns = vulnCount > 0;

// Generate vulnerability sections
local generateVulnSummary() = 
  if !hasVulns then ""
  else if vulnCount == 1 then
    local vuln = payload[0];
    local severity = formatSeverity(getField(vuln, "riskAssessment", 0));
    local rawRisk = getField(vuln, "rawRiskAssessment", 0);
    local cveId = getField(vuln, "cveID", "No CVE");
    local packageName = extractPackageName(getField(vuln, "componentPurl"));
    local fixedVersion = getField(vuln, "componentFixedVersion");
    local cveObj = getField(vuln, "cve");
    local cveDescription = if cveObj != null then getField(cveObj, "description") else null;
    local cvssScore = if cveObj != null then getField(cveObj, "cvss", 0) else 0;
    
    "\n\n" + severity + " " + cveId +
    (if cvssScore > 0 then " (CVSS: " + std.toString(cvssScore) + ")" else "") +
    (if rawRisk > 0 then " [Risk: " + formatRiskAssessment(rawRisk) + "]" else "") +
    "\n📦 " + packageName + 
    (if fixedVersion != null then " → Fix: " + fixedVersion else " (No fix available)") +
    (if cveDescription != null then "\n💬 " + cveDescription else "")
  else
    local criticalCount = std.length([v for v in payload if getField(v, "riskAssessment", 0) >= 90]);
    local highCount = std.length([v for v in payload if getField(v, "riskAssessment", 0) >= 70 && getField(v, "riskAssessment", 0) < 90]);
    local mediumCount = std.length([v for v in payload if getField(v, "riskAssessment", 0) >= 40 && getField(v, "riskAssessment", 0) < 70]);
    local lowCount = std.length([v for v in payload if getField(v, "riskAssessment", 0) < 40]);
    
    "\n\n📊 " + std.toString(vulnCount) + " vulnerabilities detected:" +
    (if criticalCount > 0 then "\n🔴 " + std.toString(criticalCount) + " Critical" else "") +
    (if highCount > 0 then "\n🟠 " + std.toString(highCount) + " High" else "") +
    (if mediumCount > 0 then "\n🟡 " + std.toString(mediumCount) + " Medium" else "") +
    (if lowCount > 0 then "\n🟢 " + std.toString(lowCount) + " Low" else "");

local generateVulnDetails() =
  if !hasVulns then ""
  else if vulnCount == 1 then ""  // Details already in summary for single vuln
  else
    local sortedVulns = std.reverse(std.sort(payload, function(v) getField(v, "riskAssessment", 0)));
    "\n\n🔍 Top vulnerabilities:" +
    std.join("", [
      "\n• " + formatSeverity(getField(v, "riskAssessment", 0)) + " " + 
      getField(v, "cveID", "No CVE") + " in " + 
      extractPackageName(getField(v, "componentPurl")) +
      (local cveObj = getField(v, "cve"); local cvssScore = if cveObj != null then getField(cveObj, "cvss", 0) else 0;
       if cvssScore > 0 then " (CVSS: " + std.toString(cvssScore) + ")" else "") +
      (local rawRisk = getField(v, "rawRiskAssessment", 0);
       if rawRisk > 0 then " [Risk: " + formatRiskAssessment(rawRisk) + "]" else "")
      for v in sortedVulns[0:std.min(5, vulnCount)]
    ]);

// Context information
local contextInfo = "📁 " + projectName + " / " + assetName + 
  (if versionName != "unknown" then " (" + versionName + ")" else "");

// Generate title
local title = if eventType == "dependencyVulnerabilities" then 
  "🛡️ DevGuard Security Scan"
else "🛡️ DevGuard Alert";

// Generate vulnerability summary
local vulnSummary = generateVulnSummary();
local vulnDetails = generateVulnDetails();

// Plain text format
local plainTitle = title;
local plainBody = contextInfo + vulnSummary + vulnDetails;

// HTML format
local htmlTitle = "<b>" + title + "</b>";
local htmlContextInfo = "📁 <b>" + projectName + "</b> / <b>" + assetName + "</b>" + 
  (if versionName != "unknown" then " <i>(" + versionName + ")</i>" else "");

local htmlVulnSummary = if !hasVulns then ""
  else if vulnCount == 1 then
    local vuln = payload[0];
    local severity = formatSeverity(getField(vuln, "riskAssessment", 0));
    local rawRisk = getField(vuln, "rawRiskAssessment", 0);
    local cveId = getField(vuln, "cveID", "No CVE");
    local packageName = extractPackageName(getField(vuln, "componentPurl"));
    local fixedVersion = getField(vuln, "componentFixedVersion");
    local cveObj = getField(vuln, "cve");
    local cveDescription = if cveObj != null then getField(cveObj, "description") else null;
    local cvssScore = if cveObj != null then getField(cveObj, "cvss", 0) else 0;
    
    "<br/><br/>" + severity + " <code>" + cveId + "</code>" +
    (if cvssScore > 0 then " <i>(CVSS: " + std.toString(cvssScore) + ")</i>" else "") +
    (if rawRisk > 0 then " <i>[Risk: " + formatRiskAssessment(rawRisk) + "]</i>" else "") +
    "<br/>📦 <code>" + packageName + "</code>" + 
    (if fixedVersion != null then " → Fix: <code>" + fixedVersion + "</code>" else " <i>(No fix available)</i>") +
    (if cveDescription != null then "<br/>💬 <i>" + cveDescription + "</i>" else "")
  else
    local criticalCount = std.length([v for v in payload if getField(v, "riskAssessment", 0) >= 90]);
    local highCount = std.length([v for v in payload if getField(v, "riskAssessment", 0) >= 70 && getField(v, "riskAssessment", 0) < 90]);
    local mediumCount = std.length([v for v in payload if getField(v, "riskAssessment", 0) >= 40 && getField(v, "riskAssessment", 0) < 70]);
    local lowCount = std.length([v for v in payload if getField(v, "riskAssessment", 0) < 40]);
    
    "<br/><br/>📊 <b>" + std.toString(vulnCount) + " vulnerabilities detected:</b>" +
    (if criticalCount > 0 then "<br/>🔴 " + std.toString(criticalCount) + " Critical" else "") +
    (if highCount > 0 then "<br/>🟠 " + std.toString(highCount) + " High" else "") +
    (if mediumCount > 0 then "<br/>🟡 " + std.toString(mediumCount) + " Medium" else "") +
    (if lowCount > 0 then "<br/>🟢 " + std.toString(lowCount) + " Low" else "");

local htmlVulnDetails = if !hasVulns then ""
  else if vulnCount == 1 then ""
  else
    local sortedVulns = std.reverse(std.sort(payload, function(v) getField(v, "riskAssessment", 0)));
    "<br/><br/>🔍 <b>Top vulnerabilities:</b><ul>" +
    std.join("", [
      "<li>" + formatSeverity(getField(v, "riskAssessment", 0)) + " <code>" + 
      getField(v, "cveID", "No CVE") + "</code> in <code>" + 
      extractPackageName(getField(v, "componentPurl")) + "</code>" +
      (local cveObj = getField(v, "cve"); local cvssScore = if cveObj != null then getField(cveObj, "cvss", 0) else 0;
       if cvssScore > 0 then " <i>(CVSS: " + std.toString(cvssScore) + ")</i>" else "") +
      (local rawRisk = getField(v, "rawRiskAssessment", 0);
       if rawRisk > 0 then " <i>[Risk: " + formatRiskAssessment(rawRisk) + "]</i>" else "") + "</li>"
      for v in sortedVulns[0:std.min(5, vulnCount)]
    ]) + "</ul>";

local htmlBody = htmlContextInfo + htmlVulnSummary + htmlVulnDetails;

{
  plain: plainTitle + "\n" + plainBody,
  html: htmlTitle + "<br/>" + htmlBody
}
