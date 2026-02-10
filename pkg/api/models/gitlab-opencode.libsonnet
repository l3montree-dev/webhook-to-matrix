local input = std.parseJson(std.extVar('input'));

// List of users to ignore (add usernames here)
local ignoredUsers = [
  // "example-user",
  // "bot-user",
];

// Check if user should be ignored
local shouldIgnoreUser(username) = 
  std.member(ignoredUsers, username);

// Generate message for issue events
local generateIssueMessage() = 
  local attrs = input.object_attributes;
  local action = attrs.action;
  local user = input.user.username;
  local repo_name = input.project.path_with_namespace;
  local issue_title = attrs.title;
  local issue_body = if std.objectHas(attrs, 'description') && attrs.description != null then attrs.description else "";
  local issue_url = attrs.url;
  local issue_iid = std.toString(attrs.iid);
  
  if shouldIgnoreUser(user) then
    null
  else
    local actionText = if action == "open" then "🆕 opened"
                      else if action == "close" then "✅ closed"
                      else if action == "reopen" then "🔄 reopened"
                      else if action == "update" then "✏️ updated"
                      else action;
    
    // Include description only for "open" action and if it exists
    local includeDescription = action == "open" && issue_body != "";
    local truncatedDescription = if includeDescription && std.length(issue_body) > 300 then
                                   std.substr(issue_body, 0, 300) + "..."
                                 else if includeDescription then
                                   issue_body
                                 else
                                   "";
    
    local plainMessage = "#### OpenCode\n[" + repo_name + "] Issue #" + issue_iid + " " + actionText + " by " + user + ": " + issue_title + 
                         (if includeDescription then "\n\nDescription: " + truncatedDescription else "") +
                         "\n" + issue_url;
    
    local htmlMessage = "<h4>OpenCode</h4><strong>[" + repo_name + "]</strong> Issue #" + issue_iid + " " + actionText + " by <strong>" + user + "</strong>: <a href=\"" + issue_url + "\">" + issue_title + "</a>" +
                        (if includeDescription then "<br/><br/><blockquote>" + truncatedDescription + "</blockquote>" else "");
    
    {
      plain: plainMessage,
      html: htmlMessage
    };

// Generate message for note (comment) events on issues
local generateNoteMessage() = 
  local attrs = input.object_attributes;
  local user = input.user.username;
  local repo_name = input.project.path_with_namespace;
  local note_body = attrs.note;
  local note_url = attrs.url;
  local noteable_type = attrs.noteable_type;
  
  if shouldIgnoreUser(user) then
    null
  else if noteable_type == "Issue" then
    local issue = input.issue;
    local issue_title = issue.title;
    local issue_url = input.project.web_url + "/-/issues/" + std.toString(issue.iid);
    
    // Truncate comment if too long
    local truncatedComment = if std.length(note_body) > 200 then
                               std.substr(note_body, 0, 200) + "..."
                             else
                               note_body;
    
    {
      plain: "#### OpenCode\n[" + repo_name + "] " + user + " 💬 commented on issue #" + std.toString(issue.iid) + ": " + issue_title + "\nComment: " + truncatedComment + "\n" + note_url,
      html: "<h4>OpenCode</h4><strong>[" + repo_name + "]</strong> <strong>" + user + "</strong> 💬 commented on issue: <a href=\"" + issue_url + "\">#" + std.toString(issue.iid) + " " + issue_title + "</a><br/><blockquote>" + truncatedComment + "</blockquote><a href=\"" + note_url + "\">View comment</a>"
    }
  else if noteable_type == "MergeRequest" then
    local mr = input.merge_request;
    local mr_title = mr.title;
    local mr_url = input.project.web_url + "/-/merge_requests/" + std.toString(mr.iid);
    
    // Truncate comment if too long
    local truncatedComment = if std.length(note_body) > 200 then
                               std.substr(note_body, 0, 200) + "..."
                             else
                               note_body;
    
    {
      plain: "#### OpenCode\n[" + repo_name + "] " + user + " 💬 commented on MR !" + std.toString(mr.iid) + ": " + mr_title + "\nComment: " + truncatedComment + "\n" + note_url,
      html: "<h4>OpenCode</h4><strong>[" + repo_name + "]</strong> <strong>" + user + "</strong> 💬 commented on MR: <a href=\"" + mr_url + "\">!" + std.toString(mr.iid) + " " + mr_title + "</a><br/><blockquote>" + truncatedComment + "</blockquote><a href=\"" + note_url + "\">View comment</a>"
    }
  else
    null;

// Main logic to determine event type and generate appropriate message
local object_kind = if std.objectHas(input, 'object_kind') then input.object_kind else null;

if object_kind == "issue" then
  generateIssueMessage()
else if object_kind == "note" then
  generateNoteMessage()
else
  // Unsupported event type, return null to ignore
  null