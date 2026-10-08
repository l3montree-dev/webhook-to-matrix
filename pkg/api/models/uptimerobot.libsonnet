local input = std.parseJson(std.extVar("input"));

// UptimeRobot leaves variables that do not apply to an alert unreplaced (e.g. "*sslExpiryDate*" on a Down alert)
local isPlaceholder(value) =
  std.length(value) >= 2 && std.startsWith(value, "*") && std.endsWith(value, "*");

// Helper function for safe field access - empty and unreplaced values are treated as missing
local getField(obj, field, default=null) =
  if std.objectHas(obj, field) && obj[field] != null then
    local value = std.toString(obj[field]);
    if value == "" || isPlaceholder(value) then default else value
  else default;

// Extract information
local monitorName = getField(input, "monitorFriendlyName", "Unknown Monitor");
local monitorURL = getField(input, "monitorURL");
local alertType = getField(input, "alertType");
local alertTypeName = getField(input, "alertTypeFriendlyName", "Alert");
local alertDetails = getField(input, "alertDetails");
local alertDuration = getField(input, "alertDuration");
local sslExpiryDate = getField(input, "sslExpiryDate");
local sslExpiryDaysLeft = getField(input, "sslExpiryDaysLeft");
local domainExpireDate = getField(input, "domainExpireDate");
local dashboardUrl = getField(input, "dashboardUrl");

// Format status based on alert type (1: down, 2: up, 3: SSL & domain expiry)
local getStatus(alertType) =
  if alertType == "1" then "🔴 DOWN"
  else if alertType == "2" then "✅ UP"
  else if alertType == "3" then "⚠️ " + std.asciiUpper(alertTypeName)
  else "🔵 " + std.asciiUpper(alertTypeName);

local status = getStatus(alertType);

// Format a duration in seconds as e.g. "7d 5m"
local formatDuration(seconds) =
  local days = std.floor(seconds / 86400);
  local hours = std.floor((seconds % 86400) / 3600);
  local minutes = std.floor((seconds % 3600) / 60);
  local parts =
    (if days > 0 then [std.toString(days) + "d"] else []) +
    (if hours > 0 then [std.toString(hours) + "h"] else []) +
    (if minutes > 0 then [std.toString(minutes) + "m"] else []);
  if std.length(parts) > 0 then std.join(" ", parts) else std.toString(seconds) + "s";

local downtime = if alertType == "2" && alertDuration != null then formatDuration(std.parseInt(alertDuration)) else null;

// Only link the monitor URL if it is an actual URL (ping / port monitors only have a host)
local isLink(url) = std.startsWith(url, "http://") || std.startsWith(url, "https://");

local sslInfo = if sslExpiryDate != null then
  sslExpiryDate + (if sslExpiryDaysLeft != null then " (" + sslExpiryDaysLeft + " days left)" else "")
else null;

// Plain text format
local plainLines = [
  status + " " + monitorName,
] + (if monitorURL != null then ["🌐 " + monitorURL] else [])
  + (if alertDetails != null then ["💬 " + alertDetails] else [])
  + (if downtime != null then ["⏱️ Down for " + downtime] else [])
  + (if sslInfo != null then ["🔒 SSL certificate expires " + sslInfo] else [])
  + (if domainExpireDate != null then ["🌍 Domain expires " + domainExpireDate] else [])
  + (if dashboardUrl != null then ["🔗 " + dashboardUrl] else []);

// HTML format
local htmlLines = [
  "<b>" + status + " " + monitorName + "</b>",
] + (if monitorURL != null then
       ["🌐 " + (if isLink(monitorURL) then "<a href=\"" + monitorURL + "\">" + monitorURL + "</a>" else "<code>" + monitorURL + "</code>")]
     else [])
  + (if alertDetails != null then ["💬 " + alertDetails] else [])
  + (if downtime != null then ["⏱️ Down for <b>" + downtime + "</b>"] else [])
  + (if sslInfo != null then ["🔒 SSL certificate expires " + sslInfo] else [])
  + (if domainExpireDate != null then ["🌍 Domain expires " + domainExpireDate] else [])
  + (if dashboardUrl != null then ["🔗 <a href=\"" + dashboardUrl + "\">View in UptimeRobot</a>"] else []);

{
  plain: std.join("\n", plainLines),
  html: std.join("<br/>", htmlLines)
}
