local input = std.parseJson(std.extVar("input"));

{
  plain: input.text,
  attachment: input.attachments[0],
  project: [
    f.value
    for f in self.attachment.fields
    if f.title == "Project"
  ][0],
  details: self.project + ": " + self.attachment.title,
  html: "<b>" + self.plain + ":</b> " + self.details + " (<a href=\"" + self.attachment.title_link + "\">View Issue<a/>)",
}
