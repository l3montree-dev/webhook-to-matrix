local input = std.parseJson(std.extVar("input"));

local message = if std.objectHas(input, "message") then input.message else "No message provided";
local link = if std.objectHas(input, "link") then input.link else "";

local plain = if link != "" then message + " " + link  else message;

local html = if link != "" then message + " <a href=\"" + link + "\">"+link+"</a>" else message;

{
  plain: plain,
  html: html
}