local input = std.parseJson(std.extVar("input"));

{
  plain: input.text,
  details: input.attachments[0].title,
  html: "<b>" + self.plain + ":</b> " + self.details + " (<a href=\"" + input.attachments[0].title_link + "\">View Issue<a/>)",
}
