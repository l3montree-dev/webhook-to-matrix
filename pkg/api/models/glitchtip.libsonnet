local input = std.parseJson(std.extVar("input"));

// Helper functions for safe field access
local getField(fields, title) = 
  local matches = [f.value for f in fields if f.title == title];
  if std.length(matches) > 0 then matches[0] else null;

local hasAttachments = std.objectHas(input, "attachments") && std.length(input.attachments) > 0;
local attachment = if hasAttachments then input.attachments[0] else {};
local hasFields = std.objectHas(attachment, "fields") && std.length(attachment.fields) > 0;
local fields = if hasFields then attachment.fields else [];

// Extract information
local alertText = if std.objectHas(input, "text") then input.text else "GlitchTip Alert";
local errorTitle = if std.objectHas(attachment, "title") then attachment.title else "Unknown Error";
local titleLink = if std.objectHas(attachment, "title_link") then attachment.title_link else null;
local project = getField(fields, "Project");
local environment = getField(fields, "Environment");
local release = getField(fields, "Release");
local color = if std.objectHas(attachment, "color") then attachment.color else "#e52b50";

// Format severity based on color
local getSeverity(color) =
  if color == "#e52b50" then "🔴 ERROR"
  else if color == "#ff9500" then "⚠️ WARNING"
  else if color == "#36a64f" then "✅ INFO"
  else "🔵 ALERT";

local severity = getSeverity(color);

// Build context information
local contextInfo = 
  (if project != null then "📦 " + project else "") +
  (if environment != null then " (" + environment + ")" else "") +
  (if release != null then "\n🏷️ " + release else "");

// Generate issue link
local issueLink = if titleLink != null then titleLink else "";
local hasLink = issueLink != "";

// Plain text format
local plainTitle = severity + " " + alertText;
local plainBody = (if contextInfo != "" then contextInfo + "\n\n" else "") +
  "🐛 " + errorTitle +
  (if hasLink then "\n🔗 " + issueLink else "");

// HTML format
local htmlTitle = "<b>" + severity + " " + alertText + "</b>";
local htmlContextInfo = 
  (if project != null then "<code>" + project + "</code>" else "") +
  (if environment != null then " <i>(" + environment + ")</i>" else "") +
  (if release != null then "<br/>🏷️ <code>" + release + "</code>" else "");

local htmlBody = (if contextInfo != "" then htmlContextInfo + "<br/><br/>" else "") +
  "🐛 " + errorTitle +
  (if hasLink then "<br/>🔗 <a href=\"" + issueLink + "\">View Issue</a>" else "");

{
  plain: plainTitle + "\n" + plainBody,
  html: htmlTitle + "<br/>" + htmlBody
}
