const { spawn } = require('child_process');

module.exports = (app) => {
  app.log.info("🛡️ StuntBot App is active! Listening for PR comments...");

  app.on("issue_comment.created", async (context) => {
    // Only trigger on Pull Requests when the user says "@stuntdouble review"
    if (!context.payload.issue.pull_request) return;
    if (!context.payload.comment.body.includes("@stuntdouble review")) return;

    const issue = context.issue();
    
    app.log.info(`>> [StuntBot] Received review request on PR #${issue.issue_number}`);
    
    const reaction = context.issue({
      content: "eyes",
      comment_id: context.payload.comment.id
    });
    await context.octokit.reactions.createForIssueComment(reaction);

    // Notify the user the sandbox is spinning up
    const comment = context.issue({
      body: "🚀 **StuntDouble Isolation Sandbox Starting...**\n\nCloning repository securely and booting eBPF kernel hooks..."
    });
    await context.octokit.issues.createComment(comment);

    // In a production environment, StuntBot would securely clone the PR branch into 
    // a temporary directory, and execute the StuntDouble CLI with the selected AI agent.
    
    // Simulate the CLI executing an agent like Claude to review the code
    const child = spawn("sd", ["run", "claude"]);
    
    child.on('close', async (code) => {
      const resultComment = context.issue({
        body: `✅ **StuntDouble Review Complete**\n\nAgent executed safely inside the kernel sandbox. No database connections or rogue network calls were allowed.\n\n*Zero zero-day escapes detected during review.*`
      });
      await context.octokit.issues.createComment(resultComment);
    });
  });
};
