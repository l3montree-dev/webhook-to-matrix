local input = std.parseJson(std.extVar("input"));

{
  plain: "Botkube Alert",
  html: "<b>" + self.plain + ":</b> " + input.source + " - " + input.data.Kind + " - " + input.data.Type + " - " + input.data.Recommendations
}
