local input = std.parseJson(std.extVar('input'));

// List of users to ignore (add usernames here)
local ignoredUsers = [
  // "example-user",
  // "bot-user",
];

// Check if user should be ignored
local shouldIgnoreUser(username) = 
  std.member(ignoredUsers, username);

// Format issue URL
local formatIssueUrl(repo_url, issue_number) = 
  repo_url + "/issues/" + issue_number;

// Generate message for issue events
local generateIssueMessage() = 
  local action = input.action;
  local issue = input.issue;
  local user = input.sender.login;
  local repo_name = input.repository.full_name;
  local issue_title = issue.title;
  local issue_body = if std.objectHas(issue, 'body') && issue.body != null then issue.body else "";
  local issue_url = issue.html_url;
  local issue_number = std.toString(issue.number);
  
  if shouldIgnoreUser(user) then
    null
  else
    local actionText = if action == "opened" then "🆕 opened"
                      else if action == "closed" then "✅ closed"
                      else if action == "reopened" then "🔄 reopened"
                      else action;
    
    // Include description only for "opened" action and if it exists
    local includeDescription = action == "opened" && issue_body != "";
    local truncatedDescription = if includeDescription && std.length(issue_body) > 300 then
                                   std.substr(issue_body, 0, 300) + "..."
                                 else if includeDescription then
                                   issue_body
                                 else
                                   "";
    
    local plainMessage = "[" + repo_name + "] Issue " + actionText + " by " + user + ": " + issue_title + 
                         (if includeDescription then "\n\nDescription: " + truncatedDescription else "") +
                         "\n" + issue_url;
    
    local htmlMessage = "<strong>[" + repo_name + "]</strong> Issue " + actionText + " by <strong>" + user + "</strong>: <a href=\"" + issue_url + "\">" + issue_title + "</a>" +
                        (if includeDescription then "<br/><br/><blockquote>" + truncatedDescription + "</blockquote>" else "");
    
    {
      plain: plainMessage,
      html: htmlMessage
    };

// Generate message for issue comment events
local generateIssueCommentMessage() = 
  local action = input.action;
  local comment = input.comment;
  local issue = input.issue;
  local user = input.sender.login;
  local repo_name = input.repository.full_name;
  local issue_title = issue.title;
  local issue_url = issue.html_url;
  local comment_url = comment.html_url;
  local comment_body = comment.body;
  
  if shouldIgnoreUser(user) then
    null
  else
    local actionText = if action == "created" then "💬 commented on"
                      else if action == "edited" then "✏️ edited comment on"
                      else if action == "deleted" then "🗑️ deleted comment on"
                      else action + " comment on";
    
    // Truncate comment if too long
    local truncatedComment = if std.length(comment_body) > 200 then
                               std.substr(comment_body, 0, 200) + "..."
                             else
                               comment_body;
    
    {
      plain: "[" + repo_name + "] " + user + " " + actionText + " issue: " + issue_title + "\nComment: " + truncatedComment + "\n" + comment_url,
      html: "<strong>[" + repo_name + "]</strong> <strong>" + user + "</strong> " + actionText + " issue: <a href=\"" + issue_url + "\">" + issue_title + "</a><br/><blockquote>" + truncatedComment + "</blockquote><a href=\"" + comment_url + "\">View comment</a>"
    };

// Main logic to determine event type and generate appropriate message
if std.objectHas(input, 'issue') && std.objectHas(input, 'comment') then
  // Issue comment event
  generateIssueCommentMessage()
else if std.objectHas(input, 'issue') && std.objectHas(input, 'action') then
  // Issue event
  generateIssueMessage()
else
  // Unsupported event type, return null to ignore
  null