local input = std.parseJson(std.extVar("input"));

local data = input.data;
local cluster = if std.objectHas(data, "Cluster") then data.Cluster else "unknown";
local namespace = if std.objectHas(data, "Namespace") then data.Namespace else "default";
local kind = if std.objectHas(data, "Kind") then data.Kind else "Resource";
local name = if std.objectHas(data, "Name") then data.Name else "unknown";
local level = if std.objectHas(data, "Level") then data.Level else "info";
local type = if std.objectHas(data, "Type") then data.Type else "event";

// Helper functions
local formatLevel(level) = 
  if level == "error" then "🔴 ERROR"
  else if level == "warning" then "⚠️ WARNING"  
  else if level == "success" then "✅ SUCCESS"
  else if level == "info" then "ℹ️ INFO"
  else "📢 " + std.asciiUpper(level);

local formatType(type) =
  if type == "error" then "Error"
  else if type == "create" then "Created"
  else if type == "update" then "Updated"
  else if type == "delete" then "Deleted"
  else std.asciiUpper(type[0:1]) + type[1:];

local hasMessages = std.objectHas(data, "Messages") && data.Messages != null && std.length(data.Messages) > 0;
local hasRecommendations = std.objectHas(data, "Recommendations") && data.Recommendations != null && std.length(data.Recommendations) > 0;
local hasWarnings = std.objectHas(data, "Warnings") && data.Warnings != null && std.length(data.Warnings) > 0;

// Generate content sections
local messagesSection = if hasMessages then 
  if std.length(data.Messages) == 1 then
    "\n\n📋 " + data.Messages[0]
  else
    "\n\n📋 Messages:\n" + std.join("\n", ["• " + msg for msg in data.Messages])
  else "";

local recommendationsSection = if hasRecommendations then
  if std.length(data.Recommendations) == 1 then
    "\n\n💡 " + data.Recommendations[0]
  else
    "\n\n💡 Recommendations:\n" + std.join("\n", ["• " + rec for rec in data.Recommendations])
  else "";

local warningsSection = if hasWarnings then
  if std.length(data.Warnings) == 1 then
    "\n\n⚠️ " + data.Warnings[0]
  else
    "\n\n⚠️ Warnings:\n" + std.join("\n", ["• " + warn for warn in data.Warnings])
  else "";

// Resource info - highlight the most important parts
local resourceInfo = "📦 **" + kind + "/" + name + "** in **" + namespace + "**@" + cluster;

// Plain text format
local plainTitle = formatLevel(level) + " Kubernetes " + formatType(type);
local plainBody = resourceInfo + messagesSection + recommendationsSection + warningsSection;

// HTML format  
local htmlTitle = "<b>" + formatLevel(level) + " Kubernetes " + formatType(type) + "</b>";
local htmlResourceInfo = "📦 <b>" + kind + "/" + name + "</b> in <b>" + namespace + "</b>@<code>" + cluster + "</code>";
local htmlBody = htmlResourceInfo +
  (if hasMessages then 
    if std.length(data.Messages) == 1 then
      "<br/><br/>📋 " + data.Messages[0]
    else
      "<br/><br/><b>📋 Messages:</b><ul>" + std.join("", ["<li>" + msg + "</li>" for msg in data.Messages]) + "</ul>"
  else "") +
  (if hasRecommendations then 
    if std.length(data.Recommendations) == 1 then
      "<br/><br/>💡 " + data.Recommendations[0]
    else
      "<br/><b>💡 Recommendations:</b><ul>" + std.join("", ["<li>" + rec + "</li>" for rec in data.Recommendations]) + "</ul>"
  else "") +
  (if hasWarnings then 
    if std.length(data.Warnings) == 1 then
      "<br/><br/>⚠️ " + data.Warnings[0]
    else
      "<br/><b>⚠️ Warnings:</b><ul>" + std.join("", ["<li>" + warn + "</li>" for warn in data.Warnings]) + "</ul>"
  else "");

{
  plain: plainTitle + "\n" + plainBody,
  html: htmlTitle + "<br/>" + htmlBody
}
